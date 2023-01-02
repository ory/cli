// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package comments

import "strings"

// a comment format known to this app
type Format struct {
	// converts the given text into a comment in this format
	startToken string
	// converts the given beginning of a text line into the beginning of a comment line
	endToken string
}

// removes the comment block in the given format containing the given token from the given text
func (f Format) remove(text string, token string) string {
	commentWithToken := f.renderLineStart(token)
	inComment := false
	result := []string{}
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, commentWithToken) {
			inComment = true
		}
		if inComment && line == "" {
			// the type of comment blocks we remove here is separated by an empty line
			// --> empty line marks the end of our comment block
			inComment = false
			continue
		}
		if inComment && !strings.HasPrefix(line, f.startToken) {
			inComment = false
		}
		if !inComment {
			result = append(result, line)
		}
	}
	return strings.Join(result, "\n")
}

func (f Format) SplitHeaderFromContent(text string, token string) (header, content string) {
	commentWithToken := f.renderLineStart(token)
	inComment := false
	content_lines := []string{}
	header_lines := []string{}
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, commentWithToken) {
			inComment = true
		}
		if inComment && line == "" {
			// the type of comment blocks we remove here is separated by an empty line
			// --> empty line marks the end of our comment block
			inComment = false
			continue
		}
		if inComment && !strings.HasPrefix(line, f.startToken) {
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
	"md":   htmlComments,
	"php":  doubleSlashComments,
	"py":   poundComments,
	"rb":   poundComments,
	"rs":   doubleSlashComments,
	"ts":   doubleSlashComments,
	"vue":  htmlComments,
	"yml":  poundComments,
}

func GetFormat(path string) (Format, bool) {
	fmt, ok := commentFormats[GetFileType(path)]
	return fmt, ok
}
