package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

// createPR creates a pull request. Based on: https://godoc.org/github.com/google/go-github/github#example-PullRequestsService-Create
func createPR(title, baseBranch, prBranch, owner, repo string, client *github.Client) (err error) {

	newPR := &github.NewPullRequest{
		Title:               &title,
		Head:                &prBranch,
		Base:                &baseBranch,
		Body:                &title, // TODO
		MaintainerCanModify: github.Bool(true),
	}

	pr, _, err := client.PullRequests.Create(context.Background(), owner, repo, newPR)
	if err != nil {
		return err
	}

	fmt.Printf("PR created: %s\n", pr.GetHTMLURL())
	return nil
}

func newGithubClient(accessToken string) *github.Client {
	var ctx = context.Background()
	var ts = oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	var tc = oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func getPRs(client *github.Client) ([]string, error) {

	var prs, _, err = client.PullRequests.List(context.Background(), "kmulvey", "text2speech", nil)
	if err != nil {
		return nil, err
	}

	var prTitles = make([]string, len(prs))
	for i, pr := range prs {
		prTitles[i] = *pr.Title
	}

	return prTitles, nil
}
