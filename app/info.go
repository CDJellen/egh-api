package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/cdjellen/egh-api/domain"
	"github.com/cdjellen/egh-api/server/remote"
)

type ReadInfo func(context.Context, domain.Owner, domain.Repo) (domain.Contribution, error)

func NewReadInfo(cache domain.ExploreApi) ReadInfo {
	return func(ctx context.Context, owner domain.Owner, repo domain.Repo) (domain.Contribution, error) {

		// check the cache
		item, err := cache.ReadInfo(ctx, owner, repo)
		if err == nil {
			return item, nil
		}
		if strings.Contains(err.Error(), "cache miss") {
			// read from remote
			item, err = infoRequest(ctx, owner, repo)
			if err != nil {
				log.Printf("failed to get INFO with error %+v", err)
				return item, err
			}

			// persist to cache
			err = cache.CreateInfo(ctx, owner, repo, item)
			if err != nil {
				return item, err
			}

			return item, nil

		}

		return item, err
	}
}

func infoUrl(o domain.Owner, r domain.Repo) string {
	baseUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s", o, r)

	return baseUrl
}

func infoRequest(ctx context.Context, o domain.Owner, r domain.Repo) (domain.Contribution, error) {
	headers := remote.GetGitHubHeaders()
	url := infoUrl(o, r)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return domain.Contribution{}, err
	}
	req.Header = *headers

	client := remote.GetHttpClient()
	resp, err := client.Do(req)
	if err != nil {
		return domain.Contribution{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return domain.Contribution{}, err
	}

	var item domain.Contribution
	if err != nil {
		return domain.Contribution{}, err
	}
	err = json.Unmarshal(body, &item)
	if err != nil {
		log.Printf("\n\n%+v\n", err)
		return domain.Contribution{}, err
	}

	return item, nil
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
