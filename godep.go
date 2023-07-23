package main

import (
	"fmt"
	"time"
)

type Dependency struct {
	Repo string
	Version
}

func (d *Dependency) String() string {
	return d.Repo + " " + d.Version.String()
}

type Version struct {
	Major uint8
	Minor uint8
	Patch uint8
	Date  time.Time
	Sha   string
}

func (v *Version) String() string {
	if v.Sha != "" {
		return fmt.Sprintf("%d.%d.%d-%s-%s", v.Major, v.Minor, v.Patch, v.Date.Format("20060102150405"), v.Sha)
	}
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

type Upgrade struct {
	From Version
	To   Version
}

func Run(owner, repo, githubToken string) error {

	var err = backupOriginalGoMod()
	if err != nil {
		return nil
	}

	return nil
}
