package readme

import (
	"context"

	"github.com/cdjellen/egh-api/app"
	"github.com/cdjellen/egh-api/domain"
	pb "github.com/cdjellen/egh-api/pb/proto"
)

type Read func(context.Context, *pb.ReadReadMeRequest) (*pb.ReadReadMeResponse, error)

func NewRead(handler app.ReadReadMe) Read {
	return func(ctx context.Context, req *pb.ReadReadMeRequest) (*pb.ReadReadMeResponse, error) {
		item, err := handler(ctx, domain.Owner(req.Owner), domain.Repo(req.Repo), domain.MainBranch(req.MainBranch), domain.FileExt(req.FileExt))
		if err != nil {
			return &pb.ReadReadMeResponse{Message: ToPb(item)}, err
		}

		return &pb.ReadReadMeResponse{Message: ToPb(item)}, nil
	}
}
