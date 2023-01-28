package store

import (
	"context"
	"fmt"
	"sync"

	"github.com/cdjellen/egh-api/domain"
)

var (
	instance ExploreApiCache
	once     sync.Once
)

type ExploreApiCache struct {
	mu                 sync.Mutex
	infoCache          map[string]domain.Contribution
	contributorsCache  map[string]domain.RepoContributors
	contributionsCache map[string]domain.Contributions
	readMeCache        map[string]domain.ReadMe
}

func NewExploreApiCache() *ExploreApiCache {
	once.Do(func() {
		instance = ExploreApiCache{
			infoCache:          make(map[string]domain.Contribution),
			contributionsCache: make(map[string]domain.Contributions),
			contributorsCache:  make(map[string]domain.RepoContributors),
			readMeCache:        make(map[string]domain.ReadMe),
		}
	})

	return &instance
}

func (cache *ExploreApiCache) ReadInfo(ctx context.Context, owner domain.Owner, repo domain.Repo) (item domain.Contribution, err error) {
	key := fmt.Sprintf("%s/%s", owner, repo)
	if item, ok := cache.infoCache[key]; ok {
		return item, nil
	}
	return item, fmt.Errorf("cache miss on repo %s", key)
}

func (cache *ExploreApiCache) ListInfo(ctx context.Context) (items []domain.Contribution, err error) {
	for _, v := range cache.infoCache {
		items = append(items, v)
	}
	return items, nil
}

func (cache *ExploreApiCache) CreateInfo(ctx context.Context, owner domain.Owner, repo domain.Repo, item domain.Contribution) error {
	key := fmt.Sprintf("%s/%s", owner, repo)

	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.infoCache[key] = item
	return nil
}

func (cache *ExploreApiCache) UpdateInfo(ctx context.Context, owner domain.Owner, repo domain.Repo, item domain.Contribution) error {
	_, err := cache.ReadInfo(ctx, owner, repo)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s/%s", owner, repo)

	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.infoCache[key] = item
	return nil
}

func (cache *ExploreApiCache) ReadContributors(ctx context.Context, owner domain.Owner, repo domain.Repo) (item domain.RepoContributors, err error) {
	key := fmt.Sprintf("%s/%s", owner, repo)
	if item, ok := cache.contributorsCache[key]; ok {
		return item, nil
	}
	return item, fmt.Errorf("cache miss on repo %s", key)
}

func (cache *ExploreApiCache) ListContributors(ctx context.Context) (items []domain.RepoContributors, err error) {
	for _, v := range cache.contributorsCache {
		items = append(items, v)
	}
	return items, nil
}

func (cache *ExploreApiCache) CreateContributors(ctx context.Context, owner domain.Owner, repo domain.Repo, item domain.RepoContributors) error {
	key := fmt.Sprintf("%s/%s", owner, repo)

	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.contributorsCache[key] = item
	return nil
}

func (cache *ExploreApiCache) UpdateContributors(ctx context.Context, owner domain.Owner, repo domain.Repo, item domain.RepoContributors) error {
	_, err := cache.ReadContributors(ctx, owner, repo)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s/%s", owner, repo)

	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.contributorsCache[key] = item
	return nil
}

func (cache *ExploreApiCache) ReadContributions(ctx context.Context, login domain.Login) (item domain.Contributions, err error) {
	if item, ok := cache.contributionsCache[string(login)]; ok {
		return item, nil
	}
	return item, fmt.Errorf("cache miss on user %s", login)
}

func (cache *ExploreApiCache) ListContributions(ctx context.Context) (items []domain.Contributions, err error) {
	for _, v := range cache.contributionsCache {
		items = append(items, v)
	}
	return items, nil
}

func (cache *ExploreApiCache) CreateContributions(ctx context.Context, login domain.Login, item domain.Contributions) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.contributionsCache[string(login)] = item
	return nil
}

func (cache *ExploreApiCache) UpdateContributions(ctx context.Context, login domain.Login, item domain.Contributions) error {
	_, err := cache.ReadContributions(ctx, login)
	if err != nil {
		return err
	}

	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.contributionsCache[string(login)] = item
	return nil
}

func (cache *ExploreApiCache) ReadReadMe(ctx context.Context, owner domain.Owner, repo domain.Repo, main domain.MainBranch, ext domain.FileExt) (item domain.ReadMe, err error) {
	key := fmt.Sprintf("%s/%s/%s.%s", owner, repo, main, ext)

	if item, ok := cache.readMeCache[key]; ok {
		return item, nil
	}
	return item, fmt.Errorf("cache miss on readme %s", key)
}

func (cache *ExploreApiCache) ListReadMe(ctx context.Context) (items []domain.ReadMe, err error) {
	for _, v := range cache.readMeCache {
		items = append(items, v)
	}
	return items, nil
}

func (cache *ExploreApiCache) CreateReadMe(ctx context.Context, owner domain.Owner, repo domain.Repo, main domain.MainBranch, ext domain.FileExt, item domain.ReadMe) error {
	key := fmt.Sprintf("%s/%s/%s.%s", owner, repo, main, ext)

	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.readMeCache[key] = item
	return nil
}

func (cache *ExploreApiCache) UpdateReadMe(ctx context.Context, owner domain.Owner, repo domain.Repo, main domain.MainBranch, ext domain.FileExt, item domain.ReadMe) error {
	_, err := cache.ReadReadMe(ctx, owner, repo, main, ext)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s/%s/%s.%s", owner, repo, main, ext)

	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.readMeCache[key] = item
	return nil
}
