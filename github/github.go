package github

import (
	"context"
	"net/http"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

const (
	perPage = 10
)

type Client interface {
	GetAllReposoties(ctx context.Context) ([]*github.Repository, error)
	GetListTeamsByRepo(ctx context.Context, repoName string) ([]*github.Team, error)
	GetBranchProtection(ctx context.Context, repoName string) *github.Protection
}

type BaseClient struct {
	Owner        string
	GithubClient *github.Client
}

func NewBaseClient(ctx context.Context, token string) *BaseClient {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})

	return &BaseClient{
		GithubClient: github.NewClient(oauth2.NewClient(ctx, ts)),
	}
}

type UserClient struct {
	*BaseClient
}

func NewUserClient(ctx context.Context, token string, user string) Client {
	b := NewBaseClient(ctx, token)
	b.Owner = user

	return &UserClient{
		BaseClient: b,
	}
}

func (c *UserClient) GetAllReposoties(ctx context.Context) ([]*github.Repository, error) {
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: perPage},
	}
	var allRepos []*github.Repository
	for {

		repos, resp, err := c.GithubClient.Repositories.List(ctx, c.Owner, opt)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allRepos, nil
}

type OrgClient struct {
	*BaseClient
}

func NewOrgClient(ctx context.Context, token string, user string) Client {
	b := NewBaseClient(ctx, token)
	b.Owner = user

	return &OrgClient{
		BaseClient: b,
	}
}

func (c *OrgClient) GetAllReposoties(ctx context.Context) ([]*github.Repository, error) {
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: perPage},
	}
	var allRepos []*github.Repository
	for {
		repos, resp, err := c.GithubClient.Repositories.ListByOrg(ctx, c.Owner, opt)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allRepos, nil
}

func (c *BaseClient) GetListTeamsByRepo(ctx context.Context, repoName string) ([]*github.Team, error) {
	teams, _, err := c.GithubClient.Repositories.ListTeams(ctx, c.Owner, repoName, nil)
	if err != nil {
		return nil, err
	}

	return teams, nil
}

func (c *BaseClient) GetBranchProtection(ctx context.Context, repoName string) *github.Protection {
	protection, resp, err := c.GithubClient.Repositories.GetBranchProtection(ctx, c.Owner, repoName, "main")
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			protection, _, _ = c.GithubClient.Repositories.GetBranchProtection(ctx, c.Owner, repoName, "master")
		}
	}

	if protection == nil {
		protection = &github.Protection{}
	}

	return protection
}
