package main

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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
	}); err != nil {
		return err
	}

	return nil
}
