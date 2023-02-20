package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/cdjellen/egh-api/server"
	"github.com/cdjellen/egh-api/server/telemetry/tracing"
	"github.com/cdjellen/egh-api/store/mem"
	"github.com/cdjellen/egh-api/store/redis"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

var (
	name     = flag.String("name", "egh-api", "Server name for logging and tracing")
	protocol = flag.String("protocol", "tcp", "protocol type")
	rpc      = flag.String("rpc", ":50051", "gRPC server endpoint")
	gw       = flag.String("gw", ":8080", "REST gateway endpoint")
	trace    = flag.String("tp", ":4317", "OpenTelemetry tracing endpoint")
	store    = flag.String("store", "redis", "backend cache for remote requests")
	rds      = flag.String("redis", "localhost:6379", "address for optional redis cluster")
	user     = flag.String("user", "", "optional redis username")
	pass     = flag.String("pass", "", "optional redis password")
	db       = flag.Int("db", 0, "optional redis DB index")
	sampling = flag.Float64("sampling", 1.0, "tracing sampling ratio")
)

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	{
		provider, err := tracing.NewGrpcTraceProvider(ctx, *name, *trace, *sampling)
		if err != nil {
			log.Printf("failed to stand up new otel tracing provider at address %s with sampling ratio %v", *trace, *sampling)
		}
		defer provider.Shutdown(ctx)
	}
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
