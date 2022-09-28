package markdown

import (
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/pkg"
)

func init() {
	Main.AddCommand(render)
}

var render = &cobra.Command{
	Use:   "render <file.md>",
	Args:  cobra.ExactArgs(1),
	Short: "Render a Markdown file",
	Run: func(cmd *cobra.Command, args []string) {
		changelogRaw, err := os.ReadFile(args[0])
		pkg.Check(err)

		fmt.Println(renderMarkdown(changelogRaw))
	},
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
