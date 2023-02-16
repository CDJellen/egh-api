package contributions

import (
	"context"

	"github.com/cdjellen/egh-api/app"
	"github.com/cdjellen/egh-api/domain"
	pb "github.com/cdjellen/egh-api/pb/proto"
)

type Read func(context.Context, *pb.ReadContributionsRequest) (*pb.ReadContributionsResponse, error)

func NewRead(handler app.ReadContributions) Read {
	return func(ctx context.Context, req *pb.ReadContributionsRequest) (*pb.ReadContributionsResponse, error) {
		item, err := handler(ctx, domain.Login(req.GetLogin()), req.GetNumContributions())
		if err != nil {
			return &pb.ReadContributionsResponse{Message: ToPb(item)}, err
		}

		return &pb.ReadContributionsResponse{Message: ToPb(item)}, nil
	}
}
