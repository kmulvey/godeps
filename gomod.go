package main

import (
	"os/exec"
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

	return nil
}
