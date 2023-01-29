package app

import (
	"context"

	"github.com/cdjellen/egh-api/domain"
)

type ReadHealth func(context.Context) (domain.Health, error)

func NewReadHealth() ReadHealth {
	return func(ctx context.Context) (domain.Health, error) {

		item := domain.Health{
			Status: "alive",
		}
		return item, nil
	}
}
