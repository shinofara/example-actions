package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/shinofara/example-actions/github"
	"github.com/shinofara/example-actions/policy"
)

var (
	githubAuthToken = os.Getenv("GITHUB_AUTH_TOKEN")
)

func main() {
	var orgName, userName string
	flag.StringVar(&orgName, "org", "", "your github org name")
	flag.StringVar(&userName, "user", "", "your github user name")
	flag.Parse()

	if orgName != "" && userName != "" {
		panic("one of org or user")
	}

	ctx := context.Background()

	var gc github.Client

	switch {
	case orgName != "" && userName != "":
		panic("one of org or user")
	case orgName != "":
		gc = github.NewOrgClient(ctx, githubAuthToken, orgName)
	case userName != "":
		gc = github.NewUserClient(ctx, githubAuthToken, userName)
	default:
		panic("require one of org or user")
	}

	re := policy.New(policy.NewOption(gc),
		policy.NewAccess(gc),
		policy.NewProtection(gc))

	allRepos, err := gc.GetAllReposoties(ctx)
	if err != nil {
		panic(err)
	}

	for _, repo := range allRepos {
		if repo.Archived != nil && *repo.Archived {
			continue
		}

		result, err := re.Do(ctx, repo)
		if err != nil {
			panic(err)
		}

		if result != nil {
			b, err := json.Marshal(result)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(b))
		}
	}
}
