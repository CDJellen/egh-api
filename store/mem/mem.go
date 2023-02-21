package mem

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
	key := fmt.Sprintf("i-%s/%s", owner, repo)
	if item, ok := cache.infoCache[key]; ok {
		return item, nil
	}

	return item, fmt.Errorf("cache miss on repo info %s", key)

}

func (cache *ExploreApiCache) ListInfo(ctx context.Context) (items []domain.Contribution, err error) {
	for _, v := range cache.infoCache {
		items = append(items, v)
	}

	return items, nil

}

func (cache *ExploreApiCache) CreateInfo(ctx context.Context, owner domain.Owner, repo domain.Repo, item domain.Contribution) error {
	key := fmt.Sprintf("i-%s/%s", owner, repo)

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

	key := fmt.Sprintf("i-%s/%s", owner, repo)

	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.infoCache[key] = item

	return nil

}

func (cache *ExploreApiCache) ReadContributors(ctx context.Context, owner domain.Owner, repo domain.Repo, topN int32) (item domain.RepoContributors, err error) {
	key := fmt.Sprintf("c-%s/%s", owner, repo)
	if item, ok := cache.contributorsCache[key]; ok {
		// get the appropriate number of items
		if len(item.RepoContributors) > int(topN) {
			topNItems := item.RepoContributors[0:topN]
			item = domain.RepoContributors{RepoContributors: topNItems}
		}

		return item, nil
	}

	return item, fmt.Errorf("cache miss on repo contributors %s", key)

}

func (cache *ExploreApiCache) ListContributors(ctx context.Context) (items []domain.RepoContributors, err error) {
	for _, v := range cache.contributorsCache {
		items = append(items, v)
	}

	return items, nil

}

func (cache *ExploreApiCache) CreateContributors(ctx context.Context, owner domain.Owner, repo domain.Repo, item domain.RepoContributors) error {
	key := fmt.Sprintf("c-%s/%s", owner, repo)

	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.contributorsCache[key] = item

	return nil

}

func (cache *ExploreApiCache) UpdateContributors(ctx context.Context, owner domain.Owner, repo domain.Repo, item domain.RepoContributors) error {
	_, err := cache.ReadContributors(ctx, owner, repo, 100)
	if err != nil {
		return err
	}

	cache.mu.Lock()
	defer cache.mu.Unlock()

	err = cache.CreateContributors(ctx, owner, repo, item)
	if err != nil {
		return err
	}

	return nil

}

func (cache *ExploreApiCache) ReadContributions(ctx context.Context, login domain.Login, numContributions int32) (item domain.Contributions, err error) {
	key := fmt.Sprintf("l-%s", login)

	if item, ok := cache.contributionsCache[key]; ok {
		// get the appropriate number of items
		if len(item.Contributions) > int(numContributions) {
			name := item.Name
			url := item.Url
			avatarUrl := item.AvatarUrl
			topNItems := item.Contributions[0:numContributions]
			// set item to include only `numContributions` contributions
			item = domain.Contributions{Name: name, Url: url, AvatarUrl: avatarUrl, Contributions: topNItems}
		}
	}

	return item, fmt.Errorf("cache miss on login %s", key)

}

func (cache *ExploreApiCache) ListContributions(ctx context.Context) (items []domain.Contributions, err error) {
	for _, v := range cache.contributionsCache {
		items = append(items, v)
	}

	return items, nil

}

func (cache *ExploreApiCache) CreateContributions(ctx context.Context, login domain.Login, item domain.Contributions) error {
	key := fmt.Sprintf("l-%s", login)

	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.contributionsCache[key] = item

	return nil

}

func (cache *ExploreApiCache) UpdateContributions(ctx context.Context, login domain.Login, item domain.Contributions) error {
	_, err := cache.ReadContributions(ctx, login, 100)
	if err != nil {
		return err
	}

	cache.mu.Lock()
	defer cache.mu.Unlock()

	err = cache.CreateContributions(ctx, login, item)
	if err != nil {
		return err
	}

	return nil

}

func (cache *ExploreApiCache) ReadReadMe(ctx context.Context, owner domain.Owner, repo domain.Repo, main domain.MainBranch, ext domain.FileExt) (item domain.ReadMe, err error) {
	key := fmt.Sprintf("r-%s/%s/%s.%s", owner, repo, main, ext)

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
	key := fmt.Sprintf("r-%s/%s/%s.%s", owner, repo, main, ext)

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

	cache.mu.Lock()
	defer cache.mu.Unlock()

	err = cache.CreateReadMe(ctx, owner, repo, main, ext, item)
	if err != nil {
		return err
	}

	return nil

}
