package monorepo

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var listMode string
var parent string
var includeFiles bool

var sampleOutput = `303599b feat: improve create migrations command (#16)
cmd/dev/pop/migration/create.go
ca21cb5 autogen: pin v0.0.24 release commit
232143e autogen: pin v0.0.24 release commit
97ba127 autogen: pin v0.0.24 release commit
0845e2b fix: register schema command to dev parent command (#15)
cmd/dev/main.go
5e47981 autogen: pin v0.0.23 release commit
f0a888d ci: bump orbs
.circleci/config.yml
0d6d4ed fix: use tap for brew instruction
.goreleaser.yml
662216d ci: remove test action
.github/workflows/test.yml
41ae16f ci: remove checks-go action
.github/workflows/checks-go.yml
ac0326d fix: revert brews.0.tap to brews.0.github
.goreleaser.yml
6404955 autogen: pin v0.0.22 release commit
9d51695 autogen: pin v0.0.21 release commit
bd49b2d feat: add render version json schema command (#13)
.schema/version_meta.schema.json
cmd/dev/release/publish.go
cmd/dev/schema/expected/render_version_test/.schema/version.schema.json
cmd/dev/schema/fixtures/render_version_test/.schema/version.schema.json
cmd/dev/schema/main.go
cmd/dev/schema/render_version.go
cmd/dev/schema/render_version_test.go
cmd/pkg/repos.go
go.mod
go.sum
63fd21e chore: update deprecated goreleaser config and add goimports linter (#14)
.github/workflows/checks-go.yml
.golangci.yml
.goreleaser.yml
cmd/dev/ci/github/env.go
cmd/dev/pop/migration/render.go
cmd/dev/release/compile.go
cmd/dev/release/publish.go
go_mod_version_pins.go
349145f autogen: pin v0.0.20 release commit
c6e59ad fix: auto-detect branch if on pr (#12)
cmd/dev/ci/github/env.go
06f2b8c autogen: pin v0.0.19 release commit
3cefe84 fix: revert goreleaser changes (#10)
.goreleaser.yml
ba956fb feat: add github helpers (#9)
cmd/dev/ci/github/env.go
cmd/dev/ci/github/main.go
cmd/dev/ci/main.go
ca1c3a3 ci: move to github actions (#3)
.github/workflows/checks-go.yml
.github/workflows/test.yml
.goreleaser.yml
cmd/dev/ci/orbs/bump.go
cmd/dev/pop/migration/sync.go
11c21d6 feat: add 60m timeout to goreleaser task
cmd/dev/release/compile.go
79f8f5d feat: add xgoreleaser compile helpers
cmd/dev/release/compile.go
b62ab90 autogen: pin v0.0.18 release commit
1ebd00c chore: update install.sh
install.sh
a4fd7e7 fix: use dedicated name template for libmusl
.goreleaser.yml
52ad42c autogen: pin v0.0.17 release commit
33b4a88 fix: resolve goreleaser archive issues
.goreleaser.yml
install.sh
ebfd612 fix: update install script
install.sh
0397148 autogen: pin v0.0.16 release commit
ab54199 fix: name ory binary ory
.goreleaser.yml
459b2f1 autogen: pin v0.0.15 release commit
9e8742c fix: update install script
install.sh
0fa040a fix: use libc build pipeline
.goreleaser.yml
4851f7a autogen: pin v0.0.14 release commit
037a8b5 fix: resolve goreleaser build issues
.goreleaser.yml
4895ade fix: escape go build flags
.goreleaser.yml
1c5af28 autogen: pin v0.0.13 release commit
df3a678 refactor: rename --from-version to --include-changelog-since (#7)
cmd/dev/release/notify/draft.go
cmd/dev/release/notify/draft_test.go
cmd/dev/release/publish.go
f76ef03 feat: add flag to select dialects for migration rendering (#6)
cmd/dev/pop/migration/render.go
5f84f2c fix: add workaround for sqlite pragma (#5)
cmd/dev/pop/migration/fizzx/file_migrator.go
cmd/dev/pop/migration/render.go
38a525e fix: resolve broken overwrite logic (#4)
cmd/dev/pop/migration/fizzx/file_migrator.go
go.mod
go.sum
go_mod_version_pins.go
3661f58 feat: implement fizz renderer (#2)
.goreleaser.yml
cmd/dev/newsletter/send.go
cmd/dev/pop/migration/fizzx/file_migrator.go
cmd/dev/pop/migration/render.go
cmd/dev/pop/migration/sync.go
cmd/dev/release/notify/draft.go
cmd/dev/release/notify/draft_test.go
cmd/dev/release/publish.go
cmd/pkg/cli.go
go.mod
go.sum
3893677 fix: use autogen prefix for git tag pin
cmd/pkg/git.go
61446dd fix: add dir to git tag confirm message
cmd/pkg/git.go
f595d22 chore: pin v0.0.12 release commit
2665e62 chore: pin v0.0.11 release commit
1ebc76c feat: add helper for bumping CI orbs (#1)
cmd/dev/ci/main.go
cmd/dev/ci/orbs/bump.go
cmd/dev/ci/orbs/main.go
cmd/dev/main.go
cmd/pkg/git.go
0dfba54 chore: pin v0.0.10 release commit
3936682 fix: resolve issues with from-version
cmd/dev/release/publish.go
c96c715 chore: pin v0.0.9 release commit
c9e7fd0 ci: bump goreleaser orb
.circleci/config.yml
aee9d9c chore: pin v0.0.8 release commit
a863735 chore: update install script
install.sh
78d7ee6 fix: reference correct brew archive
.goreleaser.yml
cb4017a chore: pin v0.0.7 release commit
298e9e8 ci: bump goreleaser orb
.circleci/config.yml
c27859e chore: pin v0.0.6 release commit
f4cf65f fix: bump goreleaser version
.circleci/config.yml
c849c1f chore: pin v0.0.5 release commit
0968c78 chore: pin v0.0.4 release commit
fda090e fix: add newline after release success message
cmd/dev/release/publish.go
7a439c8 chore: pin v0.0.3 release commit
dc79d1f fix: resolve broken tests
cmd/dev/markdown/render_test.go
cmd/dev/release/notify/draft_test.go
cmd/dev/release/publish.go
cmd/dev/swagger/stub/expected.json
cmd/dev/swagger/stub/in.json
3c8cc05 chore: pin v0.0.2 release commit
61b7fcd fix: trim git current tag
cmd/pkg/git.go
021018f ci: build tags and add test runner
.circleci/config.yml
7ec10c9 chore: pin v0.0.1 release commit
e3deb92 fix: fall back to v0.0.0 if no previous tag exists
cmd/pkg/git.go
3cbd939 feat: add debug trace to fatals
cmd/pkg/cli.go
1552d26 fix: use NewCommand in CommandGetOutput
cmd/pkg/git.go
895d645 fix: whitelist cli in releasable projects
cmd/dev/release/publish.go
ebb629f docs: add install instructions
README.md
87192d9 feat: add bash install script
install.sh
a81e31f feat: implement ory CLI
.circleci/config.yml
.editorconfig
.gitignore
.goreleaser.yml
README.md
cmd/dev.go
cmd/dev/main.go
cmd/dev/markdown/main.go
cmd/dev/markdown/render.go
cmd/dev/markdown/render_test.go
cmd/dev/newsletter/draft.go
cmd/dev/newsletter/main.go
cmd/dev/newsletter/pkg.go
cmd/dev/newsletter/send.go
cmd/dev/pop/main.go
cmd/dev/pop/migration/create.go
cmd/dev/pop/migration/main.go
cmd/dev/pop/migration/sync.go
cmd/dev/release/main.go
cmd/dev/release/notify/draft.go
cmd/dev/release/notify/draft_test.go
cmd/dev/release/notify/main.go
cmd/dev/release/notify/send.go
cmd/dev/release/publish.go
cmd/dev/swagger/main.go
cmd/dev/swagger/sanitize.go
cmd/dev/swagger/sanitize_test.go
cmd/pkg/circle.go
cmd/pkg/cli.go
cmd/pkg/git.go
cmd/root.go
cmd/version.go
go.mod
go.sum
go_mod_version_pins.go
main.go
test/changelog.md
test/note.md
view/mail-body.html
28296b4 Initial commit
LICENSE
README.md`

