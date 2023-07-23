package main

import (
	"io/ioutil"
	"os/exec"
	"strings"
)

func backupOriginalGoMod() error {

	var _, err = exec.Command("/bin/bash", "-c", "cp go.mod go.mod.original").Output()
	if err != nil {
		return err
	}

	_, err = exec.Command("/bin/bash", "-c", "go get -u -v ./...").Output()
	if err != nil {
		return err
	}

	_, err = exec.Command("/bin/bash", "-c", "go mod tidy").Output()
	if err != nil {
		return err
	}

	_, err = exec.Command("/bin/bash", "-c", "git diff go.mod > diff").Output()
	if err != nil {
		return err
	}

	_, err = exec.Command("/bin/bash", "-c", "git checkout go.mod go.sum").Output()
	if err != nil {
		return err
	}

	return nil
}

func buildPatchedGoModFile(dep Dependency) error {
	input, err := ioutil.ReadFile("go.mod")
	if err != nil {
		return err
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, dep.Repo) {
			lines[i] = dep.String()
		}
	}
	output := strings.Join(lines, "\n")
	return ioutil.WriteFile("go.mod", []byte(output), 0644)
}
