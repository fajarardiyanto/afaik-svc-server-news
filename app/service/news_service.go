package service

import (
	"context"
	"github.com/fajarardiyanto/afaik-svc-server-news/internal/config"
	"github.com/fajarardiyanto/afaik-svc-server-news/internal/repository"
	"github.com/fajarardiyanto/flt-go-tracer/lib/jaeger"
	"github.com/fajarardiyanto/flt-go-utils/validation"
	pb "github.com/fajarardiyanto/module-proto/go/services/news"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync"
)

type GameService struct {
	validator *validator.Validate
	model     repository.NewsRepository
	pb.UnimplementedNewsServiceServer
	sync.RWMutex
}

func NewGameService(model repository.NewsRepository) pb.NewsServiceServer {
	return &GameService{
		model:     model,
		validator: utils.GetValidator(),
	}
}

func (c *GameService) GetNews(ctx context.Context, req *pb.GetNewsRequest) (res *pb.NewsResponse, err error) {
	span, _ := jaeger.RootSpan(ctx)
	defer span.Finish()

	res = &pb.NewsResponse{}

	if err = c.validator.Struct(req); err != nil {
		err = utils.ParseValidatorError(err)
		config.GetLogger().Error(err)
		return res, status.Error(codes.FailedPrecondition, err.Error())
	}

	return c.model.Get(ctx, req, span)
}

func (c *GameService) CreateNews(ctx context.Context, req *pb.NewsRequest) (res *emptypb.Empty, err error) {
	span, _ := jaeger.RootSpan(ctx)
	defer span.Finish()

	if err = c.validator.Struct(req); err != nil {
		err = utils.ParseValidatorError(err)
		config.GetLogger().Error(err)
		return res, status.Error(codes.FailedPrecondition, err.Error())
	}

	return c.model.Store(ctx, req, span)
}

func (c *GameService) GetTotalNews(ctx context.Context, req *emptypb.Empty) (res *pb.TotalNews, err error) {
	span, _ := jaeger.RootSpan(ctx)
	defer span.Finish()

	if err = c.validator.Struct(req); err != nil {
		err = utils.ParseValidatorError(err)
		config.GetLogger().Error(err)
		return res, status.Error(codes.FailedPrecondition, err.Error())
	}

	return c.model.GetTotalNews(ctx, span)
}
