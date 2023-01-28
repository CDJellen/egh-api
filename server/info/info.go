package info

import (
	"github.com/cdjellen/egh-api/domain"
	pb "github.com/cdjellen/egh-api/pb/proto"
)

func ToPb(c domain.Contribution) *pb.RepoContribution {
	o := &pb.RepoOwner{
		Login:     c.Owner.Login,
		Url:       c.Owner.Url,
		AvatarUrl: c.Owner.AvatarUrl,
	}

	return &pb.RepoContribution{
		NameWithOwner: c.NameWithOwner,
		Name:          c.Name,
		Url:           c.Url,
		Owner:         o,
	}
}

func ToDomain(c *pb.RepoContribution) domain.Contribution {
	o := c.GetOwner()

	return domain.Contribution{
		NameWithOwner: c.GetNameWithOwner(),
		Name:          c.GetName(),
		Url:           c.GetUrl(),
		Owner: domain.RepoOwner{
			Login:     o.GetLogin(),
			Url:       o.GetUrl(),
			AvatarUrl: o.GetAvatarUrl(),
		},
	}
}
