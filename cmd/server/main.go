package main

import (
	"github.com/fajarardiyanto/afaik-svc-server-news/cmd/server/transport"
)

func main() {
	transport.RunServerGRPCWithProxy()
}
