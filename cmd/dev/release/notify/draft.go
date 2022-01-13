package notify

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/cobra"

	"github.com/ory/x/flagx"

	"github.com/ory/cli/cmd/dev/newsletter"
	"github.com/ory/cli/cmd/pkg"
)

var gitCommitMessageBaseRegex = regexp.MustCompile("(?im)^" + pkg.GitCommitMessagePreviousVersion + "\\sv([0-9a-zA-Z\\-\\._]+)$")

var draft = &cobra.Command{
	Use:   "draft [list-id]",
	Short: "Create a Mailchimp draft campaign for the release notification",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		gitHash := pkg.GitHubSHA()
		circleTag := pkg.GitHubTag()

		// Required by conventional-changelog-generator
		if _, err := os.Stat("package.json"); os.IsNotExist(err) {
			pkg.Check(ioutil.WriteFile("package.json",
				[]byte(`{"private": true, "version": "0.0.0"}`), 0600))
		}

		presetDir := pkg.GitClone("git@github.com:ory/changelog.git")
		pkg.Check(pkg.NewCommandIn(presetDir, "npm", "i").Run())

		commitMessage := pkg.CommandGetOutput("git", "log", "--format=%B", "-n", "1", gitHash)
		tagMessage := pkg.CommandGetOutput("git", "tag", "-l", "--format=%(contents)", circleTag)

		pkg.Check(pkg.NewCommand("npm", "--no-git-tag-version", "--allow-same-version", "version", circleTag).Run())

		changelogFile, err := ioutil.TempDir(os.TempDir(), "ory-release-cf-*")
		pkg.Check(err)
		changelogFile = path.Join(changelogFile, "changelog-email.md")

		count := 2
		if cliFromVersion := flagx.MustGetString(cmd, "from-version"); len(cliFromVersion) > 0 {
			cv, err := semver.StrictNewVersion(strings.TrimPrefix(cliFromVersion, "v"))
			pkg.Check(err)
			count = changelogGeneratorReleaseCount(cv, pkg.GitListTags())
		} else if cliFromTag, ok := getPreviousVersionFromGitCommitMessage(commitMessage); ok {
			count = changelogGeneratorReleaseCount(cliFromTag, pkg.GitListTags())
		}

		pkg.Check(pkg.NewCommand("npx", "conventional-changelog-cli@v2.1.1", "--config",
			path.Join(presetDir, "email.js"), "-r", strconv.Itoa(count), "-o", changelogFile).Run())

		pkg.Check(pkg.NewCommand("npx", "prettier", "-w", changelogFile).Run())

		changelog, err := ioutil.ReadFile(changelogFile)
		pkg.Check(err)

		if strings.TrimSpace(tagMessage) == strings.TrimSpace(commitMessage) {
			fmt.Println("Git tag does not include any release notes.")
			if strings.Contains(string(changelog), "no significant changes") {
				fmt.Println("Changelog would be empty, skipping campaign send!")
				return nil
			}
		}

		campaign, err := newsletter.Draft(
			args[0],
			flagx.MustGetInt(cmd, "segment"),
			[]byte(tagMessage),
			changelog,
		)
		pkg.Check(err)
		fmt.Printf(`Created campaign "%s" (%s)`, campaign.Settings.Title, campaign.ID)
		fmt.Println()

		fmt.Printf(`Campaign drafted with contents:

## Release Notes

%s

## Changelog

%s
`, tagMessage, changelog)

		return nil
	},
}

func init() {
	Main.AddCommand(draft)
	draft.Flags().Int("segment", 0, "The Mailchimp segment ID.")
	draft.Flags().String("from-version", "", "Use this as the previous version for changelog generation.")
}

func getPreviousVersionFromGitCommitMessage(message string) (*semver.Version, bool) {
	matches := gitCommitMessageBaseRegex.FindAllStringSubmatch(message, -1)
	if len(matches) != 1 {
		return nil, false
	}

	if len(matches[0]) != 2 {
		return nil, false
	}

	version, err := semver.StrictNewVersion(strings.TrimPrefix(matches[0][1], "v"))
	pkg.Check(err, `%s: "%s"`, err, matches[0][1])

	return version, true
}

// changelogGeneratorReleaseCount returns the `--release-count <count>` for `npx changelog-generator-cli`.
// The count works as follows:
//
// - `-r 1` all changes since the latest git tag. For fresh git tags without any newer commits this is always empty.
// - `-r 2` all changes since the latest git tag, plus the changes from the latest git tag.
func changelogGeneratorReleaseCount(tag *semver.Version, listOfTags string) int {
	var count = 0
	for _, line := range strings.Split(listOfTags, "\n") {
		if line == fmt.Sprintf("v%s", tag.String()) {
			count = 1
		} else if count > 0 && strings.TrimSpace(line) != "" {
			count++
		}
	}

	if count == 0 {
		return 2
	}

	return count + 1
}
