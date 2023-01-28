package contributions

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cdjellen/egh-api/app"
	"github.com/cdjellen/egh-api/domain"
	pb "github.com/cdjellen/egh-api/pb/proto"
)

type Update func(context.Context, *pb.UpdateContributionsRequest) (*pb.UpdateContributionsResponse, error)

func NewUpdate(handler app.UpdateContributions) Update {
	return func(ctx context.Context, req *pb.UpdateContributionsRequest) (*pb.UpdateContributionsResponse, error) {

		err := handler(ctx, domain.Login(req.Login), ToDomain(req.GetContributions()))
		if err != nil {
			return &pb.UpdateContributionsResponse{}, status.Error(codes.Internal, err.Error())
		}

		return &pb.UpdateContributionsResponse{}, nil
	}
}
