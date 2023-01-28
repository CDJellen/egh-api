package readme

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cdjellen/egh-api/app"
	"github.com/cdjellen/egh-api/domain"
	pb "github.com/cdjellen/egh-api/pb/proto"
)

type Update func(context.Context, *pb.UpdateReadMeRequest) (*pb.UpdateReadMeResponse, error)

func NewUpdate(handler app.UpdateReadMe) Update {
	return func(ctx context.Context, req *pb.UpdateReadMeRequest) (*pb.UpdateReadMeResponse, error) {

		err := handler(ctx, domain.Owner(req.Owner), domain.Repo(req.Repo), domain.MainBranch(req.MainBranch), domain.FileExt(req.FileExt), ToDomain(req.GetReadMeHtml()))
		if err != nil {
			return &pb.UpdateReadMeResponse{}, status.Error(codes.Internal, err.Error())
		}

		return &pb.UpdateReadMeResponse{}, nil
	}
}
