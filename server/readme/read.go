package readme

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cdjellen/egh-api/app"
	"github.com/cdjellen/egh-api/domain"
	pb "github.com/cdjellen/egh-api/pb/proto"
	"github.com/cdjellen/egh-api/server/remote"
	"github.com/cdjellen/egh-api/store"
)

type Read func(context.Context, *pb.ReadReadMeRequest) (*pb.ReadReadMeResponse, error)

func NewRead(handler app.ReadReadMe) Read {
	return func(ctx context.Context, req *pb.ReadReadMeRequest) (*pb.ReadReadMeResponse, error) {

		// check the cache
		item, err := handler(ctx, domain.Owner(req.Owner), domain.Repo(req.Repo), domain.MainBranch(req.MainBranch), domain.FileExt(req.FileExt))
		if err == nil {
			return &pb.ReadReadMeResponse{Message: ToPb(item)}, nil
		} else {
			if !strings.Contains(err.Error(), "cache miss") {
				return nil, status.Error(codes.Internal, err.Error())
			}
		}

		// read from remote
		item, err = readMeRequest(ctx, domain.Owner(req.Owner), domain.Repo(req.Repo), domain.MainBranch(req.MainBranch), domain.FileExt(req.FileExt))
		if err != nil {
			fmt.Printf("failed to get README")
		}

		// save to cache
		cacher := app.NewCreateReadMe(store.NewExploreApiCache())
		cacher(ctx, domain.Owner(req.Owner), domain.Repo(req.Repo), domain.MainBranch(req.MainBranch), domain.FileExt(req.FileExt), item)

		return &pb.ReadReadMeResponse{Message: ToPb(item)}, nil
	}
}

func readMeUrl(o domain.Owner, r domain.Repo, m domain.MainBranch, e domain.FileExt) string {
	encodedElem := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/README.%s", o, r, m, e)
	return fmt.Sprintf("%s%s", remote.GetReadMeEndpoint(), remote.UrlEncode(encodedElem))
}

func readMeNotFound(o domain.Owner, r domain.Repo) (resp string) {
	return fmt.Sprintf("<h1>Failed to Load: <a href=github.com/%s/%s>github.com/%s/%s</a></h1>", o, r, o, r)
}

func readMeMdRequest(ctx context.Context, o domain.Owner, r domain.Repo, m domain.MainBranch, e domain.FileExt) (status int, body []byte, err error) {
	headers := remote.GetReadMeHeaders()
	url := readMeUrl(o, r, m, e)
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
