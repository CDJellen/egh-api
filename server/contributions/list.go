package contributions

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cdjellen/egh-api/app"
	pb "github.com/cdjellen/egh-api/pb/proto"
)

type List func(context.Context, *pb.ListContributionsRequest) (*pb.ListContributionsResponse, error)

func NewList(handler app.ListContributions) List {
	return func(ctx context.Context, req *pb.ListContributionsRequest) (*pb.ListContributionsResponse, error) {

		items, err := handler(ctx)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		var r []*pb.Contributions

		for _, item := range items {
			r = append(r, ToPb(item))
		}

		return &pb.ListContributionsResponse{Messages: r}, nil
	}
}
