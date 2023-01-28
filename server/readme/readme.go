package readme

import (
	"github.com/cdjellen/egh-api/domain"
	pb "github.com/cdjellen/egh-api/pb/proto"
)

func ToPb(c domain.ReadMe) *pb.ReadMeHtml {
	return &pb.ReadMeHtml{
		Html: c.Html,
	}
}

func ToDomain(c *pb.ReadMeHtml) domain.ReadMe {
	return domain.ReadMe{
		Html: c.GetHtml(),
	}
}
