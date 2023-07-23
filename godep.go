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
	return d.Repo + " " + d.Version.String()
}

type Upgrade struct {
	From Version
	To   Version
}

func Run(owner, repo, githubToken string) error {

	var githubClient = newGithubClient(githubToken)

	var err = backupOriginalGoMod()
	if err != nil {
		return nil
	}

	upgrades, err := findNewVersions()
	if err != nil {
		return nil
	}
	if len(upgrades) == 0 {
		return nil
	}

	existingPrs, err := getExistingPRs(githubClient, owner, repo)
	if err != nil {
		return nil
	}

	for repo, upgrade := range upgrades {

		var prTitle = fmt.Sprintf("Bump %s from %s to %s", repo, &upgrade.From, &upgrade.To)
		if _, exists := existingPrs[prTitle]; exists {
			continue
		}

		err = createUpgradePR(repo, owner, upgrade, githubClient)
		if err != nil {
			return nil
		}
	}

	return nil
}

func createUpgradePR(repo, owner string, upgrade Upgrade, githubClient *github.Client) error {

	var prTitle = fmt.Sprintf("Bump %s from %s to %s", repo, &upgrade.From, &upgrade.To)
	var prBranch = fmt.Sprintf("bump-%s-from-%s-to-%s", repo, &upgrade.From, &upgrade.To)
	prBranch = strings.ReplaceAll(prBranch, "/", "-")
	var newDep = Dependency{Repo: repo, Version: upgrade.To}

	if err := createPrBranch(prBranch); err != nil {
		return nil
	}

	if err := buildPatchedGoModFile(newDep); err != nil {
		return nil
	}

	if _, err := exec.Command("/bin/bash", "-c", "go mod tidy").Output(); err != nil {
		return err
	}

	if err := commitAndPush(newDep); err != nil {
		return err
	}

	if err := createPR(prTitle, "main", prBranch, owner, repo, githubClient); err != nil {
		return err
	}

	return nil
}
