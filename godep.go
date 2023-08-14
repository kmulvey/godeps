package godeps

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/google/go-github/v53/github"
)

type Dependency struct {
	Repo string
	Version
}

func (d *Dependency) String() string {
	return d.Repo + " v" + d.Version.String()
}

type Upgrade struct {
	From Version
	To   Version
}

func Run(owner, repo, githubToken string) error {

	var githubClient = newGithubClient(githubToken)

	var err = backupOriginalGoMod()
	if err != nil {
		return fmt.Errorf("backupOriginalGoMod: %w", err)
	}

	upgrades, err := findNewVersions()
	if err != nil {
		return fmt.Errorf("findNewVersions: %w", err)
	}
	if len(upgrades) == 0 {
		return nil
	}

	existingPrs, err := getExistingPRs(githubClient, owner, repo)
	if err != nil {
		return fmt.Errorf("getExistingPRs: %w", err)
	}

	for depRepo, upgrade := range upgrades {

		var prTitle = fmt.Sprintf("Bump %s from %s to %s", depRepo, &upgrade.From, &upgrade.To)
		if _, exists := existingPrs[prTitle]; exists {
			continue
		}

		err = createUpgradePR(depRepo, repo, owner, upgrade, githubClient)
		if err != nil {
			return fmt.Errorf("createUpgradePR: %w", err)
		}
	}

	return nil
}

func createUpgradePR(depRepo, thisRepo, repoOwner string, upgrade Upgrade, githubClient *github.Client) error {

	var ctx = context.Background()
	var prTitle = fmt.Sprintf("Bump %s from %s to %s", depRepo, &upgrade.From, &upgrade.To)
	var prBranch = fmt.Sprintf("bump-%s-from-%s-to-%s", depRepo, &upgrade.From, &upgrade.To)
	prBranch = strings.ReplaceAll(prBranch, "/", "-")
	var newDep = Dependency{Repo: depRepo, Version: upgrade.To}

	var mainBranch, err = getMainBranch(ctx, githubClient, repoOwner, thisRepo)
	if err != nil {
		log.Fatalf("Unable to main branch: %s\n", err)
	}

	ref, err := getRef(ctx, githubClient, prBranch, mainBranch, repoOwner, thisRepo)
	if err != nil {
		log.Fatalf("Unable to get/create the commit reference: %s\n", err)
	}
	if ref == nil {
		log.Fatalf("No error where returned but the reference is nil")
	}

	if err := buildPatchedGoModFile(newDep); err != nil {
		return fmt.Errorf("buildPatchedGoModFile: %w", err)
	}

	if _, err := exec.Command("/bin/bash", "-c", "go mod tidy").Output(); err != nil {
		return fmt.Errorf("go mod tidy: %w", err)
	}

	tree, err := getTree(ctx, githubClient, ref, repoOwner, thisRepo)
	if err != nil {
		log.Fatalf("Unable to create the tree based on the provided files: %s\n", err)
	}

	if err := pushCommit(ctx, githubClient, ref, tree, newDep, repoOwner, thisRepo); err != nil {
		log.Fatalf("Unable to create the commit: %s\n", err)
	}

	if err := createPR(ctx, githubClient, thisRepo, repoOwner, prBranch, mainBranch, prTitle, ""); err != nil { // fill description
		log.Fatalf("Error while creating the pull request: %s", err)
	}

	return nil
}
