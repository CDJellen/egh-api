package app

import (
	"context"
	"strings"

	"github.com/cdjellen/egh-api/domain"
)

type ReadContributions func(context.Context, domain.Login) (domain.Contributions, error)

func NewReadContributions(cache domain.ExploreApi) ReadContributions {
	return func(ctx context.Context, login domain.Login) (domain.Contributions, error) {

		item, err := cache.ReadContributions(ctx, login)
		if err == nil {
			return item, nil
		}
		if strings.Contains(err.Error(), "cache miss") {
			err = cache.CreateContributions(ctx, login, item)
			if err != nil {
				return item, err
			}

			return item, nil

		}

		return item, err
	}
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
