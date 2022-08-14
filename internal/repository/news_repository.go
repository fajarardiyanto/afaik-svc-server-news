package repository

import (
	"context"
	pb "github.com/fajarardiyanto/module-proto/go/services/news"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/protobuf/types/known/emptypb"
)

type NewsRepository interface {
	GetTotalNews(context.Context, opentracing.Span) (*pb.TotalNews, error)
	Get(context.Context, *pb.GetNewsRequest, opentracing.Span) (*pb.NewsResponse, error)
	Store(context.Context, *pb.NewsRequest, opentracing.Span) (*emptypb.Empty, error)
}
