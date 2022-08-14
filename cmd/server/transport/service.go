package transport

import (
	"github.com/fajarardiyanto/afaik-svc-server-news/internal/config"
	"github.com/fajarardiyanto/flt-go-listener/lib/proxy"
)

func RunServerGRPCWithProxy() {
	server := NewServerGRPC().CreateServer()

	cc := make(chan bool)
	go proxy.NewProxy(server, config.GetConfig().Server).Start()
	<-cc
}

func RunServerGRPC() {
	NewServerGRPC().CreateServer()
}
