package health

import (
	"context"

	"github.com/cdjellen/egh-api/app"
	pb "github.com/cdjellen/egh-api/pb/proto"
)

type Read func(context.Context, *pb.HealthRequest) (*pb.HealthResponse, error)

func NewRead(handler app.ReadHealth) Read {
	return func(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {

		item, err := handler(ctx)
		if err == nil {
			return ToPb(item), nil
		}

		return &pb.HealthResponse{Message: "dead"}, err
	}
}
