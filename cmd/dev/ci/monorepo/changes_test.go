// Copyright © 2022 Ory Corp

package monorepo

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var sampleGitOutput = `303599b feat: improve create migrations command (#16)
cmd/dev/pop/migration/create.go
ca21cb5 autogen: pin v0.0.24 release commit
cmd/dev/schema/render_version.go
cmd/dev/schema/render_version_test.go
cmd/pkg/repos.go
README.md
go.mod
go.sum
63fd21e chore: update deprecated goreleaser config and add goimports linter (#14)
.github/workflows/checks-go.yml
go.sum
test/changelog.md`

var changelogWithoutCommitMessages = `cmd/dev/pop/migration/create.go
cmd/dev/schema/render_version.go
cmd/dev/schema/render_version_test.go
cmd/pkg/repos.go
README.md
go.mod
go.sum
.github/workflows/checks-go.yml
go.sum
test/changelog.md`

var deduplicatedChangelogOutput = `303599b feat: improve create migrations command (#16)
cmd/dev/pop/migration/create.go
ca21cb5 autogen: pin v0.0.24 release commit
cmd/dev/schema/render_version.go
cmd/dev/schema/render_version_test.go
cmd/pkg/repos.go
README.md
go.mod
go.sum
63fd21e chore: update deprecated goreleaser config and add goimports linter (#14)
.github/workflows/checks-go.yml
test/changelog.md`

var cleansedDirectoriesChangelogOutput = `cmd/dev/pop/migration
cmd/dev/schema
cmd/pkg
.
.github/workflows
test`

var cleansedFilesChangelogOutput = `cmd/dev/pop/migration/create.go
cmd/dev/schema/render_version.go
cmd/dev/schema/render_version_test.go
cmd/pkg/repos.go
README.md
go.mod
go.sum
.github/workflows/checks-go.yml
test/changelog.md`

var sortedCleansedFilesChangelogOutput = `.github/workflows/checks-go.yml
cmd/dev/pop/migration/create.go
cmd/dev/schema/render_version.go
cmd/dev/schema/render_version_test.go
cmd/pkg/repos.go
go.mod
go.sum
README.md
test/changelog.md`

func TestDeduplication(t *testing.T) {
	changelogArray := strings.Split(sampleGitOutput, "\n")
	deduplicateChangelog(&changelogArray)
	dedup := strings.Join(changelogArray, "\n")
	assert.Equalf(t, deduplicatedChangelogOutput, dedup, "Deduplicated!", nil)
}

func TestRemovingCommitMessages(t *testing.T) {
	cleansedChangelog := removeCommitMessages(sampleGitOutput)
	assert.Equalf(t, changelogWithoutCommitMessages, cleansedChangelog, "Cleansed!", nil)
}

func TestCleanseChangelogFiles(t *testing.T) {
	changelogArray := strings.Split(changelogWithoutCommitMessages, "\n")
	cleanseRepositoryChanges(&changelogArray, true, true)
	cleansedFilesChangelog := strings.Join(changelogArray, "\n")
	assert.Equalf(t, cleansedFilesChangelogOutput, cleansedFilesChangelog, "Cleansed Files Changelog!", nil)
}

func TestCleanseChangelogDirectories(t *testing.T) {
	changelogArray := strings.Split(changelogWithoutCommitMessages, "\n")
	cleanseRepositoryChanges(&changelogArray, false, true)
	cleansedDirectoryChangelog := strings.Join(changelogArray, "\n")
	assert.Equalf(t, cleansedDirectoriesChangelogOutput, cleansedDirectoryChangelog, "Cleansed Directory Changelog!", nil)
}

func TestCaseInsensitiveSorting(t *testing.T) {
	changelogArray := strings.Split(cleansedFilesChangelogOutput, "\n")
	caseInsensitiveSort(changelogArray)
	sortedChangelog := strings.Join(changelogArray, "\n")
	assert.Equalf(t, sortedCleansedFilesChangelogOutput, sortedChangelog, "Sorted Cleansed Files Changelog!", nil)
}
