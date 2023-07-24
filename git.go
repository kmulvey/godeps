package godeps

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func createPrBranch(branchName string) error {
	var repo, err = git.PlainOpen("./")
	if err != nil {
		return err
	}

	headRef, err := repo.Head()
	if err != nil {
		return err
	}

	var prBrnach = plumbing.ReferenceName("refs/heads/" + branchName)
	ref := plumbing.NewHashReference(prBrnach, headRef.Hash())

	// The created reference is saved in the storage.
	if err := repo.Storer.SetReference(ref); err != nil {
		return err
	}

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	if err := w.Checkout(&git.CheckoutOptions{
		Branch: prBrnach,
		Keep:   true,
	}); err != nil {
		return err
	}

	return nil
}

func commitAndPush(dep Dependency) error {
	var repo, err = git.PlainOpen("./")
	if err != nil {
		return err
	}

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	_, err = w.Add("go.mod")
	if err != nil {
		return err
	}

	_, err = w.Add("go.sum")
	if err != nil {
		return err
	}

	commit, err := w.Commit(fmt.Sprintf("upgrade %s to %s", dep.Repo, dep.Version.String()), &git.CommitOptions{
		Author: &object.Signature{
			Name:  "GoDep",
			Email: "godep@godep.null",
			When:  time.Now(),
		},
	})

	_, err = repo.CommitObject(commit)
	if err != nil {
		return err
	}

	return repo.Push(&git.PushOptions{})
}
