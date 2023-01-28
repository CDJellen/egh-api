package readme

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cdjellen/egh-api/app"
	pb "github.com/cdjellen/egh-api/pb/proto"
)

type List func(context.Context, *pb.ListReadMeRequest) (*pb.ListReadMeResponse, error)

func NewList(handler app.ListReadMe) List {
	return func(ctx context.Context, req *pb.ListReadMeRequest) (*pb.ListReadMeResponse, error) {

		items, err := handler(ctx)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		var r []*pb.ReadMeHtml

		for _, item := range items {
			r = append(r, ToPb(item))
		}

		return &pb.ListReadMeResponse{Messages: r}, nil
	}
}
