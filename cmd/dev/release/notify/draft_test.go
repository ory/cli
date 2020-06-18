package notify

import (
	"fmt"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/cli/cmd/pkg"
)

func TestHowManyVersionsSince(t *testing.T) {
	const list = `v0.0.1-alpha.1
v0.0.1-alpha.2
v0.0.1-alpha.3
v0.0.1-alpha.5
v0.0.1-alpha.6
v0.0.1-alpha.7
v0.0.1-alpha.8
v0.0.2-alpha.1
v0.0.1-alpha.10+oryOS.15
v0.0.1-alpha.11
v0.0.1-alpha.9
v0.0.3-alpha.1
v0.0.3-alpha.2
v0.0.3-alpha.3
v0.0.3-alpha.4
v0.0.3-alpha.5
v0.0.3-alpha.7
v0.0.3-alpha.8+oryOS.15
v0.0.3-alpha.9
v0.0.3-alpha.10
v0.0.3-alpha.11
v0.0.3-alpha.12
v0.0.3-alpha.13
v0.0.3-alpha.14
v0.0.3-alpha.15
v0.1.0-alpha.1
v0.1.0-alpha.2
v0.1.0-alpha.3
v0.1.0-alpha.4
v0.1.0-alpha.5
v0.1.0-alpha.6
v0.1.1-alpha.1
v0.2.0-alpha.2
v0.2.1-alpha.1
v0.3.0-alpha.1
`

	assert.Equal(t, 2, howManyVersionsSince(semver.MustParse("0.3.0-alpha.1"), list))
	assert.Equal(t, 3, howManyVersionsSince(semver.MustParse("0.2.1-alpha.1"), list))
	assert.Equal(t, 5, howManyVersionsSince(semver.MustParse("0.1.1-alpha.1"), list))
	assert.Equal(t, 2, howManyVersionsSince(semver.MustParse("9.9.9"), list))
}

func TestGetPreviousVersionFromGitCommitMessage(t *testing.T) {
	for k, tc := range []struct {
		m string
		v string
	}{
		{
			m: "not the message you're looking for",
		},
		{
			m: pkg.GitCommitMessagePreviousVersion + " v1.2.3-beta.1",
			v: "v1.2.3-beta.1",
		},
		{
			m: pkg.GitCommitMessagePreviousVersion + " v1.2.3-beta1",
			v: "v1.2.3-beta1",
		},
		{
			m: pkg.GitCommitMessagePreviousVersion + " v1.2.3",
			v: "v1.2.3",
		},
		{
			m: `chore: some header

Some message

` + pkg.GitCommitMessagePreviousVersion + ` v1.2.3
`,
			v: "v1.2.3",
		},
		{
			m: `chore: some header

Some message

` + pkg.GitCommitMessagePreviousVersion + ` v1.2.3

Super multiline
`,
			v: "v1.2.3",
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			v, o := getPreviousVersionFromGitCommitMessage(tc.m)
			if tc.v == "" {
				require.False(t, o)
			} else {
				require.True(t, o)
				assert.EqualValues(t, tc.v, "v"+v.String())
			}
		})
	}
}
