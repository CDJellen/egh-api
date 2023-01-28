package contributions

import (
	"github.com/cdjellen/egh-api/domain"
	pb "github.com/cdjellen/egh-api/pb/proto"
)

func ToPb(c domain.Contributions) *pb.Contributions {
	var s []*pb.RepoContribution

	for i := 0; i < len(c.Contributions); i++ {
		have := c.Contributions[i]

		needOwner := pb.RepoOwner{
			Login:     have.Owner.Login,
			Url:       have.Owner.Url,
			AvatarUrl: have.Owner.AvatarUrl,
		}
		need := pb.RepoContribution{
			NameWithOwner: have.NameWithOwner,
			Name:          have.Name,
			Url:           have.Url,
			Owner:         &needOwner,
		}
		s = append(s, &need)
	}

	return &pb.Contributions{
		Name:          &c.Name,
		Url:           &c.Url,
		AvatarUrl:     &c.AvatarUrl,
		Contributions: s,
	}
}

func ToDomain(c *pb.Contributions) domain.Contributions {
	var s []domain.Contribution
	contributions := c.GetContributions()

	for i := 0; i < len(contributions); i++ {
		have := *contributions[i]
		haveOwner := *have.GetOwner()

		needOwner := domain.RepoOwner{
			Login:     haveOwner.GetLogin(),
			Url:       haveOwner.GetUrl(),
			AvatarUrl: haveOwner.GetAvatarUrl(),
		}
		need := domain.Contribution{
			NameWithOwner: have.GetNameWithOwner(),
			Name:          have.GetName(),
			Url:           have.GetUrl(),
			Owner:         needOwner,
		}

		s = append(s, need)
	}

	return domain.Contributions{
		Name:          c.GetName(),
		Url:           c.GetUrl(),
		AvatarUrl:     c.GetAvatarUrl(),
		Contributions: s,
	}
}
