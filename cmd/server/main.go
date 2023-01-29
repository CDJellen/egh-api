package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/cdjellen/egh-api/server"
	"github.com/cdjellen/egh-api/store"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

var (
	rpc      = flag.String("rpc", ":5000", "gRPC server endpoint")
	protocol = flag.String("protocol", "tcp", "protocol type")
	gw       = flag.String("gw", ":8080", "REST gateway endpoint")
)

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	defer func() {
		signal.Stop(quit)
		cancel()
	}()

	cache := store.NewExploreApiCache()
	infoServer := server.NewInfoServer(cache)
	contributionsServer := server.NewContributionsServer(cache)
	contributorsServer := server.NewContributorsServer(cache)
	readMeServer := server.NewReadMeServer(cache)

	go func() {
		if err := server.Run(ctx, *protocol, *rpc, infoServer, contributionsServer, contributorsServer, readMeServer); err != nil {
			return
		}
	}()

	go func() {
		opts := []runtime.ServeMuxOption{}
		if err := server.RunInProcessGateway(ctx, *gw, infoServer, contributionsServer, contributorsServer, readMeServer, opts...); err != nil {
			return
		}
	}()

	select {
	case <-quit:
		fmt.Println("shutting down server")
		cancel()
	case <-ctx.Done():
		fmt.Println("shutting down server")
		cancel()
	}
}
