package app

import (
	"context"
	"strings"

	"github.com/cdjellen/egh-api/domain"
)

type ReadReadMe func(context.Context, domain.Owner, domain.Repo, domain.MainBranch, domain.FileExt) (domain.ReadMe, error)

func NewReadReadMe(cache domain.ExploreApi) ReadReadMe {
	return func(ctx context.Context, owner domain.Owner, repo domain.Repo, main domain.MainBranch, ext domain.FileExt) (domain.ReadMe, error) {

		item, err := cache.ReadReadMe(ctx, owner, repo, main, ext)
		if err == nil {
			return item, nil
		}
		if strings.Contains(err.Error(), "cache miss") {
			err = cache.CreateReadMe(ctx, owner, repo, main, ext, item)
			if err != nil {
				return item, err
			}

			return item, nil

		}

		return item, err
	}
}

type ListReadMe func(context.Context) ([]domain.ReadMe, error)

func NewListReadMe(cache domain.ExploreApi) ListReadMe {
	return func(ctx context.Context) ([]domain.ReadMe, error) {
		return cache.ListReadMe(ctx)
	}
}

type CreateReadMe func(context.Context, domain.Owner, domain.Repo, domain.MainBranch, domain.FileExt, domain.ReadMe) error

func NewCreateReadMe(cache domain.ExploreApi) CreateReadMe {
	return func(ctx context.Context, owner domain.Owner, repo domain.Repo, main domain.MainBranch, ext domain.FileExt, readMe domain.ReadMe) error {

		err := cache.CreateReadMe(ctx, owner, repo, main, ext, readMe)
		if err != nil {
			return err
		}

		return nil
	}
}

type UpdateReadMe func(context.Context, domain.Owner, domain.Repo, domain.MainBranch, domain.FileExt, domain.ReadMe) error

func NewUpdateReadMe(cache domain.ExploreApi) UpdateReadMe {
	return func(ctx context.Context, owner domain.Owner, repo domain.Repo, main domain.MainBranch, ext domain.FileExt, readMe domain.ReadMe) error {

		err := cache.UpdateReadMe(ctx, owner, repo, main, ext, readMe)
		if err != nil {
			return err
		}
		return nil
	}
}
