package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/redis/go-redis/v9"

	"github.com/cdjellen/egh-api/domain"
)

var (
	instance RedisCache
	once     sync.Once
)

type RedisCache struct {
	mu     sync.Mutex
	Client *redis.Client
}

func NewRedisCache(addr string, user string, pass string, DB int) *RedisCache {
	once.Do(func() {
		opts := redis.Options{Addr: addr}
		if user != "" {
			opts.Username = user
		}
		if pass != "" {
			opts.Password = pass
		}
		if DB != 0 {
			opts.DB = DB
		}
		client := redis.NewClient(&opts)
		instance = RedisCache{Client: client}
	})

	return &instance
}

func (cache *RedisCache) ReadInfo(ctx context.Context, owner domain.Owner, repo domain.Repo) (item domain.Contribution, err error) {
	key := fmt.Sprintf("i-%s/%s", owner, repo)
	val, err := cache.Client.Get(ctx, key).Bytes()
	switch {
	case err == redis.Nil:
		return item, fmt.Errorf("cache miss on repo info %s", key)
	case err != nil:
		return item, err
	}
	err = json.Unmarshal(val, &item)
	if err != nil {
		return item, err
	}
	return item, nil
}

func (cache *RedisCache) ListInfo(ctx context.Context) (items []domain.Contribution, err error) {
	// @TODO move from KEYS to SCAN MATCH
	// obtain all keys from the redis cache
	//i := 0 // redis SCAN MATCH iterator
	//keys := []string{}
	//for {
	//	vals, err := cache.Client.Do(ctx, "SCAN", "MATCH", "*/*").Result()
	//	if err != nil {
	//		return items, err
	//	}
	//	// update the SCAN iterator and keys
	//	i, _ = redis.Int(vals[0], nil)
	//	newKeys, _ := redis.Strings(vals[1], nil)
	//
	//	// clean to ensure no readme keys are appended
	//	infoKeys := []string{}
	//	for _, k := range newKeys {
	//		// README keys have a dot character, info keys do not
	//		if !strings.Contains(k, ".") {
	//			infoKeys = append(infoKeys, k)
	//		}
	//	}
	//	// append only info keys to our keys slice
	//	keys = append(keys, infoKeys...)

	//	// check redis SCAN stopping condition
	//	if i == 0 {
	//		break
	//	}
	//}

	keys := []string{}
	allKeys, err := cache.Client.Do(ctx, "KEYS", "i-*/*").StringSlice()

	for _, k := range allKeys {
		if !strings.Contains(k, ".") {
			keys = append(keys, k)
		}
	}

	// iterate through the obtained keys
	for _, k := range keys {
		k = strings.TrimLeft(k, "i-")
		owner := domain.Owner(strings.Split(k, "/")[0])
		repo := domain.Repo(strings.Split(k, "/")[1])

		item, err := cache.ReadInfo(ctx, owner, repo)
		if err != nil {
			log.Printf("Failed to read item with key %s from redis cache.", k)
			continue
		}
		items = append(items, item)
	}

	return items, nil
}

func (cache *RedisCache) CreateInfo(ctx context.Context, owner domain.Owner, repo domain.Repo, item domain.Contribution) error {
	key := fmt.Sprintf("i-%s/%s", owner, repo)

	cache.mu.Lock()
	defer cache.mu.Unlock()

	val, err := json.Marshal(item)
	if err != nil {
		log.Printf("Failed to marshall item of type %t with value %+v to byte slice", item, item)
		return err
	}
	res, err := cache.Client.Set(ctx, key, val, 0).Result()
	if err != nil {
		log.Printf("Failed to set key %s to value %+v", key, val)
		return err
	}

	log.Printf("Set key %s with result %s", key, res)

	return nil
}

func (cache *RedisCache) UpdateInfo(ctx context.Context, owner domain.Owner, repo domain.Repo, item domain.Contribution) error {
	_, err := cache.ReadInfo(ctx, owner, repo)
	if err != nil {
		return err
	}

	err = cache.CreateInfo(ctx, owner, repo, item)
	if err != nil {
		return err
	}

	return nil
}

func (cache *RedisCache) ReadContributors(ctx context.Context, owner domain.Owner, repo domain.Repo) (item domain.RepoContributors, err error) {
	key := fmt.Sprintf("c-%s/%s", owner, repo)
	val, err := cache.Client.Get(ctx, key).Bytes()
	switch {
	case err == redis.Nil:
		return item, fmt.Errorf("cache miss on repo contributors %s", key)
	case err != nil:
		return item, err
	}
	err = json.Unmarshal(val, &item)
	if err != nil {
		return item, err
	}

	return item, nil

}

func (cache *RedisCache) ListContributors(ctx context.Context) (items []domain.RepoContributors, err error) {
	keys := []string{}
	allKeys, err := cache.Client.Do(ctx, "KEYS", "c-*/*").StringSlice()

	for _, k := range allKeys {
		if !strings.Contains(k, ".") {
			keys = append(keys, k)
		}
	}

	// iterate through the obtained keys
	for _, k := range keys {
		k = strings.TrimLeft(k, "c-")
		owner := domain.Owner(strings.Split(k, "/")[0])
		repo := domain.Repo(strings.Split(k, "/")[1])

		item, err := cache.ReadContributors(ctx, owner, repo)
		if err != nil {
			log.Printf("Failed to read item with key %s from redis cache.", k)
			continue
		}
		items = append(items, item)
	}

	return items, nil

}

