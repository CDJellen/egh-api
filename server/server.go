package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/rs/cors"

	"github.com/cdjellen/egh-api/app"
	pb "github.com/cdjellen/egh-api/pb/proto"
	"github.com/cdjellen/egh-api/server/contributions"
	"github.com/cdjellen/egh-api/server/contributors"
	"github.com/cdjellen/egh-api/server/health"
	"github.com/cdjellen/egh-api/server/info"
	"github.com/cdjellen/egh-api/server/readme"
	"github.com/cdjellen/egh-api/store"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type HealthServer struct {
	HealthEndpoint health.Read
}

type InfoServer struct {
	CreateInfoEndpoint info.Create
	ReadInfoEndpoint   info.Read
	UpdateInfoEndpoint info.Update
	ListInfoEndpoint   info.List
}

type ReadMeServer struct {
	CreateReadMeEndpoint readme.Create
	ReadReadMeEndpoint   readme.Read
	UpdateReadMeEndpoint readme.Update
	ListReadMeEndpoint   readme.List
}

type ContributionsServer struct {
	CreateContributionsEndpoint contributions.Create
	ReadContributionsEndpoint   contributions.Read
	UpdateContributionsEndpoint contributions.Update
	ListContributionsEndpoint   contributions.List
}

type ContributorsServer struct {
	ReadContributorsEndpoint   contributors.Read
	ListContributorsEndpoint   contributors.List
	CreateContributorsEndpoint contributors.Create
	UpdateContributorsEndpoint contributors.Update
}

func NewHealthServer() *HealthServer {
	s := HealthServer{
		HealthEndpoint: health.NewRead(app.NewReadHealth()),
	}

	return &s
}

func (s *HealthServer) ReadHealth(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
	return s.HealthEndpoint(ctx, req)
}

func NewInfoServer(cache *store.ExploreApiCache) *InfoServer {
	s := InfoServer{
		ReadInfoEndpoint:   info.NewRead(app.NewReadInfo(cache)),
		ListInfoEndpoint:   info.NewList(app.NewListInfo(cache)),
		CreateInfoEndpoint: info.NewCreate(app.NewCreateInfo(cache)),
		UpdateInfoEndpoint: info.NewUpdate(app.NewUpdateInfo(cache)),
	}

	return &s
}

func (s *InfoServer) ReadInfo(ctx context.Context, req *pb.ReadInfoRequest) (*pb.ReadInfoResponse, error) {
	return s.ReadInfoEndpoint(ctx, req)
}

func (s *InfoServer) ListInfo(ctx context.Context, req *pb.ListInfoRequest) (*pb.ListInfoResponse, error) {
	return s.ListInfoEndpoint(ctx, req)
}

func (s *InfoServer) CreateInfo(ctx context.Context, req *pb.CreateInfoRequest) (*pb.CreateInfoResponse, error) {
	return s.CreateInfoEndpoint(ctx, req)
}

func (s *InfoServer) UpdateInfo(ctx context.Context, req *pb.UpdateInfoRequest) (*pb.UpdateInfoResponse, error) {
	return s.UpdateInfoEndpoint(ctx, req)
}

func NewReadMeServer(cache *store.ExploreApiCache) *ReadMeServer {
	s := ReadMeServer{
		ReadReadMeEndpoint:   readme.NewRead(app.NewReadReadMe(cache)),
		ListReadMeEndpoint:   readme.NewList(app.NewListReadMe(cache)),
		CreateReadMeEndpoint: readme.NewCreate(app.NewCreateReadMe(cache)),
		UpdateReadMeEndpoint: readme.NewUpdate(app.NewUpdateReadMe(cache)),
	}

	return &s
}

func (s *ReadMeServer) ReadReadMe(ctx context.Context, req *pb.ReadReadMeRequest) (*pb.ReadReadMeResponse, error) {
	return s.ReadReadMeEndpoint(ctx, req)
}

func (s *ReadMeServer) ListReadMe(ctx context.Context, req *pb.ListReadMeRequest) (*pb.ListReadMeResponse, error) {
	return s.ListReadMeEndpoint(ctx, req)
}

func (s *ReadMeServer) CreateReadMe(ctx context.Context, req *pb.CreateReadMeRequest) (*pb.CreateReadMeResponse, error) {
	return s.CreateReadMeEndpoint(ctx, req)
}

func (s *ReadMeServer) UpdateReadMe(ctx context.Context, req *pb.UpdateReadMeRequest) (*pb.UpdateReadMeResponse, error) {
	return s.UpdateReadMeEndpoint(ctx, req)
}

func NewContributionsServer(cache *store.ExploreApiCache) *ContributionsServer {
	s := ContributionsServer{
		ReadContributionsEndpoint:   contributions.NewRead(app.NewReadContributions(cache)),
		ListContributionsEndpoint:   contributions.NewList(app.NewListContributions(cache)),
		CreateContributionsEndpoint: contributions.NewCreate(app.NewCreateContributions(cache)),
		UpdateContributionsEndpoint: contributions.NewUpdate(app.NewUpdateContributions(cache)),
	}

	return &s
}

