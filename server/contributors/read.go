package contributors

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cdjellen/egh-api/app"
	"github.com/cdjellen/egh-api/domain"
	pb "github.com/cdjellen/egh-api/pb/proto"
	"github.com/cdjellen/egh-api/server/remote"
)

type Read func(context.Context, *pb.ReadContributorsRequest) (*pb.ReadContributorsResponse, error)

func NewRead(handler app.ReadContributors) Read {
	return func(ctx context.Context, req *pb.ReadContributorsRequest) (*pb.ReadContributorsResponse, error) {

		// check the cache
		item, err := handler(ctx, domain.Owner(req.Owner), domain.Repo(req.Repo))
		if err == nil {
			return &pb.ReadContributorsResponse{Message: ToPb(item)}, nil
		} else {
			if !strings.Contains(err.Error(), "cache miss") {
				return nil, status.Error(codes.Internal, err.Error())
			}
		}

		// read from remote
		item, err = contributorRequest(ctx, domain.Owner(req.Owner), domain.Repo(req.Repo), req.GetAnon(), int(req.GetPerPage()), int(req.GetPage()))
		if err != nil {
			fmt.Printf("failed to get CONTRIBUTORS with error %+v", err)
			return &pb.ReadContributorsResponse{Message: ToPb(item)}, err
		}

		return &pb.ReadContributorsResponse{Message: ToPb(item)}, nil
	}
}

func contributorsUrl(o domain.Owner, r domain.Repo, anon string, perPage int, page int) string {
	baseUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/contributors?per_page=%s", o, r, fmt.Sprint(perPage))
	if anon != "" {
		baseUrl = fmt.Sprintf("%s?anon=%s", baseUrl, anon)
	}
	if page != 0 {
		baseUrl = fmt.Sprintf("%s?page=%s", baseUrl, fmt.Sprint(page))
	}

	return baseUrl
}

func contributorRequest(ctx context.Context, o domain.Owner, r domain.Repo, anon string, perPage int, page int) (domain.RepoContributors, error) {
	if perPage == 0 {
		perPage = 50
	}

	headers := remote.GetGitHubHeaders()
	url := contributorsUrl(o, r, anon, perPage, page)
	fmt.Printf("URL: %s", url)
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
		fmt.Printf("\n\n%+v\n", err)
		return domain.RepoContributors{}, err
	}

	return domain.RepoContributors{RepoContributors: item}, nil
}
