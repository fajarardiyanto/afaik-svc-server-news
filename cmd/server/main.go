package main

import (
	"github.com/fajarardiyanto/prometheus-svc-server-news/cmd/server/transport"
)

func main() {
	transport.RunServerGRPCWithProxy()
}
