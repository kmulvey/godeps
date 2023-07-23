package main

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

	var err = backupOriginalGoMod()
	if err != nil {
		return nil
	}

	return nil
}
