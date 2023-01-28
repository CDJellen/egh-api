package contributors

import (
	"github.com/cdjellen/egh-api/domain"
	pb "github.com/cdjellen/egh-api/pb/proto"
)

func ToPb(c domain.RepoContributors) *pb.RepoContributors {
	var s []*pb.RepoContributor

	for i := 0; i < len(c.RepoContributors); i++ {
		have := c.RepoContributors[i]
		need := pb.RepoContributor{
			Login:             have.Login,
			Id:                have.ID,
			NodeId:            have.NodeID,
			AvatarUrl:         have.AvatarUrl,
			GravatarId:        &have.GravatarID,
			Url:               have.Url,
			HtmlUrl:           have.HtmlUrl,
			FollowersUrl:      have.FollowersUrl,
			FollowingUrl:      have.FollowingUrl,
			GistsUrl:          have.GistsUrl,
			StarredUrl:        have.StarredUrl,
			SubscriptionsUrl:  have.SubscriptionsUrl,
			OrganizationsUrl:  have.OrganizationsUrl,
			ReposUrl:          have.ReposUrl,
			EventsUrl:         have.EventsUrl,
			ReceivedEventsUrl: have.ReceivedEventsUrl,
			Type:              have.Type,
			SiteAdmin:         have.SiteAdmin,
			Contributions:     have.Contributions,
		}

		s = append(s, &need)
	}

	return &pb.RepoContributors{
		Contributors: s,
	}
}

func ToDomain(c *pb.RepoContributors) domain.RepoContributors {
	var s []domain.Contributors
	contributors := c.GetContributors()

	for i := 0; i < len(contributors); i++ {
		have := *contributors[i]
		need := domain.Contributors{
			Login:             have.GetLogin(),
			ID:                have.GetId(),
			NodeID:            have.GetNodeId(),
			AvatarUrl:         have.GetAvatarUrl(),
			GravatarID:        have.GetGravatarId(),
			Url:               have.GetUrl(),
			HtmlUrl:           have.GetHtmlUrl(),
			FollowersUrl:      have.GetFollowersUrl(),
			FollowingUrl:      have.GetFollowingUrl(),
			GistsUrl:          have.GetGistsUrl(),
			StarredUrl:        have.GetStarredUrl(),
			SubscriptionsUrl:  have.GetSubscriptionsUrl(),
			OrganizationsUrl:  have.GetOrganizationsUrl(),
			ReposUrl:          have.GetReposUrl(),
			EventsUrl:         have.GetEventsUrl(),
			ReceivedEventsUrl: have.GetReceivedEventsUrl(),
			Type:              have.GetType(),
			SiteAdmin:         have.GetSiteAdmin(),
			Contributions:     have.GetContributions(),
		}

		s = append(s, need)
	}

	return domain.RepoContributors{
		RepoContributors: s,
	}
}
