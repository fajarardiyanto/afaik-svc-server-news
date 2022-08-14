package transport

import (
	"github.com/fajarardiyanto/flt-go-listener/lib/proxy"
	"github.com/fajarardiyanto/prometheus-svc-server-news/internal/config"
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
