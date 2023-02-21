package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/cdjellen/egh-api/app/remote"
	"github.com/cdjellen/egh-api/domain"
)

type ReadContributors func(context.Context, domain.Owner, domain.Repo, int32) (domain.RepoContributors, error)

func NewReadContributors(cache domain.ExploreApi) ReadContributors {
	return func(ctx context.Context, owner domain.Owner, repo domain.Repo, topN int32) (domain.RepoContributors, error) {

		// check the cache
		item, err := cache.ReadContributors(ctx, owner, repo, topN)
		if err == nil {
			return item, nil
		}
		if strings.Contains(err.Error(), "cache miss") {
			// read from remote
			item, err = contributorRequest(ctx, owner, repo, "false", 100, 1)
			if err != nil {
				log.Printf("failed to get CONTRIBUTORS with error %+v", err)
				return item, err
			}

			// persist to cache
			err = cache.CreateContributors(ctx, owner, repo, item)
			if err != nil {
				return item, err
			}
		}
		// pull from cache with params
		item, err = cache.ReadContributors(ctx, owner, repo, topN)
		if err == nil {
			return item, nil
		}

		return item, err
	}
}

func contributorsUrl(o domain.Owner, r domain.Repo, anon string, perPage int32, page int32) string {
	baseUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/contributors?per_page=%s", o, r, fmt.Sprint(perPage))
	if anon != "" {
		baseUrl = fmt.Sprintf("%s?anon=%s", baseUrl, anon)
	}
	if page != 0 {
		baseUrl = fmt.Sprintf("%s?page=%s", baseUrl, fmt.Sprint(page))
	}

	return baseUrl
}

func contributorRequest(ctx context.Context, o domain.Owner, r domain.Repo, anon string, perPage int32, page int32) (domain.RepoContributors, error) {
	if perPage == 0 {
		perPage = 25
	}

	headers := remote.GetGitHubHeaders()
	url := contributorsUrl(o, r, anon, perPage, page)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return domain.RepoContributors{}, err
	}
	req.Header = *headers

	client := remote.GetHttpClient()
	resp, err := client.Do(req)
	if err != nil {
		return domain.RepoContributors{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return domain.RepoContributors{}, err
	}

	var item []domain.Contributors
	if err != nil {
		return domain.RepoContributors{}, err
	}
	log.Printf("body:\n%+v\n%s", body, body)
	err = json.Unmarshal(body, &item)
	if err != nil {
		log.Printf("\n\n%+v\n", err)
		return domain.RepoContributors{}, err
	}

	return domain.RepoContributors{RepoContributors: item}, nil
}

type ListContributors func(context.Context) ([]domain.RepoContributors, error)

func NewListContributors(cache domain.ExploreApi) ListContributors {
	return func(ctx context.Context) ([]domain.RepoContributors, error) {
		return cache.ListContributors(ctx)
	}
}

type CreateContributors func(context.Context, domain.Owner, domain.Repo, domain.RepoContributors) error

func NewCreateContributors(cache domain.ExploreApi) CreateContributors {
	return func(ctx context.Context, owner domain.Owner, repo domain.Repo, repoContributors domain.RepoContributors) error {

		err := cache.CreateContributors(ctx, owner, repo, repoContributors)
		if err != nil {
			return err
		}

		return nil
	}
}

type UpdateContributors func(context.Context, domain.Owner, domain.Repo, domain.RepoContributors) error

func NewUpdateContributors(cache domain.ExploreApi) UpdateContributors {
	return func(ctx context.Context, owner domain.Owner, repo domain.Repo, repoContributors domain.RepoContributors) error {

		err := cache.UpdateContributors(ctx, owner, repo, repoContributors)
		if err != nil {
			return err
		}
		return nil
	}
}
