package main

import (
	"io/ioutil"
	"log"
	"strings"
)

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
