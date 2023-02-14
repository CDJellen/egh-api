package app

import (
	"context"

	"github.com/cdjellen/egh-api/domain"
)

type ReadContributors func(context.Context, domain.Owner, domain.Repo) (domain.RepoContributors, error)

func NewReadContributors(cache domain.ExploreApi) ReadContributors {
	return func(ctx context.Context, owner domain.Owner, repo domain.Repo) (domain.RepoContributors, error) {

		item, err := cache.ReadContributors(ctx, owner, repo)
		if err != nil {
			return item, err
		}

		if err.Error() == "cache miss" {
			cache.CreateContributors(ctx, owner, repo, item)
		}

		return item, nil
	}
}

type ListContributors func(context.Context) ([]domain.RepoContributors, error)

func NewListContributors(cache domain.ExploreApi) ListContributors {
	return func(ctx context.Context) ([]domain.RepoContributors, error) {
		return cache.ListContributors(ctx)
	}
}

type CreateContributors func(context.Context, domain.Owner, domain.Repo, domain.RepoContributors) error

func NewCreateContributors(cache domain.ExploreApi) CreateContributors {
	return func(ctx context.Context, owner domain.Owner, repo domain.Repo, repoContributors domain.RepoContributors) error {

		err := cache.CreateContributors(ctx, owner, repo, repoContributors)
		if err != nil {
			return err
		}

		return nil
	}
}

type UpdateContributors func(context.Context, domain.Owner, domain.Repo, domain.RepoContributors) error

func NewUpdateContributors(cache domain.ExploreApi) UpdateContributors {
	return func(ctx context.Context, owner domain.Owner, repo domain.Repo, repoContributors domain.RepoContributors) error {

		err := cache.UpdateContributors(ctx, owner, repo, repoContributors)
		if err != nil {
			return err
		}
		return nil
	}
}