func (s *ContributionsServer) ReadContributions(ctx context.Context, req *pb.ReadContributionsRequest) (*pb.ReadContributionsResponse, error) {
	return s.ReadContributionsEndpoint(ctx, req)
}

func (s *ContributionsServer) ListContributions(ctx context.Context, req *pb.ListContributionsRequest) (*pb.ListContributionsResponse, error) {
	return s.ListContributionsEndpoint(ctx, req)
}

func (s *ContributionsServer) CreateContributions(ctx context.Context, req *pb.CreateContributionsRequest) (*pb.CreateContributionsResponse, error) {
	return s.CreateContributionsEndpoint(ctx, req)
}

func (s *ContributionsServer) UpdateContributions(ctx context.Context, req *pb.UpdateContributionsRequest) (*pb.UpdateContributionsResponse, error) {
	return s.UpdateContributionsEndpoint(ctx, req)
}

func NewContributorsServer(cache *store.ExploreApiCache) *ContributorsServer {
	s := ContributorsServer{
		ReadContributorsEndpoint:   contributors.NewRead(app.NewReadContributors(cache)),
		ListContributorsEndpoint:   contributors.NewList(app.NewListContributors(cache)),
		CreateContributorsEndpoint: contributors.NewCreate(app.NewCreateContributors(cache)),
		UpdateContributorsEndpoint: contributors.NewUpdate(app.NewUpdateContributors(cache)),
	}

	return &s
}

func (s *ContributorsServer) ReadContributors(ctx context.Context, req *pb.ReadContributorsRequest) (*pb.ReadContributorsResponse, error) {
	return s.ReadContributorsEndpoint(ctx, req)
}

func (s *ContributorsServer) ListContributors(ctx context.Context, req *pb.ListContributorsRequest) (*pb.ListContributorsResponse, error) {
	return s.ListContributorsEndpoint(ctx, req)
}

func (s *ContributorsServer) CreateContributors(ctx context.Context, req *pb.CreateContributorsRequest) (*pb.CreateContributorsResponse, error) {
	return s.CreateContributorsEndpoint(ctx, req)
}

func (s *ContributorsServer) UpdateContributors(ctx context.Context, req *pb.UpdateContributorsRequest) (*pb.UpdateContributorsResponse, error) {
	return s.UpdateContributorsEndpoint(ctx, req)
}

func Run(ctx context.Context, network string, address string, name string, healthServer *HealthServer, infoServer *InfoServer, contributionsServer *ContributionsServer, contributorsServer *ContributorsServer, readMeServer *ReadMeServer) error {
	l, err := net.Listen(network, address)
	if err != nil {
		return err
	}

	defer func() {
		if err := l.Close(); err != nil {
			fmt.Printf("server failed to run with error: %+v", err)
		}
	}()

	s := grpc.NewServer()
	pb.RegisterHealthServiceServer(s, healthServer)
	pb.RegisterInfoServiceServer(s, infoServer)
	pb.RegisterContributionsServiceServer(s, contributionsServer)
	pb.RegisterContributorsServiceServer(s, contributorsServer)
	pb.RegisterReadMeServiceServer(s, readMeServer)

	go func() {
		defer s.GracefulStop()
		<-ctx.Done()
	}()
	return s.Serve(l)
}

func RunInProcessGateway(ctx context.Context, addr string, name string, healthServer *HealthServer, infoServer *InfoServer, contributionsServer *ContributionsServer, contributorsServer *ContributorsServer, readMeServer *ReadMeServer, opts ...runtime.ServeMuxOption) error {
	gw := runtime.NewServeMux(opts...)
	docs := http.StripPrefix("/api/docs/", http.FileServer(http.Dir("./swagger/")))
	frontend := http.StripPrefix("/", http.FileServer(http.Dir("./frontend/")))

	err := pb.RegisterHealthServiceHandlerServer(ctx, gw, healthServer)
	if err != nil {
		panic(err)
	}
	err = pb.RegisterInfoServiceHandlerServer(ctx, gw, infoServer)
	if err != nil {
		panic(err)
	}
	err = pb.RegisterContributionsServiceHandlerServer(ctx, gw, contributionsServer)
	if err != nil {
		panic(err)
	}
	err = pb.RegisterContributorsServiceHandlerServer(ctx, gw, contributorsServer)
	if err != nil {
		panic(err)
	}
	err = pb.RegisterReadMeServiceHandlerServer(ctx, gw, readMeServer)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", frontend)
	mux.Handle("/api/v1/", gw)
	mux.Handle("/api/docs/", docs)

	// CORS
	withCors := cors.New(cors.Options{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"ACCEPT", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler(mux)

	h := &http.Server{
		Addr:              addr,
		Handler:           withCors,
		ReadHeaderTimeout: 15 * time.Second,
	}

	go func() {
		<-ctx.Done()
		if err := h.Shutdown(context.Background()); err != nil {
			fmt.Printf("server failed to run with error: %+v", err)
		}
	}()

	if err := h.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}
