package contributors

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cdjellen/egh-api/app"
	"github.com/cdjellen/egh-api/domain"
	pb "github.com/cdjellen/egh-api/pb/proto"
)

type Update func(context.Context, *pb.UpdateContributorsRequest) (*pb.UpdateContributorsResponse, error)

func NewUpdate(handler app.UpdateContributors) Update {
	return func(ctx context.Context, req *pb.UpdateContributorsRequest) (*pb.UpdateContributorsResponse, error) {

		err := handler(ctx, domain.Owner(req.Owner), domain.Repo(req.Repo), ToDomain(req.GetRepoContributors()))
		if err != nil {
			return &pb.UpdateContributorsResponse{}, status.Error(codes.Internal, err.Error())
		}

		return &pb.UpdateContributorsResponse{}, nil
	}
}
