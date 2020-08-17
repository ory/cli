package monorepo

import (
	"fmt"
	"os/exec"
	"strings"
)

func changes(rootDirectory string, parent string) (string, error) {
	//cmd := exec.Command("git", "--no-pager", "log") //, "origin/master", "--name-only", "--oneline", "|", "sed", "'/ /d'", "|", "sed", `/\//!d'`, "|", "sed", `'s/\/.*//'`, "|", "sort", "|", "uniq")
	cmd := exec.Command("git", "--no-pager", "log", parent, "--name-only", "--oneline")
	cmd.Dir = rootDirectory
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out[:]), nil
}

func cleanseChangeList(lines *[]string, includeFiles bool) {
	fmt.Printf("cleanseChangeList: includeFiles '%t'", includeFiles)
	for i, line := range *lines {
		if strings.Contains(line, "/") {
			//if line contains '/', it represent a path with filename
			//we remove everything starting from the last '/'
			if !includeFiles {
				(*lines)[i] = line[:strings.LastIndex(line, "/")]
			}
		} else {
			//if line does not any '/', it represent filename in the root directory
			//we substitute the filename with '.'
			(*lines)[i] = "."
		}
	}
}

func removeDuplicateLines(lines *[]string) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *lines {
		if !found[x] {
			found[x] = true
			(*lines)[j] = (*lines)[i]
			j++
		}
	}
	*lines = (*lines)[:j]
}
