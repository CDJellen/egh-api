package info

import (
	"context"

	"github.com/cdjellen/egh-api/app"
	"github.com/cdjellen/egh-api/domain"
	pb "github.com/cdjellen/egh-api/pb/proto"
)

type Read func(context.Context, *pb.ReadInfoRequest) (*pb.ReadInfoResponse, error)

func NewRead(handler app.ReadInfo) Read {
	return func(ctx context.Context, req *pb.ReadInfoRequest) (*pb.ReadInfoResponse, error) {
		item, err := handler(ctx, domain.Owner(req.Owner), domain.Repo(req.Repo))
		if err != nil {
			return &pb.ReadInfoResponse{Message: ToPb(item)}, err
		}

		return &pb.ReadInfoResponse{Message: ToPb(item)}, nil
	}
}
