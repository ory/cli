// Copyright © 2022 Ory Corp

package newsletter

import (
	"bytes"
	"html/template"
	"os"
	"strings"
	"testing"

	"github.com/ory/x/snapshotx"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestRenderMarkdown(t *testing.T) {
	assert.EqualValues(t, `<p><img width="600" border="0" style="display: block; max-width: 100%; min-width: 100px; width: 100%;" src=".foo/bar" alt="Some image" /></p>

<pre style="word-break: break-all; white-space: pre-wrap"><code>foo
</code></pre>`,
		strings.TrimSpace(string(renderMarkdown([]byte(`
![Some image](.foo/bar)

`+"```\nfoo\n```"+`

`)))))
}

func TestRenderMarkdownLong(t *testing.T) {
	cl, err := os.ReadFile("stub/changelog.md.expected")
	require.NoError(t, err)

	tmplRaw, err := os.ReadFile("../../../view/mail-body.html")
	require.NoError(t, err)
	tmpl, err := template.New("view").Parse(string(tmplRaw))
	require.NoError(t, err)
	var body bytes.Buffer
	require.NoError(t, tmpl.Execute(&body, struct {
		Version     string
		GitTag      string
		ProjectName string
		RepoName    string
		Changelog   template.HTML
		Message     template.HTML
		BrandColor  string
	}{
		Version:     "v0.1.0",
		GitTag:      "v0.1.0",
		ProjectName: "Ory Kratos",
		RepoName:    "ory/kratos",
		Changelog:   renderMarkdown(cl),
		Message:     "iuaw4hri",
		BrandColor:  "#5528FF",
	}))

	require.NoError(t, os.WriteFile("stub/changelog.html.tmp", body.Bytes(), 0644))
	snapshotx.SnapshotT(t, body.String())
}
