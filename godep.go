package godeps

import (
	"fmt"
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

func createUpgradePR(depRepo, thisRepo, owner string, upgrade Upgrade, githubClient *github.Client) error {

	var prTitle = fmt.Sprintf("Bump %s from %s to %s", depRepo, &upgrade.From, &upgrade.To)
	var prBranch = fmt.Sprintf("bump-%s-from-%s-to-%s", depRepo, &upgrade.From, &upgrade.To)
	prBranch = strings.ReplaceAll(prBranch, "/", "-")
	var newDep = Dependency{Repo: depRepo, Version: upgrade.To}

	if err := createPrBranch(prBranch); err != nil {
		return fmt.Errorf("createPrBranch: %w", err)
	}

	if err := buildPatchedGoModFile(newDep); err != nil {
		return fmt.Errorf("buildPatchedGoModFile: %w", err)
	}

	if _, err := exec.Command("/bin/bash", "-c", "go mod tidy").Output(); err != nil {
		return fmt.Errorf("go mod tidy: %w", err)
	}

	if err := commitAndPush(newDep); err != nil {
		return fmt.Errorf("commitAndPush: %w", err)
	}

	if err := createPR(prTitle, "main", prBranch, owner, thisRepo, githubClient); err != nil {
		return fmt.Errorf("createPR: %w", err)
	}

	return nil
}
