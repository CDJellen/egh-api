package info

import (
	"context"
	"encoding/json"
	"errors"
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
)

type Read func(context.Context, *pb.ReadInfoRequest) (*pb.ReadInfoResponse, error)

func NewRead(handler app.ReadInfo) Read {
	return func(ctx context.Context, req *pb.ReadInfoRequest) (*pb.ReadInfoResponse, error) {

		// check the cache
		item, err := handler(ctx, domain.Owner(req.Owner), domain.Repo(req.Repo))
		if err == nil {
			return &pb.ReadInfoResponse{Message: ToPb(item)}, nil
		} else {
			if !strings.Contains(err.Error(), "cache miss") {
				return nil, status.Error(codes.Internal, err.Error())
			}
		}

		// read from remote
		item, err = infoRequest(ctx, domain.Owner(req.Owner), domain.Repo(req.Repo))
		if err != nil {
			fmt.Printf("failed to get INFO with error %+v", err)
			return &pb.ReadInfoResponse{Message: ToPb(item)}, err
		}

		return &pb.ReadInfoResponse{Message: ToPb(item)}, errors.New("cache miss")
	}
}

func infoUrl(o domain.Owner, r domain.Repo) string {
	baseUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s", o, r)

	return baseUrl
}

func infoRequest(ctx context.Context, o domain.Owner, r domain.Repo) (domain.Contribution, error) {
	headers := remote.GetGitHubHeaders()
	url := infoUrl(o, r)
	fmt.Printf("URL: %s", url)
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
		fmt.Printf("\n\n%+v\n", err)
		return domain.Contribution{}, err
	}

	return item, nil
}
