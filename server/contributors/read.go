package contributors

import (
	"context"

	"github.com/cdjellen/egh-api/app"
	"github.com/cdjellen/egh-api/domain"
	pb "github.com/cdjellen/egh-api/pb/proto"
)

type Read func(context.Context, *pb.ReadContributorsRequest) (*pb.ReadContributorsResponse, error)

func NewRead(handler app.ReadContributors) Read {
	return func(ctx context.Context, req *pb.ReadContributorsRequest) (*pb.ReadContributorsResponse, error) {
		item, err := handler(ctx, domain.Owner(req.GetOwner()), domain.Repo(req.GetRepo()), req.GetAnon(), req.GetPerPage(), req.GetPage())
		if err != nil {
			return &pb.ReadContributorsResponse{Message: ToPb(item)}, err
		}

		return &pb.ReadContributorsResponse{Message: ToPb(item)}, nil
	}
}
