package app

import (
	"context"

	"github.com/cdjellen/egh-api/domain"
)

type ReadInfo func(context.Context, domain.Owner, domain.Repo) (domain.Contribution, error)

func NewReadInfo(cache domain.ExploreApi) ReadInfo {
	return func(ctx context.Context, owner domain.Owner, repo domain.Repo) (domain.Contribution, error) {

		item, err := cache.ReadInfo(ctx, owner, repo)
		if err != nil {
			return item, err
		}

		if err.Error() == "cache miss" {
			cache.CreateInfo(ctx, owner, repo, item)
		}

		return item, nil
	}
}

type ListInfo func(context.Context) ([]domain.Contribution, error)

func NewListInfo(cache domain.ExploreApi) ListInfo {
	return func(ctx context.Context) ([]domain.Contribution, error) {
		return cache.ListInfo(ctx)
	}
}

type CreateInfo func(context.Context, domain.Owner, domain.Repo, domain.Contribution) error

func NewCreateInfo(cache domain.ExploreApi) CreateInfo {
	return func(ctx context.Context, owner domain.Owner, repo domain.Repo, info domain.Contribution) error {

		err := cache.CreateInfo(ctx, owner, repo, info)
		if err != nil {
			return err
		}

		return nil
	}
}

type UpdateInfo func(context.Context, domain.Owner, domain.Repo, domain.Contribution) error

func NewUpdateInfo(cache domain.ExploreApi) UpdateInfo {
	return func(ctx context.Context, owner domain.Owner, repo domain.Repo, info domain.Contribution) error {

		err := cache.UpdateInfo(ctx, owner, repo, info)
		if err != nil {
			return err
		}
		return nil
	}
}
