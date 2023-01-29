package domain

import (
	"context"
)

type ExploreApiReader interface {
	ReadInfo(context.Context, Owner, Repo) (Contribution, error)
	ListInfo(context.Context) ([]Contribution, error)
	ReadContributors(context.Context, Owner, Repo) (RepoContributors, error)
	ListContributors(context.Context) ([]RepoContributors, error)
	ReadContributions(context.Context, Login) (Contributions, error)
	ListContributions(context.Context) ([]Contributions, error)
	ReadReadMe(context.Context, Owner, Repo, MainBranch, FileExt) (ReadMe, error)
	ListReadMe(context.Context) ([]ReadMe, error)
}

type ExploreApiWriter interface {
	CreateInfo(context.Context, Owner, Repo, Contribution) error
	UpdateInfo(context.Context, Owner, Repo, Contribution) error
	CreateContributors(context.Context, Owner, Repo, RepoContributors) error
	UpdateContributors(context.Context, Owner, Repo, RepoContributors) error
	CreateContributions(context.Context, Login, Contributions) error
	UpdateContributions(context.Context, Login, Contributions) error
	CreateReadMe(context.Context, Owner, Repo, MainBranch, FileExt, ReadMe) error
	UpdateReadMe(context.Context, Owner, Repo, MainBranch, FileExt, ReadMe) error
}

type ExploreApi interface {
	ExploreApiReader
	ExploreApiWriter
}

type Owner string

type Repo string

type Login string

type MainBranch string

type FileExt string

type RepoOwner struct {
	Login     string `json:"login"`
	Url       string `json:"url"`
	AvatarUrl string `json:"avatar_url"`
}

type Contribution struct {
	NameWithOwner string `json:"full_name"`
	Name          string `json:"name"`
	Url           string `json:"url"`
	Owner         RepoOwner
}

type Contributions struct {
	Name          string `json:"name"`
	Url           string `json:"url"`
	AvatarUrl     string `json:"avatar_url"`
	Contributions []Contribution
}

type RepoContributors struct {
	RepoContributors []Contributors
}

type Contributors struct {
	Login             string `json:"login"`
	ID                int64  `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarUrl         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_url"`
	Url               string `json:"url"`
	HtmlUrl           string `json:"html_url"`
	FollowersUrl      string `json:"followers_url"`
	FollowingUrl      string `json:"following_url"`
	GistsUrl          string `json:"gists_url"`
	StarredUrl        string `json:"starred_url"`
	SubscriptionsUrl  string `json:"subscriptions_url"`
	OrganizationsUrl  string `json:"organizations_url"`
	ReposUrl          string `json:"repos_url"`
	EventsUrl         string `json:"events_url"`
	ReceivedEventsUrl string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
	Contributions     int32  `json:"contributions"`
}

type ReadMe struct {
	Html string
}

type Health struct {
	Status string
}
