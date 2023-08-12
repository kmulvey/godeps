package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/kmulvey/godeps"
	log "github.com/sirupsen/logrus"
)

func main() {
	var githubToken, owner, repo string

	// use cli params first
	flag.StringVar(&githubToken, "token", "", "your github token")
	flag.StringVar(&owner, "owner", "", "your github username")
	flag.StringVar(&repo, "repo", "", "repo name")
	flag.Parse()

	// env vars override cli params, because these are most likely to be used
	if githubRepo := os.Getenv("GITHUB_REPOSITORY"); githubRepo != "" {
		var repoArr = strings.Split(githubRepo, "/")
		if len(repoArr) == 2 {
			owner = repoArr[0]
			repo = repoArr[1]
		}
	}

	githubToken = os.Getenv("GITHUB_TOKEN")

	fmt.Printf("token: %s \n", githubToken)
	fmt.Printf("owner: %s \n", owner)
	fmt.Printf("repo: %s \n", repo)

	var err = godeps.Run(owner, repo, githubToken)
	if err != nil {
		log.Fatal(err)
	}
}
