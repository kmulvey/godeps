package godeps

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

var (
	authorName  = "Godep"
	authorEmail = "godep@godep.null"
)

func newGithubClient(accessToken string) *github.Client {
	var ctx = context.Background()
	var ts = oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	var tc = oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

// getMainBranch returns the main branch (main or master)
func getMainBranch(ctx context.Context, client *github.Client, repoOwner, repoName string) (string, error) {

	// return it if it already exists
	var _, _, err = client.Git.GetRef(ctx, repoOwner, repoName, "refs/heads/main")
	if err == nil {
		return "main", nil
	} else if _, _, err := client.Git.GetRef(ctx, repoOwner, repoName, "refs/heads/master"); err == nil {
		return "master", nil
	}

	return "", err
}

// getRef returns the commit branch reference object if it exists or creates it
// from the base branch before returning it.
func getRef(ctx context.Context, client *github.Client, prBranch, mainBranch, repoOwner, repoName string) (ref *github.Reference, err error) {

	// return it if it already exists
	if ref, _, err = client.Git.GetRef(ctx, repoOwner, repoName, "refs/heads/"+prBranch); err == nil {
		return ref, nil
	}

	// get a ref to main
	var baseRef *github.Reference
	if baseRef, _, err = client.Git.GetRef(ctx, repoOwner, repoName, "refs/heads/"+mainBranch); err != nil {
		return nil, err
	}

	// pr ref
	newRef := &github.Reference{Ref: github.String("refs/heads/" + prBranch), Object: &github.GitObject{SHA: baseRef.Object.SHA}}
	ref, _, err = client.Git.CreateRef(ctx, repoOwner, repoName, newRef)
	return ref, err
}

// getTree generates the tree to commit based on the given files and the commit
// of the ref you got in getRef.
func getTree(ctx context.Context, client *github.Client, ref *github.Reference, repoOwner, repoName string) (tree *github.Tree, err error) {
	// Create a tree with what to commit.
	var entries = make([]*github.TreeEntry, 1)

	// Load each file into the tree.
	file, content, err := getFileContent("commitfile")
	if err != nil {
		return nil, err
	}
	entries[0] = &github.TreeEntry{
		Path:    github.String(file),
		Type:    github.String("blob"),
		Content: github.String(string(content)),
		Mode:    github.String("100644"),
		Size:    github.Int(len(content)),
	}

	tree, _, err = client.Git.CreateTree(ctx, repoOwner, repoName, *ref.Object.SHA, entries)
	return tree, err
}

// getFileContent loads the local content of a file and return the target name
// of the file in the target repository and its contents.
func getFileContent(fileArg string) (targetName string, b []byte, err error) {
	var localFile string
	files := strings.Split(fileArg, ":")
	switch {
	case len(files) < 1:
		return "", nil, errors.New("empty `-files` parameter")
	case len(files) == 1:
		localFile = files[0]
		targetName = files[0]
	default:
		localFile = files[0]
		targetName = files[1]
	}

	b, err = os.ReadFile(localFile)
	return targetName, b, err
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

// pushCommit creates the commit in the given reference using the given tree.
func pushCommit(ctx context.Context, client *github.Client, ref *github.Reference, tree *github.Tree, dep Dependency, repoOwner, repoName string) (err error) {
	// Get the parent commit to attach the commit to.
	parent, _, err := client.Repositories.GetCommit(ctx, repoOwner, repoName, *ref.Object.SHA, nil)
	if err != nil {
		return err
	}
	// This is not always populated, but is needed.
	parent.Commit.SHA = parent.SHA

	// Create the commit using the tree.
	date := time.Now()
	author := &github.CommitAuthor{Date: &github.Timestamp{Time: date}, Name: &authorName, Email: &authorEmail}
	var commitMessage = fmt.Sprintf("upgrade %s to %s", dep.Repo, dep.Version.String())
	commit := &github.Commit{Author: author, Message: &commitMessage, Tree: tree, Parents: []*github.Commit{parent.Commit}}
	newCommit, _, err := client.Git.CreateCommit(ctx, repoOwner, repoName, commit)
	if err != nil {
		return err
	}

	// Attach the commit to the master branch.
	ref.Object.SHA = newCommit.SHA
	_, _, err = client.Git.UpdateRef(ctx, repoOwner, repoName, ref, false)
	return err
}

// createPR creates a pull request. Based on: https://godoc.org/github.com/google/go-github/github#example-PullRequestsService-Create
func createPR(ctx context.Context, client *github.Client, repoName, repoOwner, prBranch, mainBranch, prSubject, prDescription string) error {

	newPR := &github.NewPullRequest{
		Title:               &prSubject,
		Head:                &prBranch,
		Base:                &mainBranch,
		Body:                &prDescription,
		MaintainerCanModify: github.Bool(true),
	}

	var _, _, err = client.PullRequests.Create(ctx, repoOwner, repoName, newPR)
	if err != nil {
		return err
	}

	return nil
}
