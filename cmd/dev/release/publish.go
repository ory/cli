package release

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/ory/x/flagx"
	"github.com/ory/x/stringslice"
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/pkg"
)

var supported = []string{"hydra", "kratos", "keto", "oathkeeper", "cli"}

var publish = &cobra.Command{
	Use:   "publish [version]",
	Args:  cobra.ExactArgs(1),
	Short: "Publish a new release",
	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		pkg.Check(err)
		project := path.Base(wd)
		if !stringslice.Has(supported, project) {
			pkg.Fatalf(`This script is expected to run in a directory named "hydra", "keto", "oathkeeper", "kratos".`)
			return
		}

		dry := flagx.MustGetBool(cmd, "dry")
		gitCleanTags()

		currentVersion, err := semver.StrictNewVersion(strings.TrimPrefix(pkg.GitGetCurrentTag(),"v"))
		pkg.Check(err, "Unable to parse current git tag %s: %s", pkg.GitGetCurrentTag(), err)

		var nextVersion semver.Version
		switch args[0] {
		case "major":
			nextVersion = currentVersion.IncMajor()
		case "minor":
			nextVersion = currentVersion.IncMinor()
		case "patch":
			nextVersion = currentVersion.IncPatch()
		default:
			nv, err := semver.StrictNewVersion(strings.TrimPrefix(args[0],"v"))
			pkg.Check(err)
			nextVersion = *nv
		}

		knowsTag(&nextVersion)

		pkg.Check(pkg.NewCommand("goreleaser", "check").Run())
		pkg.Check(pkg.NewCommand("circleci", "config", "check").Run())

		if dry {
			fmt.Println("Don't worry, this is a dry run!")
		}
		pkg.Confirm("Are you sure you want to bump to v%s? Previous version was v%s.", nextVersion, currentVersion)

		switch project {
		case "hydra":
			pkg.Confirm("This will also release hydra-login-consent-node:v%s. Previous version was v%s. Is that ok?", nextVersion, currentVersion)
			pkg.GitTagRelease(pkg.GitClone("git@github.com:ory/hydra-login-consent-node.git"), false, dry, nextVersion, nil)
		case "kratos":
			pkg.Confirm("This will also release kratos-selfservice-ui-node:v%s. Previous version was v%s. Is that ok?", nextVersion, currentVersion)
			pkg.GitTagRelease(pkg.GitClone("git@github.com:ory/kratos-selfservice-ui-node.git"), false, dry, nextVersion, nil)
		}

		var fromVersion *semver.Version
		if ov := flagx.MustGetString(cmd, "from-version"); len(ov) > 0 {
			fromVersion, err = semver.StrictNewVersion(strings.TrimPrefix(ov,"v"))
			pkg.Check(err, "Unable to parse from-version git tag v%s: %s", ov, err)
			knowsTag(fromVersion)
		}
		pkg.GitTagRelease(wd, true, dry, nextVersion, fromVersion)

		fmt.Printf("Successfully released version: v%s", nextVersion.String())
	},
}

func knowsTag(v *semver.Version) {
	if stringslice.Has(strings.Split(pkg.GitListTags(), "\n"), fmt.Sprintf("v%s",v)) {
		pkg.Fatalf(`Version v%s exists already and can not be re-released!`, v.String())
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
	publish.Flags().String("from-version", "", "When set includes all release up to this release in the changelog that will be sent out.")
}
