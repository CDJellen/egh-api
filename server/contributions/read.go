package contributions

import (
	"context"
	"fmt"
	"strings"

	"github.com/machinebox/graphql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cdjellen/egh-api/app"
	"github.com/cdjellen/egh-api/domain"
	pb "github.com/cdjellen/egh-api/pb/proto"
	"github.com/cdjellen/egh-api/server/remote"
	"github.com/cdjellen/egh-api/store"
)

type Read func(context.Context, *pb.ReadContributionsRequest) (*pb.ReadContributionsResponse, error)

func NewRead(handler app.ReadContributions) Read {
	return func(ctx context.Context, req *pb.ReadContributionsRequest) (*pb.ReadContributionsResponse, error) {

		// check the cache
		item, err := handler(ctx, domain.Login(req.Login))
		if err == nil {
			return &pb.ReadContributionsResponse{Message: ToPb(item)}, nil
		} else {
			if !strings.Contains(err.Error(), "cache miss") {
				return nil, status.Error(codes.Internal, err.Error())
			}
		}

		// read from remote
		last := req.GetNumContributions()
		if last == 0 {
			last = int32(30)
		}
		item, err = contributionsRequest(ctx, domain.Login(req.GetLogin()), last)
		if err != nil {
			return &pb.ReadContributionsResponse{Message: ToPb(item)}, err
		}

		// save to cache
		cacher := app.NewCreateContributions(store.NewExploreApiCache())
		cacher(ctx, domain.Login(req.Login), item)

		return &pb.ReadContributionsResponse{Message: ToPb(item)}, nil
	}
}

func contributionsRequest(ctx context.Context, login domain.Login, last int32) (domain.Contributions, error) {

	client := remote.GetGraphQLClient()
	headers := remote.GetGitHubHeaders()

	// make a request
	req := graphql.NewRequest(remote.GqlContributionsQuery)
	req.Header = *headers
	req.Var("key", string(login))
	req.Var("last", last)

	resp := remote.GqlResponse{}
	if err := client.Run(ctx, req, &resp); err != nil {
		fmt.Printf("Failed to unpack request with error %+v", err)
		return domain.Contributions{}, err
	}
	fmt.Printf("response: %s\n%+v", resp, resp)

	reposContributedTo := []domain.Contribution{}

	for _, c := range resp.User.RepositoriesContributedTo.Edges {
		newOwner := domain.RepoOwner{
			Login:     c.Node.Owner.Login,
			Url:       c.Node.Owner.Url,
			AvatarUrl: c.Node.Owner.AvatarUrl,
		}
		newContrib := domain.Contribution{
			NameWithOwner: c.Node.NameWithOwner,
			Name:          c.Node.Name,
			Url:           c.Node.Url,
			Owner:         newOwner,
		}

		reposContributedTo = append(reposContributedTo, newContrib)
	}

	contributions := domain.Contributions{
		Name:          resp.User.Name,
		Url:           resp.User.Url,
		AvatarUrl:     resp.User.AvatarUrl,
		Contributions: reposContributedTo,
	}

	return contributions, nil
}
