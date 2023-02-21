package store

import (
	"context"

	"github.com/cdjellen/egh-api/domain"
)

type ExploreApiCache interface {
	ReadInfo(ctx context.Context, owner domain.Owner, repo domain.Repo) (item domain.Contribution, err error)
	ListInfo(ctx context.Context) (items []domain.Contribution, err error)
	CreateInfo(ctx context.Context, owner domain.Owner, repo domain.Repo, item domain.Contribution) error
	UpdateInfo(ctx context.Context, owner domain.Owner, repo domain.Repo, item domain.Contribution) error
	ReadContributors(ctx context.Context, owner domain.Owner, repo domain.Repo, topN int32) (item domain.RepoContributors, err error)
	ListContributors(ctx context.Context) (items []domain.RepoContributors, err error)
	CreateContributors(ctx context.Context, owner domain.Owner, repo domain.Repo, item domain.RepoContributors) error
	UpdateContributors(ctx context.Context, owner domain.Owner, repo domain.Repo, item domain.RepoContributors) error
	ReadContributions(ctx context.Context, login domain.Login, numContributions int32) (item domain.Contributions, err error)
	ListContributions(ctx context.Context) (items []domain.Contributions, err error)
	CreateContributions(ctx context.Context, login domain.Login, item domain.Contributions) error
	UpdateContributions(ctx context.Context, login domain.Login, item domain.Contributions) error
	ReadReadMe(ctx context.Context, owner domain.Owner, repo domain.Repo, main domain.MainBranch, ext domain.FileExt) (item domain.ReadMe, err error)
	ListReadMe(ctx context.Context) (items []domain.ReadMe, err error)
	CreateReadMe(ctx context.Context, owner domain.Owner, repo domain.Repo, main domain.MainBranch, ext domain.FileExt, item domain.ReadMe) error
	UpdateReadMe(ctx context.Context, owner domain.Owner, repo domain.Repo, main domain.MainBranch, ext domain.FileExt, item domain.ReadMe) error
}
