package info

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cdjellen/egh-api/app"
	pb "github.com/cdjellen/egh-api/pb/proto"
)

type List func(context.Context, *pb.ListInfoRequest) (*pb.ListInfoResponse, error)

func NewList(handler app.ListInfo) List {
	return func(ctx context.Context, req *pb.ListInfoRequest) (*pb.ListInfoResponse, error) {

		items, err := handler(ctx)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		var r []*pb.RepoContribution

		for _, item := range items {
			r = append(r, ToPb(item))
		}

		return &pb.ListInfoResponse{Messages: r}, nil
	}
}