func (cache *RedisCache) CreateContributors(ctx context.Context, owner domain.Owner, repo domain.Repo, item domain.RepoContributors) error {
	key := fmt.Sprintf("c-%s/%s", owner, repo)

	cache.mu.Lock()
	defer cache.mu.Unlock()

	val, err := json.Marshal(item)
	if err != nil {
		log.Printf("Failed to marshall item of type %t with value %+v to byte slice", item, item)
		return err
	}
	res, err := cache.Client.Set(ctx, key, val, 0).Result()
	if err != nil {
		log.Printf("Failed to set key %s to value %+v", key, val)
		return err
	}

	log.Printf("Set key %s with result %s", key, res)

	return nil

}

func (cache *RedisCache) UpdateContributors(ctx context.Context, owner domain.Owner, repo domain.Repo, item domain.RepoContributors) error {
	_, err := cache.ReadContributors(ctx, owner, repo)
	if err != nil {
		return err
	}

	err = cache.CreateContributors(ctx, owner, repo, item)
	if err != nil {
		return err
	}

	return nil

}

func (cache *RedisCache) ReadContributions(ctx context.Context, login domain.Login) (item domain.Contributions, err error) {
	key := fmt.Sprintf("l-%s", login)
	val, err := cache.Client.Get(ctx, key).Bytes()
	switch {
	case err == redis.Nil:
		return item, fmt.Errorf("cache miss on login %s", key)
	case err != nil:
		return item, err
	}
	err = json.Unmarshal(val, &item)
	if err != nil {
		return item, err
	}

	return item, nil
}

func (cache *RedisCache) ListContributions(ctx context.Context) (items []domain.Contributions, err error) {
	keys, err := cache.Client.Do(ctx, "KEYS", "l-*").StringSlice()

	// iterate through the obtained keys
	for _, k := range keys {
		login := domain.Login(strings.TrimLeft(k, "l-"))

		item, err := cache.ReadContributions(ctx, login)
		if err != nil {
			log.Printf("Failed to read item with key %s from redis cache.", k)
			continue
		}
		items = append(items, item)
	}

	return items, nil

}

func (cache *RedisCache) CreateContributions(ctx context.Context, login domain.Login, item domain.Contributions) error {
	key := fmt.Sprintf("l-%s", login)

	cache.mu.Lock()
	defer cache.mu.Unlock()

	val, err := json.Marshal(item)
	if err != nil {
		log.Printf("Failed to marshall item of type %t with value %+v to byte slice", item, item)
		return err
	}
	res, err := cache.Client.Set(ctx, key, val, 0).Result()
	if err != nil {
		log.Printf("Failed to set key %s to value %+v", key, val)
		return err
	}

	log.Printf("Set key %s with result %s", key, res)

	return nil
}

func (cache *RedisCache) UpdateContributions(ctx context.Context, login domain.Login, item domain.Contributions) error {
	_, err := cache.ReadContributions(ctx, login)
	if err != nil {
		return err
	}

	err = cache.CreateContributions(ctx, login, item)
	if err != nil {
		return err
	}

	return nil

}

func (cache *RedisCache) ReadReadMe(ctx context.Context, owner domain.Owner, repo domain.Repo, main domain.MainBranch, ext domain.FileExt) (item domain.ReadMe, err error) {
	key := fmt.Sprintf("r-%s/%s/%s.%s", owner, repo, main, ext)

	val, err := cache.Client.Get(ctx, key).Bytes()
	switch {
	case err == redis.Nil:
		return item, fmt.Errorf("cache miss on repo contributors %s", key)
	case err != nil:
		return item, err
	}
	err = json.Unmarshal(val, &item)
	if err != nil {
		return item, err
	}

	return item, nil

}

func (cache *RedisCache) ListReadMe(ctx context.Context) (items []domain.ReadMe, err error) {
	keys, err := cache.Client.Do(ctx, "KEYS", "r-*").StringSlice()

	// iterate through the obtained keys
	for _, k := range keys {
		k = strings.TrimLeft(k, "r-")
		components := strings.Split(k, "/")
		owner := domain.Owner(components[0])
		repo := domain.Repo(components[1])
		// file and extension must be split
		mainExt := strings.Split(components[2], ".")
		main := domain.MainBranch(mainExt[0])
		ext := domain.FileExt(mainExt[1])

		item, err := cache.ReadReadMe(ctx, owner, repo, main, ext)
		if err != nil {
			log.Printf("Failed to read item with key %s from redis cache.", k)
			continue
		}
		items = append(items, item)
	}

	return items, nil

}

func (cache *RedisCache) CreateReadMe(ctx context.Context, owner domain.Owner, repo domain.Repo, main domain.MainBranch, ext domain.FileExt, item domain.ReadMe) error {
	key := fmt.Sprintf("r-%s/%s/%s.%s", owner, repo, main, ext)

	cache.mu.Lock()
	defer cache.mu.Unlock()

	val, err := json.Marshal(item)
	if err != nil {
		log.Printf("Failed to marshall item of type %t with value %+v to byte slice", item, item)
		return err
	}
	res, err := cache.Client.Set(ctx, key, val, 0).Result()
	if err != nil {
		log.Printf("Failed to set key %s to value %+v", key, val)
		return err
	}

	log.Printf("Set key %s with result %s", key, res)

	return nil

}

func (cache *RedisCache) UpdateReadMe(ctx context.Context, owner domain.Owner, repo domain.Repo, main domain.MainBranch, ext domain.FileExt, item domain.ReadMe) error {
	_, err := cache.ReadReadMe(ctx, owner, repo, main, ext)
	if err != nil {
		return err
	}

	err = cache.CreateReadMe(ctx, owner, repo, main, ext, item)
	if err != nil {
		return err
	}

	return nil

}
