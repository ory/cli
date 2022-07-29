package release

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Masterminds/semver/v3"

	"github.com/spf13/cobra"

	"github.com/ory/x/flagx"
	"github.com/ory/x/stringslice"

	"github.com/ory/cli/cmd/pkg"
)

var isTestRelease = regexp.MustCompile(`^(([a-zA-Z0-9\.\-]+\.)|)pre\.[0-9]+$`)

var publish = &cobra.Command{
	Use:   "publish [version]",
	Args:  cobra.ExactArgs(1),
	Short: "Publish a new release",
	Long: `Performs git magic and other automated tasks such as tagging the example applications for ORY Kratos and ORY Hydra as well.

To publish a release, you first have to create a pre-release:

	ory dev release publish v0.1.0-pre.0

Once the release pipeline finished successfully (skip publishing the newsletter!) you can run:

	ory dev release publish v0.1.0 --include-changelog-since v0.1.0-pre.0

Here are some examples how to choose pre:

- v0.1.0-pre.0 -> v0.1.0
- v0.1.0-pre.1 -> v0.1.0
- v0.1.0-alpha.0.pre.0 -> v0.1.0-alpha.0
- v0.1.0-alpha.1.pre.0 -> v0.1.0-alpha.1
- v0.1.0-rc.1.pre.0 -> v0.1.0-rc.1

In case where the release pipeline failed and you re-create another release where you want to include the changelog from the failed release, perform the following:

1. Assuming release "v0.1.0" failed
2. You wish to create "v0.1.1" and include the changelog of "v0.1.0" as well
3. Run ` + "`ory dev release publish v0.1.1 --include-changelog-since v0.1.0`",
	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		pkg.Check(err)

		cfg, err := pkg.ReadConfig()
		pkg.Check(err)

		dry := flagx.MustGetBool(cmd, "dry")
		gitCleanTags()

		var latestTag string
		if cfg.IgnoreTags != nil && len(cfg.IgnoreTags.String()) > 0 {
			for {
				var o, e bytes.Buffer
				args := []string{"describe", "--abbrev=0", "--tags"}
				if latestTag != "" {
					args = append(args, latestTag+"^")
				}
				cmd := pkg.NewCommand("git", args...)
				cmd.Stdout = &o
				cmd.Stderr = &e
				if cmd.Run() != nil {
					pkg.Fatalf("could not get git tag: %s%s", o.String(), e.String())
				}
				latestTag = strings.TrimSpace(o.String())

				if !cfg.IgnoreTags.MatchString(latestTag) {
					break
				}
			}
		} else {
			latestTag = pkg.GitGetCurrentTag()
		}

		currentVersion, err := semver.StrictNewVersion(strings.TrimPrefix(latestTag, "v"))
		pkg.Check(err, "Unable to parse current git tag %s: %s", latestTag, err)

		var nextVersion semver.Version
		switch args[0] {
		case "major":
			nextVersion = currentVersion.IncMajor()
		case "minor":
			nextVersion = currentVersion.IncMinor()
		case "patch":
			nextVersion = currentVersion.IncPatch()
		default:
			nv, err := semver.StrictNewVersion(strings.TrimPrefix(args[0], "v"))
			pkg.Check(err)
			nextVersion = *nv
		}

		checkForDuplicateTag(&nextVersion)

		if !isTestRelease.MatchString(currentVersion.Prerelease()) &&
			!isTestRelease.MatchString(nextVersion.Prerelease()) {
			pkg.Confirm(`You should create a test release before publishing the real release or vice versa:

- Current version: 	v%s
- Next version: 	v%s

Please check "ory help dev release publish".

Are you sure you want to proceed without creating a pre version first?`, currentVersion, nextVersion)
		}

		if ov := flagx.MustGetString(cmd, "include-changelog-since"); len(ov) == 0 && !isTestRelease.MatchString(nextVersion.Prerelease()) {
			pkg.Confirm("You are about to release a non-test release v%s but did not include the --include-changelog-since flag. Are you sure you want to continue?", nextVersion)
		}

		pkg.Check(pkg.NewCommand("goreleaser", "check").Run())

		if dry {
			fmt.Println("Don't worry, this is a dry run!")
		}
		pkg.Confirm("Are you sure you want to bump to v%s? Previous version was v%s.", nextVersion, currentVersion)

		if len(cfg.PreReleaseHooks) != 0 {
			fmt.Println("Running pre-release hooks...")

			for _, h := range cfg.PreReleaseHooks {
				parts := strings.Split(h, " ")

				var err error
				if len(parts) > 1 {
					err = pkg.NewCommand(parts[0], parts[1:]...).Run()
				} else {
					err = pkg.NewCommand(parts[0]).Run()
				}

				if err != nil {
					pkg.Fatalf("Pre-release hook failed: %s\nAborting release.", err.Error())
				}
			}
		}

		var fromVersion *semver.Version
		if ov := flagx.MustGetString(cmd, "include-changelog-since"); len(ov) > 0 {
			fromVersion, err = semver.StrictNewVersion(strings.TrimPrefix(ov, "v"))
			pkg.Check(err, "Unable to parse include-changelog-since git tag v%s: %s", ov, err)
			checkIfTagExists(fromVersion)
		}

		pkg.GitTagRelease(wd, !isTestRelease.MatchString(nextVersion.Prerelease()), dry, nextVersion, fromVersion)

		switch cfg.Project {
		case "hydra":
			pkg.GitTagRelease(pkg.GitClone("git@github.com:ory/hydra-login-consent-node.git"), false, dry, nextVersion, nil)
		case "kratos":
			pkg.GitTagRelease(pkg.GitClone("git@github.com:ory/kratos-selfservice-ui-node.git"), false, dry, nextVersion, nil)
			pkg.GitTagRelease(pkg.GitClone("git@github.com:ory/kratos-selfservice-ui-react-native.git"), false, dry, nextVersion, nil)
		}

		fmt.Printf("Successfully released version: v%s\n", nextVersion.String())
	},
}

func checkForDuplicateTag(v *semver.Version) {
	if stringslice.Has(strings.Split(pkg.GitListTags(), "\n"), fmt.Sprintf("v%s", v)) {
		pkg.Fatalf(`Version v%s exists already and can not be re-released!`, v.String())
	}
}

func checkIfTagExists(v *semver.Version) {
	if !stringslice.Has(strings.Split(pkg.GitListTags(), "\n"), fmt.Sprintf("v%s", v)) {
		pkg.Fatalf(`Version v%s does not exist!`, v.String())
	}
}

func gitCleanTags() {
	pkg.Check(pkg.NewCommand("git", "checkout", "master").Run())
	pkg.Check(pkg.NewCommand("git",
		append([]string{"tag", "-d"}, pkg.BashPipe(pkg.GitListTags())...)...).Run())
	pkg.Check(pkg.NewCommand("git", "fetch", "origin", "--tags").Run())
	pkg.Check(pkg.NewCommand("git", "pull", "-ff").Run())
	pkg.Check(pkg.NewCommand("git", "diff", "--exit-code").Run())
}

func init() {
	Main.AddCommand(publish)
	publish.Flags().Bool("dry", false, "Make changes only locally and do not push to remotes.")
	publish.Flags().String("include-changelog-since", "", "If set includes all changelog entries for all git tags up to and including the specified git tag")
}
