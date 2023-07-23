package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

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

func getExistingPRs(client *github.Client, repoOwner, repo string) (map[string]struct{}, error) {

	prs, _, err := client.PullRequests.List(context.Background(), repoOwner, repo, nil)
	if err != nil {
		return nil, err
	}

	var existingPRs = make(map[string]struct{})
	for _, pr := range prs {
		existingPRs[*pr.Title] = struct{}{}
	}

	return existingPRs, nil
}
