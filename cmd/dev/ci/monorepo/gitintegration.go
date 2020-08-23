package monorepo

import (
	"os/exec"
	"regexp"
	"sort"
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
		//fmt.Printf("%d. [%s]\n", i, line)
		if strings.Contains(line, "/") {
			//if line contains '/', it represent a path with filename
			//we remove everything starting from the last '/'
			if !includeFiles {
				(*lines)[i] = line[:strings.LastIndex(line, "/")]
			}
		} else {
			//if line does not include any '/', it represent filename in the root directory
			//we substitute the filename with '.'
			if !includeFiles {
				(*lines)[i] = "."
			}
		}
	}
	if deduplicate {
		deduplicateChanges(lines)
	}
}

func caseInsensitiveSort(data []string) []string {
	sort.Slice(data, func(i, j int) bool { return strings.ToLower(data[i]) < strings.ToLower(data[j]) })
	return data
}

func removeCommitMessages(changeLog string) string {
	regex := regexp.MustCompile(`(?m)^([a-z0-9]{7} .*\s|\z)|\n$`)
	cleansedChangeLog := regex.ReplaceAllString(changeLog, "")
	return cleansedChangeLog
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
