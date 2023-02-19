package main // import "github.com/leighlondon/ghaudit"

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

type result struct {
	repo string
	user string
}

func main() {
	token, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok || token == "" {
		fmt.Fprintf(os.Stderr, "unable to find github api token in GITHUB_TOKEN\n")
		os.Exit(1)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	opts := github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 50},
		Affiliation: "owner",
	}
	var results []result
	for {
		repos, rsp, err := client.Repositories.List(ctx, "", &opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to list repos: %s\n", err)
			os.Exit(1)
		}
		for _, repo := range repos {
			name := repo.GetFullName()
			collabs, _, err := client.Repositories.ListCollaborators(ctx, *repo.Owner.Login, *repo.Name, nil)
			if err != nil {
				fmt.Fprintf(os.Stderr, "unable to list collabs for %s: %s\n", name, err)
				os.Exit(1)
			}
			for _, c := range collabs {
				if *c.Login == *repo.Owner.Login {
					continue
				}
				results = append(results, result{repo: name, user: *c.Login})
			}
		}

		if rsp.NextPage == 0 {
			break
		}
		opts.Page = rsp.NextPage
	}

	for _, r := range results {
		fmt.Fprintf(os.Stdout, "%s,%s\n", r.repo, r.user)
	}
}
