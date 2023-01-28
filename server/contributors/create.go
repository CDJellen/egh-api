package contributors

import (
	"context"

	"github.com/cdjellen/egh-api/app"
	"github.com/cdjellen/egh-api/domain"
	pb "github.com/cdjellen/egh-api/pb/proto"
)

type Create func(context.Context, *pb.CreateContributorsRequest) (*pb.CreateContributorsResponse, error)

func NewCreate(handler app.CreateContributors) Create {
	return func(ctx context.Context, req *pb.CreateContributorsRequest) (*pb.CreateContributorsResponse, error) {

		item := ToDomain(req.GetRepoContributors())

		err := handler(ctx, domain.Owner(req.Owner), domain.Repo(req.Repo), item)
		if err != nil {
			return nil, err
		}
		return &pb.CreateContributorsResponse{}, nil
	}
}
