package info

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cdjellen/egh-api/app"
	"github.com/cdjellen/egh-api/domain"
	pb "github.com/cdjellen/egh-api/pb/proto"
)

type Update func(context.Context, *pb.UpdateInfoRequest) (*pb.UpdateInfoResponse, error)

func NewUpdate(handler app.UpdateInfo) Update {
	return func(ctx context.Context, req *pb.UpdateInfoRequest) (*pb.UpdateInfoResponse, error) {

		err := handler(ctx, domain.Owner(req.Owner), domain.Repo(req.Repo), ToDomain(req.GetInfo()))
		if err != nil {
			return &pb.UpdateInfoResponse{}, status.Error(codes.Internal, err.Error())
		}

		return &pb.UpdateInfoResponse{}, nil
	}
}
