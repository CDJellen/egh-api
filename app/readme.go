package app

import (
	"bytes"
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

type ReadReadMe func(context.Context, domain.Owner, domain.Repo, domain.MainBranch, domain.FileExt) (domain.ReadMe, error)

func NewReadReadMe(cache domain.ExploreApi) ReadReadMe {
	return func(ctx context.Context, owner domain.Owner, repo domain.Repo, main domain.MainBranch, ext domain.FileExt) (domain.ReadMe, error) {

		// check the cache
		item, err := cache.ReadReadMe(ctx, owner, repo, main, ext)
		if err == nil {
			return item, nil
		}
		if strings.Contains(err.Error(), "cache miss") {
			// read from remote
			item, err = readMeRequest(ctx, owner, repo, main, ext)
			if err != nil {
				fmt.Printf("failed to get README with error %+v", err)
				return item, err
			}

			// persist to cache
			err = cache.CreateReadMe(ctx, owner, repo, main, ext, item)
			if err != nil {
				return item, err
			}

			return item, nil

		}

		return item, err
	}
}

func readMeUrl(o domain.Owner, r domain.Repo, m domain.MainBranch, e domain.FileExt) string {
	encodedElem := remote.UrlEncode(fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/README.%s", o, r, m, e))
	return fmt.Sprintf("%s%s", remote.GetReadMeEndpoint(), encodedElem)
}

func readMeNotFound(o domain.Owner, r domain.Repo) (resp string) {
	return fmt.Sprintf("<h1>Failed to Load: <a href=github.com/%s/%s>github.com/%s/%s</a></h1>", o, r, o, r)
}

func readMeMdRequest(ctx context.Context, o domain.Owner, r domain.Repo, m domain.MainBranch, e domain.FileExt) (status int, body []byte, err error) {
	headers := remote.GetReadMeHeaders()
	url := readMeUrl(o, r, m, e)
	log.Printf("url: %s, headers: %v", url, *headers)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 500, body, err
	}
	req.Header = *headers

	client := remote.GetHttpClient()
	resp, err := client.Do(req)
	if err != nil {
		return 500, body, err
	}
	if resp.StatusCode != 200 {
		return resp.StatusCode, body, nil
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return 500, body, err
	}

	return resp.StatusCode, body, nil
}

func readMeHtmlRequest(ctx context.Context, body []byte, owner domain.Owner, repo domain.Repo) (string, error) {
	headers := remote.GetGitHubHeaders()
	url := remote.GetMarkdownUrl()
	bodyData := map[string]string{"text": string(body)}
	postData, err := json.Marshal(bodyData)
	if err != nil {
		return readMeNotFound(owner, repo), err
	}
	postBody := bytes.NewBuffer(postData)

	req, err := http.NewRequestWithContext(ctx, "POST", url, postBody)
	if err != nil {
		return readMeNotFound(owner, repo), err
	}
	req.Header = *headers

	client := remote.GetHttpClient()
	resp, err := client.Do(req)
	if err != nil {
		return readMeNotFound(owner, repo), err
	}
	if resp.StatusCode != 200 {
		return readMeNotFound(owner, repo), nil
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return readMeNotFound(owner, repo), err
	}

	return string(body), nil

}

func readMeRequest(ctx context.Context, o domain.Owner, r domain.Repo, m domain.MainBranch, e domain.FileExt) (domain.ReadMe, error) {
	mdStatus, body, err := readMeMdRequest(ctx, o, r, m, e)
	log.Printf("status: %v, error: %+v, body: %s", mdStatus, err, body)
	if err != nil {
		return domain.ReadMe{}, err
	}
	if mdStatus != 200 {
		return domain.ReadMe{Html: readMeNotFound(o, r)}, nil
	}

	html, err := readMeHtmlRequest(ctx, body, o, r)
	if err != nil {
		return domain.ReadMe{Html: html}, err
	}

	return domain.ReadMe{Html: html}, nil

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
