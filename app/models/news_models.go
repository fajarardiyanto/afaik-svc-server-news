package models

import (
	"context"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/fajarardiyanto/afaik-svc-server-news/internal/config"
	"github.com/fajarardiyanto/afaik-svc-server-news/internal/entity"
	"github.com/fajarardiyanto/afaik-svc-server-news/internal/interfaces"
	tr "github.com/fajarardiyanto/flt-go-tracer/lib/jaeger"
	"github.com/fajarardiyanto/flt-go-utils/pagination"
	"github.com/fajarardiyanto/flt-go-utils/parser"
	pb "github.com/fajarardiyanto/module-proto/go/services/news"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"os"
	"strconv"
	"time"
)

type Game struct{}

func NewGame() interfaces.GameRepository {
	return &Game{}
}

func (g *Game) Store(ctx context.Context, req *pb.NewsRequest, span opentracing.Span) (res *emptypb.Empty, err error) {
	sp := tr.CreateSubChildSpan(span, "Store News Models")
	defer sp.Finish()

	res = &emptypb.Empty{}

	tr.LogRequest(sp, req)

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.GetConfig().Timeout)*time.Millisecond)
	defer cancel()

	finish := make(chan bool)
	errCh := hystrix.Go("store_news", func() error {
		query := `
			INSERT INTO news_bulk (rss_id, title, link, image_url, author, category_source, category, category_id, 
			                       pubdate, permalink, content, source, tags, publish, is_headline, is_popular, created_at, sorting, exclusive)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
		`

		tx := config.GetDBConn().Orm().WithContext(ctx).Debug().
			Exec(
				query,
				req.RssID,
				req.Title,
				req.Link,
				req.ImageURL,
				req.Author,
				req.CategorySource,
				req.Category,
				req.CategoryID,
				time.Now().UTC(),
				req.Permalink,
				req.Content,
				req.Source,
				req.Tags,
				req.Publish,
				req.IsHeadline,
				req.IsPopular,
				time.Now().UTC(),
				req.Sorting,
				req.Exclusive,
			)
		if tx.Error != nil {
			return tx.Error
		}
		finish <- true

		return nil
	}, nil)

	select {
	case <-finish:
		tr.LogResponse(sp, res)
		return res, nil
	case err = <-errCh:
		tr.LogError(sp, err)
		return nil, err
	}
}

func (*Game) Get(ctx context.Context, req *pb.GetNewsRequest, span opentracing.Span) (res *pb.NewsResponse, err error) {
	sp := tr.CreateSubChildSpan(span, "Get News Models")
	defer sp.Finish()

	res = &pb.NewsResponse{}

	tr.LogRequest(sp, req)

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.GetConfig().Timeout)*time.Millisecond)
	defer cancel()

	resCh := make(chan []entity.NewsModel)
	errCh := hystrix.Go("get_news", func() error {
		query := `
			SELECT 
				id, rss_id, title, link, image_url, author, category_source, category, category_id, 
				pubdate, permalink, content, source, tags, publish, is_headline, is_popular, created_at, sorting, exclusive
			FROM
				news_bulk
			LIMIT ? OFFSET ?
		`
		tr.LogObject(sp, "query", query)

		var result []entity.NewsModel
		tx := config.GetDBConn().Orm().WithContext(ctx).Debug().
			Raw(query, req.Limit, pagination.GetPage(req.Limit, req.Offset)).Model(&entity.NewsModel{}).Scan(&result)

		if tx.Error != nil {
			return status.Error(codes.FailedPrecondition, tx.Error.Error())
		}

		if len(result) == 0 {
			return status.Error(codes.NotFound, "News not found")
		}

		resCh <- result
		return nil
	}, nil)

	select {
	case result := <-resCh:
		hostName, _ := os.Hostname()
		for _, v := range result {
			var permalink string
			if v.Permalink == "" {
				permalink = "https://" + hostName + "/detail/" + parser.StringLink(v.Category) + strconv.FormatInt(v.ID, 10) + "/" + parser.StringLink(v.Title)
			} else {
				permalink = v.Permalink
			}

			res.News = append(res.News, &pb.News{
				ID:             parser.IntToStr(v.ID),
				RssID:          parser.IntToStr(v.RssID),
				Title:          v.Title,
				Cover:          v.Cover,
				Author:         v.Author,
				CategorySource: v.CategorySource,
				Category:       v.Category,
				CategoryID:     parser.IntToStr(v.CategoryID),
				PubDate:        v.PubDate,
				Permalink:      permalink,
				IsHeadLine:     parser.IntToStr(v.IsHeadline),
				CreatedAt:      v.CreatedAt,
				Sorting:        parser.IntToStr(v.Sorting),
				Exclusive:      v.Exclusive,
			})
		}

		tr.LogResponse(sp, res)

		return res, nil
	case err = <-errCh:
		tr.LogError(sp, err)
		return nil, err
	}
}

func (*Game) GetTotalNews(ctx context.Context, span opentracing.Span) (res *pb.TotalNews, err error) {
	sp := tr.CreateSubChildSpan(span, "Get Total News Models")
	defer sp.Finish()

	res = &pb.TotalNews{}

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.GetConfig().Timeout)*time.Millisecond)
	defer cancel()

	totalNews := make(chan entity.TotalNews)
	errCh := hystrix.Go("total_news", func() error {
		query := `
			SELECT
				COUNT(*) as total
			FROM
				news_bulk
		`
		tr.LogObject(sp, "query", query)

		var Total entity.TotalNews
		tx := config.GetDBConn().Orm().WithContext(ctx).Debug().
			Raw(query).Model(&entity.TotalNews{}).Scan(&Total)

		if tx.Error != nil {
			return tx.Error
		}

		totalNews <- Total
		tr.LogResponse(sp, Total)

		return nil
	}, nil)

	select {
	case result := <-totalNews:
		res.Total = result.Total
		return res, nil
	case err := <-errCh:
		tr.LogError(sp, err)
		return nil, err
	}
}
