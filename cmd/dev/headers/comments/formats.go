// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package comments

import (
	"regexp"
	"strings"
)

// a comment format known to this app
type Format struct {
	// converts the given text into a comment in this format
	startToken string
	// converts the given beginning of a text line into the beginning of a comment line
	endToken string
}

func (f Format) SplitHeaderFromContent(text string, headerRegexp *regexp.Regexp) (header, content string) {
	inComment := false
	content_lines := []string{}
	header_lines := []string{}
	for _, line := range strings.Split(text, "\n") {
		if f.isComment(line) && headerRegexp.MatchString(line) {
			inComment = true
		}
		if inComment && line == "" {
			// the type of comment blocks we remove here is separated by an empty line
			// --> empty line marks the end of our comment block
			inComment = false
			continue
		}
		if inComment && !f.isComment(line) {
			inComment = false
		}
		if !inComment {
			content_lines = append(content_lines, line)
		} else {
			header_lines = append(header_lines, line)
		}
	}
	return strings.Join(header_lines, "\n"), strings.Join(content_lines, "\n")
}

func (f Format) isComment(line string) bool {
	return strings.HasPrefix(line, f.startToken)
}

// renders the given text block (consisting of many text lines) into a comment block
func (f Format) renderBlock(text string) string {
	result := []string{}
	for _, line := range strings.Split(text, "\n") {
		if line != "" {
			line = f.renderLine(line)
		}
		result = append(result, line)
	}
	return strings.Join(result, "\n")
}

// renders the given text line into a comment line of this format
func (f Format) renderLine(text string) string {
	return f.startToken + text + f.endToken
}

// renders the given text line part into the beginning of a comment line of this format
func (f Format) renderLineStart(text string) string {
	return f.startToken + text
}

// comment format that starts with a doubleslash
var doubleSlashComments = Format{
	startToken: "// ",
	endToken:   "",
}

// comment format that is surrounded by /* */
var slashStarComments = Format{
	startToken: "/* ",
	endToken:   " */",
}

// comment format that starts with pound symbols
var poundComments = Format{
	startToken: "# ",
	endToken:   "",
}

// HTML comment format
var htmlComments = Format{
	startToken: "<!-- ",
	endToken:   " -->",
}

// all file formats that we can create comments for, and how to do it
var commentFormats = map[FileType]Format{
	"cs":   doubleSlashComments,
	"dart": doubleSlashComments,
	"go":   doubleSlashComments,
	"java": doubleSlashComments,
	"js":   doubleSlashComments,
	"jsx":  doubleSlashComments,
	"md":   htmlComments,
	"php":  doubleSlashComments,
	"py":   poundComments,
	"rb":   poundComments,
	"rs":   doubleSlashComments,
	"ts":   doubleSlashComments,
	"tsx":  doubleSlashComments,
	"vue":  htmlComments,
	"yml":  poundComments,
	"html": htmlComments,
	"css":  slashStarComments,
}

func GetFormat(path string) (Format, bool) {
	fmt, ok := commentFormats[GetFileType(path)]
	return fmt, ok
}
