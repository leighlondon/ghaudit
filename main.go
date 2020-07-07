package main // import "github.com/leighlondon/ghaudit"

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

func main() {
	token, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok || token == "" {
		fmt.Fprintf(os.Stderr, "unable to find github api token\n")
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
	var results []*github.Repository
	for {

		repos, rsp, err := client.Repositories.List(ctx, "", &opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to list repos: %s\n", err)
			os.Exit(1)
		}
		results = append(results, repos...)
		if rsp.NextPage == 0 {
			break
		}
		opts.Page = rsp.NextPage
	}
	for _, r := range results {
		fmt.Fprintf(os.Stdout, "%s\n", r.GetFullName())
	}
}
