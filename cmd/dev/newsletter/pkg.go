// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package newsletter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"

	"github.com/ory/cli/cmd/pkg"
	"github.com/ory/gochimp3"
	"github.com/ory/x/httpx"
)

var defaultRenderer = html.NewRenderer(html.RendererOptions{
	Flags: html.CommonFlags | html.HrefTargetBlank,
})

// return (ast.GoToNext, true) to tell html renderer to skip rendering this node
// (because you've rendered it)
func renderNodeHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	var b bytes.Buffer
	switch n := node.(type) {
	case *ast.Image:
		var b bytes.Buffer
		defaultRenderer.Image(&b, n, entering)
		_, _ = w.Write([]byte(strings.ReplaceAll(b.String(), "<img", "<img width=\"600\" border=\"0\" style=\"display: block; max-width: 100%; min-width: 100px; width: 100%;\"")))
	case *ast.CodeBlock:
		defaultRenderer.CodeBlock(&b, n)
		_, _ = w.Write([]byte(strings.ReplaceAll(b.String(), "<pre", "<pre style=\"word-break: break-all; white-space: pre-wrap\"")))
	default:
		return ast.GoToNext, false
	}

	return ast.GoToNext, true
}

func renderMarkdown(source []byte) template.HTML {
	var markdownRenderer = html.NewRenderer(html.RendererOptions{
		Flags:          html.CommonFlags | html.HrefTargetBlank,
		RenderNodeHook: renderNodeHook,
	})
	var markdownParser = parser.NewWithExtensions(
		parser.NoIntraEmphasis | parser.Tables | parser.FencedCode | parser.NoEmptyLineBeforeBlock |
			parser.Autolink | parser.Strikethrough | parser.SpaceHeadings | parser.DefinitionLists)

	rendered := string(markdown.ToHTML(source, markdownParser, markdownRenderer))
	//rendered = strings.ReplaceAll(rendered, "<p>", "")
	//rendered = strings.ReplaceAll(rendered, "</p>", "<br>")
	return template.HTML(rendered)
}

func newMailchimpRequest(apiKey, path string, payload interface{}) {
	u := url.URL{}
	u.Scheme = "https"
	u.Host = fmt.Sprintf(gochimp3.URIFormat, gochimp3.DatacenterRegex.FindString(apiKey))
	u.Path = filepath.Join(gochimp3.Version, path)
	req, err := retryablehttp.NewRequest("GET", u.String(), nil)
	pkg.Check(err)
	req.SetBasicAuth("gochimp3", apiKey)

	client := httpx.NewResilientClient()
	res, err := client.Do(req)
	pkg.Check(err)
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	pkg.Check(err)
	if res.StatusCode != http.StatusOK {
		pkg.Check(errors.Errorf("received unexpected status code: %d", res.StatusCode), "%s", body)
	}

	pkg.Check(json.NewDecoder(bytes.NewReader(body)).Decode(payload))
}

func campaignID() string {
	var repoName string
	ghRepo := strings.Split(os.Getenv("GITHUB_REPOSITORY"), "/")
	if len(ghRepo) == 2 {
		repoName = ghRepo[1]
	} else {
		repoName = pkg.MustGetEnv("CIRCLE_PROJECT_REPONAME")
	}

	return fmt.Sprintf("%s-%s-%s",
		repoName,
		substr(pkg.GitHubSHA(), 0, 6),
		pkg.GitHubTag())
}

func substr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}
