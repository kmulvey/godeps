package main

import (
	"io/ioutil"
	"log"
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

	return nil
}

func buildPatchedGoModFile() {
	input, err := ioutil.ReadFile("myfile")
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "]") {
			lines[i] = "LOL"
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile("myfile", []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}
