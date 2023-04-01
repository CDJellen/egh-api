package remote

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/machinebox/graphql"
)

var (
	gitHubHeaders  *http.Header
	readMeHeaders  *http.Header
	readMeEndpoint string
	c              *http.Client
	gql            *graphql.Client
	markdownUrl    = "https://api.github.com/markdown"
	gqlUrl         = "https://api.github.com/graphql"
)

var GqlContributionsQuery string = `
query GetUserContributions($key: String!, $last: Int!) {
	user(login: $key) {
	name
	url
	avatarUrl
	repositoriesContributedTo(last: $last) {
	edges {
		node {
		url
		name
		nameWithOwner
		owner {
			login
			url
			avatarUrl
		}
		}
		}
	}
	}
}
`

type GqlResponse struct {
	User struct {
		Name                      string `json:"name"`
		Url                       string `json:"url"`
		AvatarUrl                 string `json:"avatarUrl"`
		RepositoriesContributedTo struct {
			Edges []struct {
				Node struct {
					Url           string `json:"url"`
					Name          string `json:"name"`
					NameWithOwner string `json:"nameWithOwner"`
					Owner         struct {
						Login     string `json:"login"`
						Url       string `json:"url"`
						AvatarUrl string `json:"avatarUrl"`
					}
				}
			}
		}
	}
}

func newHttpClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 20,
		},
		Timeout: 10 * time.Second,
	}

	return client
}

func newGqlClient() *graphql.Client {
	client := graphql.NewClient(gqlUrl)

	return client
}

func GetGraphQLClient() *graphql.Client {
	if gql == nil {
		gql = newGqlClient()
	}

	return gql
}

func GetHttpClient() *http.Client {
	if c == nil {
		c = newHttpClient()
	}
	return c
}

func GetGitHubHeaders() *http.Header {
	if gitHubHeaders == nil {
		gitHubHeaders = &http.Header{
			"Content-Type":         []string{"application/vnd.github+json"},
			"User-Agent":           []string{"ExploreGitHub"},
			"X-GitHub-Api-Version": []string{"2022-11-28"},
			"Authorization":        []string{fmt.Sprintf("Bearer %s", os.Getenv("GITHUB_TOKEN"))},
		}
	}
	return gitHubHeaders
}

func GetReadMeHeaders() *http.Header {
	if readMeHeaders == nil {
		readMeHeaders = &http.Header{
			"Accept":     []string{"text/html"},
			"User-Agent": []string{"ExploreGitHub"},
		}
	}
	return readMeHeaders
}

func GetReadMeEndpoint() string {
	if readMeEndpoint == "" {
		readMeEndpoint = os.Getenv("README_ENDPOINT")
	}
	return readMeEndpoint
}

func GetMarkdownUrl() string {
	return markdownUrl
}

func UrlEncode(s string) string {
	return strings.Replace(url.QueryEscape(s), "+", "%20", -1)
}
