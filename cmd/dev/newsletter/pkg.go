package newsletter

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/markbates/pkger/pkging"
	"github.com/ory/gochimp3"
	"github.com/ory/x/httpx"

	"github.com/ory/cli/cmd/pkg"
)

func readTemplate(file pkging.File, err error) *template.Template {
	pkg.Check(err)
	defer pkg.Check(file.Close())

	contents, err := ioutil.ReadAll(file)
	pkg.Check(err)

	t, err := template.New(file.Name()).Parse(string(contents))
	pkg.Check(err)
	return t
}

func renderMarkdown(source []byte) template.HTML {
	var markdownRenderer = html.NewRenderer(html.RendererOptions{Flags: html.CommonFlags | html.HrefTargetBlank})
	var markdownParser = parser.NewWithExtensions(
		parser.NoIntraEmphasis | parser.Tables | parser.FencedCode | parser.NoEmptyLineBeforeBlock |
			parser.Autolink | parser.Strikethrough | parser.SpaceHeadings | parser.DefinitionLists)

	rendered := string(markdown.ToHTML(source, markdownParser, markdownRenderer))
	rendered = strings.ReplaceAll(rendered, "<p>", "")
	rendered = strings.ReplaceAll(rendered, "</p>", "<br>")
	return template.HTML(rendered)
}

func newMailchimpRequest(apiKey, path string, payload interface{}) {
	u := url.URL{}
	u.Scheme = "https"
	u.Host = fmt.Sprintf(gochimp3.URIFormat, gochimp3.DatacenterRegex.FindString(apiKey))
	u.Path = filepath.Join(gochimp3.Version, path)
	req, err := http.NewRequest("GET", u.String(), nil)
	pkg.Check(err)
	req.SetBasicAuth("gochimp3", apiKey)
	client := httpx.NewResilientClientLatencyToleranceMedium(nil)
	res, err := client.Do(req)
	pkg.Check(err)
	defer res.Body.Close()
	pkg.Check(json.NewDecoder(res.Body).Decode(payload))
}

func campaignID() string {
	return fmt.Sprintf("%s-%s-%s",
		pkg.MustGetEnv("CIRCLE_PROJECT_REPONAME"),
		pkg.CircleSHA1(),
		pkg.CircleTag())
}
