package readme

import (
	"context"

	"github.com/cdjellen/egh-api/app"
	"github.com/cdjellen/egh-api/domain"
	pb "github.com/cdjellen/egh-api/pb/proto"
)

type Create func(context.Context, *pb.CreateReadMeRequest) (*pb.CreateReadMeResponse, error)

func NewCreate(handler app.CreateReadMe) Create {
	return func(ctx context.Context, req *pb.CreateReadMeRequest) (*pb.CreateReadMeResponse, error) {

		item := ToDomain(req.GetReadMeHtml())

		err := handler(ctx, domain.Owner(req.Owner), domain.Repo(req.Repo), domain.MainBranch(req.MainBranch), domain.FileExt(req.FileExt), item)
		if err != nil {
			return nil, err
		}
		return &pb.CreateReadMeResponse{}, nil
	}
}
