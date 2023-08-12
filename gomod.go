package godeps

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func backupOriginalGoMod() error {

	var _, err = exec.Command("/bin/bash", "-c", "cp go.mod go.mod.original").Output()
	if err != nil {
		return fmt.Errorf("cp go.mod go.mod.original: %w", err)
	}

	_, err = exec.Command("/bin/bash", "-c", "go get -u -v ./...").Output()
	if err != nil {
		return fmt.Errorf("go get -u -v ./...: %w", err)
	}

	_, err = exec.Command("/bin/bash", "-c", "go mod tidy").Output()
	if err != nil {
		return fmt.Errorf("go mod tidy: %w", err)
	}

	_, err = exec.Command("/bin/bash", "-c", "git diff go.mod > diff").Output()
	if err != nil {
		return fmt.Errorf("git diff go.mod > diff: %w", err)
	}

	_, err = exec.Command("/bin/bash", "-c", "git checkout go.mod go.sum").Output()
	if err != nil {
		return fmt.Errorf("git checkout go.mod go.sum: %w", err)
	}

	return nil
}

func buildPatchedGoModFile(dep Dependency) error {
	input, err := os.ReadFile("go.mod.original")
	if err != nil {
		return fmt.Errorf("read go.mod.original: %w", err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, dep.Repo) {
			lines[i] = "\t" + dep.String()
		}
	}
	output := strings.Join(lines, "\n")

	if err := os.WriteFile("go.mod", []byte(output), 0644); err != nil {
		return fmt.Errorf("write go.mod: %w", err)
	}

	return nil
}
