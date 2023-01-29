package health

import (
	"github.com/cdjellen/egh-api/domain"
	pb "github.com/cdjellen/egh-api/pb/proto"
)

func ToPb(c domain.Health) *pb.HealthResponse {

	return &pb.HealthResponse{
		Message: c.Status,
	}
}

func ToDomain(c *pb.HealthResponse) domain.Health {

	return domain.Health{
		Status: c.GetMessage(),
	}
}
