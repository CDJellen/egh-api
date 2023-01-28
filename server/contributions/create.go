package contributions

import (
	"context"

	"github.com/cdjellen/egh-api/app"
	"github.com/cdjellen/egh-api/domain"
	pb "github.com/cdjellen/egh-api/pb/proto"
)

type Create func(context.Context, *pb.CreateContributionsRequest) (*pb.CreateContributionsResponse, error)

func NewCreate(handler app.CreateContributions) Create {
	return func(ctx context.Context, req *pb.CreateContributionsRequest) (*pb.CreateContributionsResponse, error) {

		item := ToDomain(req.GetContributions())

		err := handler(ctx, domain.Login(req.Login), item)
		if err != nil {
			return nil, err
		}
		return &pb.CreateContributionsResponse{}, nil
	}
}
