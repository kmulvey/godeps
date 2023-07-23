package main

import (
	"flag"

	"github.com/kmulvey/godeps"
	log "github.com/sirupsen/logrus"
)

func main() {

	var githubToken, owner, repo string
	flag.StringVar(&githubToken, "token", "", "your github token")
	flag.StringVar(&owner, "owner", "", "your github username")
	flag.StringVar(&repo, "repo", "", "repo name")
	flag.Parse()

	var err = godeps.Run(owner, repo, githubToken)
	if err != nil {
		log.Fatal(err)
	}
}
