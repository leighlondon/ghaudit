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
		fmt.Fprintf(os.Stderr, "unable to find github api token")
		os.Exit(1)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	repos, _, err := client.Repositories.List(ctx, "", nil)
}
