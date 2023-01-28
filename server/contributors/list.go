package contributors

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cdjellen/egh-api/app"
	pb "github.com/cdjellen/egh-api/pb/proto"
)

type List func(context.Context, *pb.ListContributorsRequest) (*pb.ListContributorsResponse, error)

func NewList(handler app.ListContributors) List {
	return func(ctx context.Context, req *pb.ListContributorsRequest) (*pb.ListContributorsResponse, error) {

		items, err := handler(ctx)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		var r []*pb.RepoContributors

		for _, item := range items {
			r = append(r, ToPb(item))
		}

		return &pb.ListContributorsResponse{Messages: r}, nil
	}
}
