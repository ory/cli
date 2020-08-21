package monorepo

import (
	"os/exec"
	"strings"
)

// getRepositoryChanges returns the changes in the specified local respository (via rootDirectory) towards the specified
func getRepositoryChanges(rootDirectory string, parent string, gitOpts []string) (string, error) {
	args := []string{"--no-pager", "log", parent}
	if gitOpts == nil {
		//apply default args
		gitOpts = []string{"--name-only", "--oneline"}
	}

	cmd := exec.Command("git", append(args, gitOpts...)...)
	cmd.Dir = rootDirectory
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out[:]), nil
}

// cleanseRepositoryChanges currently support two cleansing operations.
// * lines []string representing the change list as produced by getRepositoryChanges
// * includeFiles if true, the returned output will include a list of all changed files (with the relative path)
func cleanseRepositoryChanges(lines *[]string, includeFiles bool, deduplicate bool) {

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
	if deduplicate {
		deduplicateChanges(lines)
	}

}

func deduplicateChanges(lines *[]string) {
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
