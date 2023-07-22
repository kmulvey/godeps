package main

import (
	"fmt"
	"os/exec"
)

func backupOriginalGoMod() error {

	var output, err = exec.Command("/bin/bash", "-c", "cp go.mod go.mod.original").Output()
	if err != nil {
		return err
	}

	fmt.Println(string(output))

	return nil
}
