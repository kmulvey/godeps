package godeps

import (
	"io/ioutil"
	"strings"
)

func patchGoMod(depFile string, dep Dependency) error {

	var input, err = ioutil.ReadFile(depFile)
	if err != nil {
		return err
	}

	var lines = strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, dep.Repo) {
			lines[i] = dep.String()
		}
	}

	var output = strings.Join(lines, "\n")
	return ioutil.WriteFile(depFile, []byte(output), 0644)
}
