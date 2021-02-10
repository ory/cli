package monorepo

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var changesMode string
var changeLog string
var gitOpts string

const allChanges = "--name-only --oneline"
const lastCommitChanges = "-1 --name-only --oneline"

var changes = &cobra.Command{
	Use:   "changes",
	Short: "List changes in the repository.",
	Long:  `List changes in the repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		repoChanges := "[empty]"
		switch changesMode {
		case "full":
			if len(gitOpts) == 0 {
				if isPR() {
					gitOpts = "--pretty=full"
				} else {
					gitOpts = "-1 --pretty=full"
				}
			}

			repoChanges, _ = getChangeLog(rootDirectory, revisionRange, gitOpts)
		case "directories":
			repoChanges, _ = getChangedDirectories(rootDirectory, revisionRange, gitOpts)
		case "files":
			repoChanges, _ = getChangedFiles(rootDirectory, revisionRange, gitOpts)
		default:
			log.Fatalf("Unknown ListMode '%s'", changesMode)
		}
		fmt.Println(repoChanges)
	},
}

func getGitDefaults() string {
	if isPR() {
		return allChanges
	}
	return lastCommitChanges
}

func getChangedFiles(rootDirectory string, revisionRange string, gitOpts string) (string, error) {
	changeLog, err := getRepositoryChangeLog(rootDirectory, revisionRange, gitOpts)
	if err != nil {
		return "", fmt.Errorf("Error getting changes from Git: %v", err)
	}
	cleansedChangeLogArray := strings.Split(removeCommitMessages(changeLog), "\n")
	cleanseRepositoryChanges(&cleansedChangeLogArray, true, true)
	caseInsensitiveSort(cleansedChangeLogArray)

	changedFiles := strings.Join(cleansedChangeLogArray, "\n")
	if debug {
		fmt.Printf("getChangedFiles: \n%s\n", changedFiles)
	}
	return changedFiles, nil
}

func getChangedDirectories(rootDirectory string, revisionRange string, gitOpts string) (string, error) {
	changeLog, err := getRepositoryChangeLog(rootDirectory, revisionRange, gitOpts)
	if err != nil {
		return "", fmt.Errorf("Error getting changes from Git: %v", err)
	}
	cleansedChangeLogArray := strings.Split(removeCommitMessages(changeLog), "\n")
	cleanseRepositoryChanges(&cleansedChangeLogArray, false, true)
	caseInsensitiveSort(cleansedChangeLogArray)
	changeDirectoriesString := strings.Join(cleansedChangeLogArray, "\n")
	if debug {
		fmt.Printf("getChangedDirectories: \n%s\n\n", changeDirectoriesString)
	}
	return changeDirectoriesString, nil
}

func getChangeLog(rootDirectory string, revisionRange string, gitOpts string) (string, error) {
	repoChanges, err := getRepositoryChangeLog(rootDirectory, revisionRange, gitOpts)
	if err != nil {
		return "", fmt.Errorf("Error getting changes from Git: %v", err)
	}

	return repoChanges, nil
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
		deduplicateChangelog(lines)
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

func deduplicateChangelog(lines *[]string) {
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

// getRepositoryChanges returns the changes in the specified local respository (via repositoryPath) towards the specified
func getRepositoryChangeLog(repositoryPath string, revisionRange string, gitOptions string) (string, error) {
	//TODO: caching the changelog, ensure this does not lead to problems.
	if len(changeLog) == 0 {
		args := []string{"--no-pager", "log"}
		if len(revisionRange) > 0 {
			args = append(args, revisionRange)
		}
		if len(gitOptions) == 0 {
			//apply default args
			gitOptions = getGitDefaults()
		}
		gitOptionsArray := strings.Split(gitOptions, " ")
		args = append(args, gitOptionsArray...)
		if verbose {
			fmt.Printf("getRepositoryChanges: '$ git %s'\n", args)
		}
		cmd := exec.Command("git", args...)
		cmd.Dir = repositoryPath

		out, err := cmd.Output()
		if err != nil {
			return "", err
		}
		changeLog = string(out[:])
	}
	if debug {
		fmt.Printf("getRepositoryChangeLog: \n%s\n", changeLog)
	}
	return changeLog, nil
}

func init() {
	Main.AddCommand(changes)
	changes.Flags().StringVarP(&changesMode, "mode", "m", "directories", "Define which which type of change information you want to get listed (full, files, directories). Default is 'directories'.")
	changes.Flags().StringVarP(&gitOpts, "gitopts", "g", "", "Specify custom git arguments used to determine changes, e.g. '--pretty=oneline' Only supported if mode is 'full'.")
}
