package transport

import (
	"github.com/fajarardiyanto/flt-go-listener/lib/server"
	pb "github.com/fajarardiyanto/module-proto/go/services/news"
	"github.com/fajarardiyanto/prometheus-svc-server-news/app/models"
	"github.com/fajarardiyanto/prometheus-svc-server-news/app/service"
	"github.com/fajarardiyanto/prometheus-svc-server-news/internal/config"
	"google.golang.org/grpc"
)

type ServerGRPC struct {
	server *grpc.Server
}

func NewServerGRPC() *ServerGRPC {
	return &ServerGRPC{}
}

func (*ServerGRPC) CreateServer() *grpc.Server {
	config.Init()

	//Repository
	repoGround := models.NewGame()

	//Service
	ground := service.NewGameService(repoGround)

	//Register GRPC Server
	serv := server.NewListenerServer(config.GetLogger(), config.GetConfig().Server)
	if coreServ := serv.GetGRPCServer(); coreServ != nil {
		pb.RegisterNewsServiceServer(coreServ, ground)
		serv.Init()

		return coreServ
	}

	return nil
}
