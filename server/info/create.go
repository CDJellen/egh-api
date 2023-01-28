package info

import (
	"context"

	"github.com/cdjellen/egh-api/app"
	"github.com/cdjellen/egh-api/domain"
	pb "github.com/cdjellen/egh-api/pb/proto"
)

type Create func(context.Context, *pb.CreateInfoRequest) (*pb.CreateInfoResponse, error)

func NewCreate(handler app.CreateInfo) Create {
	return func(ctx context.Context, req *pb.CreateInfoRequest) (*pb.CreateInfoResponse, error) {

		item := ToDomain(req.GetInfo())

		err := handler(ctx, domain.Owner(req.Owner), domain.Repo(req.Repo), item)
		if err != nil {
			return nil, err
		}
		return &pb.CreateInfoResponse{}, nil
	}
}
