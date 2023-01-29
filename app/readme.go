package app

import (
	"context"

	"github.com/cdjellen/egh-api/domain"
)

type ReadReadMe func(context.Context, domain.Owner, domain.Repo, domain.MainBranch, domain.FileExt) (domain.ReadMe, error)

func NewReadReadMe(cache domain.ExploreApiReader) ReadReadMe {
	return func(ctx context.Context, owner domain.Owner, repo domain.Repo, main domain.MainBranch, ext domain.FileExt) (domain.ReadMe, error) {

		item, err := cache.ReadReadMe(ctx, owner, repo, main, ext)
		if err != nil {
			return item, err
		}

		return item, nil
	}
}

type ListReadMe func(context.Context) ([]domain.ReadMe, error)

func NewListReadMe(cache domain.ExploreApiReader) ListReadMe {
	return func(ctx context.Context) ([]domain.ReadMe, error) {
		return cache.ListReadMe(ctx)
	}
}

type CreateReadMe func(context.Context, domain.Owner, domain.Repo, domain.MainBranch, domain.FileExt, domain.ReadMe) error

func NewCreateReadMe(cache domain.ExploreApiWriter) CreateReadMe {
	return func(ctx context.Context, owner domain.Owner, repo domain.Repo, main domain.MainBranch, ext domain.FileExt, readMe domain.ReadMe) error {

		err := cache.CreateReadMe(ctx, owner, repo, main, ext, readMe)
		if err != nil {
			return err
		}

		return nil
	}
}

type UpdateReadMe func(context.Context, domain.Owner, domain.Repo, domain.MainBranch, domain.FileExt, domain.ReadMe) error

func NewUpdateReadMe(cache domain.ExploreApiWriter) UpdateReadMe {
	return func(ctx context.Context, owner domain.Owner, repo domain.Repo, main domain.MainBranch, ext domain.FileExt, readMe domain.ReadMe) error {

		err := cache.UpdateReadMe(ctx, owner, repo, main, ext, readMe)
		if err != nil {
			return err
		}
		return nil
	}
}