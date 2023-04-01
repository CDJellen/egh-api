package app

import (
	"context"
	"log"
	"strings"

	"github.com/machinebox/graphql"

	"github.com/cdjellen/egh-api/app/remote"
	"github.com/cdjellen/egh-api/domain"
)

type ReadContributions func(context.Context, domain.Login, int32) (domain.Contributions, error)

func NewReadContributions(cache domain.ExploreApi) ReadContributions {
	return func(ctx context.Context, login domain.Login, numContributions int32) (domain.Contributions, error) {

		// check the cache
		item, err := cache.ReadContributions(ctx, login, numContributions) // TODO get arbitrary number of contributions from cache
		if err == nil {
			return item, nil
		}
		if strings.Contains(err.Error(), "cache miss") {
			// read from remote
			item, err = contributionsRequest(ctx, login, 100)
			if err != nil {
				log.Printf("failed to get CONTRIBUTIONS with error %+v", err)
				return item, err
			}

			// persist to cache
			err = cache.CreateContributions(ctx, login, item)
			if err != nil {
				return item, err
			}
		}
		// pull from cache with params
		item, err = cache.ReadContributions(ctx, login, numContributions)
		if err == nil {
			return item, nil
		}

		return item, err
	}
}

func contributionsRequest(ctx context.Context, login domain.Login, last int32) (domain.Contributions, error) {

	client := remote.GetGraphQLClient()
	headers := remote.GetGitHubHeaders()

	// make a request
	req := graphql.NewRequest(remote.GqlContributionsQuery)
	req.Header = *headers
	req.Var("key", string(login))
	req.Var("last", last)

	resp := remote.GqlResponse{}
	if err := client.Run(ctx, req, &resp); err != nil {
		log.Printf("Failed to unpack request with error %+v", err)
		return domain.Contributions{}, err
	}

	reposContributedTo := []domain.Contribution{}

	for _, c := range resp.User.RepositoriesContributedTo.Edges {
		newOwner := domain.RepoOwner{
			Login:     c.Node.Owner.Login,
			Url:       c.Node.Owner.Url,
			AvatarUrl: c.Node.Owner.AvatarUrl,
		}
		newContrib := domain.Contribution{
			NameWithOwner: c.Node.NameWithOwner,
			Name:          c.Node.Name,
			Url:           c.Node.Url,
			Owner:         newOwner,
		}

		reposContributedTo = append(reposContributedTo, newContrib)
	}

	contributions := domain.Contributions{
		Name:          resp.User.Name,
		Url:           resp.User.Url,
		AvatarUrl:     resp.User.AvatarUrl,
		Contributions: reposContributedTo,
	}

	return contributions, nil
}

type ListContributions func(context.Context) ([]domain.Contributions, error)

func NewListContributions(cache domain.ExploreApi) ListContributions {
	return func(ctx context.Context) ([]domain.Contributions, error) {
		return cache.ListContributions(ctx)
	}
}

type CreateContributions func(context.Context, domain.Login, domain.Contributions) error

func NewCreateContributions(cache domain.ExploreApi) CreateContributions {
	return func(ctx context.Context, login domain.Login, contributions domain.Contributions) error {

		err := cache.CreateContributions(ctx, login, contributions)
		if err != nil {
			return err
		}

		return nil
	}
}

type UpdateContributions func(context.Context, domain.Login, domain.Contributions) error

func NewUpdateContributions(cache domain.ExploreApi) UpdateContributions {
	return func(ctx context.Context, login domain.Login, contributions domain.Contributions) error {

		err := cache.UpdateContributions(ctx, login, contributions)
		if err != nil {
			return err
		}
		return nil
	}
}
