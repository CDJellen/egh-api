package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/cdjellen/egh-api/server"
	"github.com/cdjellen/egh-api/store/mem"
	"github.com/cdjellen/egh-api/store/redis"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

var (
	name     = flag.String("name", "egh-api", "Server name for logging and tracing")
	protocol = flag.String("protocol", "tcp", "protocol type")
	rpc      = flag.String("rpc", ":50051", "gRPC server endpoint")
	gw       = flag.String("gw", ":8080", "REST gateway endpoint")
	store    = flag.String("store", "redis", "backend cache for remote requests")
	rds      = flag.String("redis", "172.17.0.2:6379", "address for optional redis cluster") //localhost:6379
	user     = flag.String("user", "", "optional redis username")
	pass     = flag.String("pass", "", "optional redis password")
	db       = flag.Int("db", 0, "optional redis DB index")
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

	var healthServer *server.HealthServer
	var infoServer *server.InfoServer
	var contributionsServer *server.ContributionsServer
	var contributorsServer *server.ContributorsServer
	var readMeServer *server.ReadMeServer

	if *store == "redis" {
		cache := redis.NewRedisCache(*rds, *user, *pass, *db)
		healthServer = server.NewHealthServer()
		infoServer = server.NewInfoServer(cache)
		contributionsServer = server.NewContributionsServer(cache)
		contributorsServer = server.NewContributorsServer(cache)
		readMeServer = server.NewReadMeServer(cache)
	} else {
		cache := mem.NewExploreApiCache()
		healthServer = server.NewHealthServer()
		infoServer = server.NewInfoServer(cache)
		contributionsServer = server.NewContributionsServer(cache)
		contributorsServer = server.NewContributorsServer(cache)
		readMeServer = server.NewReadMeServer(cache)
	}

	go func() {
		if err := server.Run(ctx, *protocol, *rpc, *name, healthServer, infoServer, contributionsServer, contributorsServer, readMeServer); err != nil {
			return
		}
	}()

	go func() {
		opts := []runtime.ServeMuxOption{}
		if err := server.RunInProcessGateway(ctx, *gw, *name, healthServer, infoServer, contributionsServer, contributorsServer, readMeServer, opts...); err != nil {
			return
		}
	}()

	select {
	case <-quit:
		log.Println("shutting down server")
		cancel()
	case <-ctx.Done():
		log.Println("shutting down server")
		cancel()
	}
}