var list = &cobra.Command{
	Use:   "list",
	Short: "List all defined components.",
	Long:  `Read dependency configs and displays dependency graph.`,
	Run: func(cmd *cobra.Command, args []string) {

		var graph ComponentGraph
		graph.readConfiguration(rootDirectory)

		switch listMode {
		case "affected":
			fmt.Println("Not implemented yet!")
		case "all":
			graph.listComponents()
		case "changed":
			graph.listComponents()
			sampleOutput2, _ := changes(rootDirectory, parent)

			outputArray := strings.Split(sampleOutput2, "\n")
			cleanseChangeList(&outputArray, includeFiles)
			removeDuplicateLines(&outputArray)
			sort.Strings(outputArray)
			sortedString := strings.Join(outputArray, "\n")
			fmt.Printf("[%v]\n", sortedString)

			if !includeFiles {
				detectChangedComponents(graph, outputArray)
			}

		case "involved":
			fmt.Println("Not implemented yet!")
		default:
			log.Fatalf("Unknown ListMode '%s'", listMode)
		}

		/*
			resolved, err := graph.resolveGraph()
			if err != nil {
				fmt.Printf("Failed to resolve dependency graph: %s\n", err)
			} else {
				fmt.Println("The dependency graph resolved successfully")
			}
		*/
	},
}

func detectChangedComponents(graph ComponentGraph, changeList []string) []*Component {
	var changedComponents []*Component
	componentPaths := graph.componentPaths
	for _, changedPath := range changeList {
		for path, component := range componentPaths {
			if strings.HasPrefix(changedPath, path) {
				changedComponents = append(changedComponents, component)
				fmt.Printf("'%s' is subpath of '%s'\n", changedPath, path)
				fmt.Printf("Adding changed component: %s\n", component.String())
			}
		}
	}
	return changedComponents
}

func init() {
	Main.AddCommand(list)
	list.Flags().StringVarP(&listMode, "mode", "m", "all", "Define which components you want to get listed (affected, all, changed, involved). Default is all.")
	list.Flags().StringVarP(&parent, "parent", "p", "origin/master", "Parent branch used to determine changes!")
	list.Flags().BoolVarP(&includeFiles, "includeFiles", "i", false, "Include files in changeset!")
}
