package comments

import (
	"fmt"
	"strings"
)

// a comment format known to this app
type Format struct {
	// converts the given text into a comment in this format
	render formatFunc
	// converts the given text into the beginning of a comment
	renderStart formatFunc
}

// comment format that starts with a doubleslash
var doubleSlashComments = Format{
	render:      prependDoubleSlash,
	renderStart: prependDoubleSlash,
}

// comment format that starts with pound symbols
var poundComments = Format{
	render:      prependPound,
	renderStart: prependPound,
}

// HTML comment format
var htmlComments = Format{
	render:      wrapInHtmlComment,
	renderStart: prependHtmlComment,
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

// signature for functions that create comments for different programming languages
type formatFunc func(text string) string

// provides a YML-style comment containing the given text
func prependPound(text string) string {
	return makeComment(text, "# %s")
}

// provides a Go-style comment containing the given text
func prependDoubleSlash(text string) string {
	return makeComment(text, "// %s")
}

func wrapInHtmlComment(text string) string {
	return makeComment(text, "<!-- %s -->")
}

func prependHtmlComment(text string) string {
	return makeComment(text, "<!-- %s")
}

// creates a comment in the given comment style containing the given text
func makeComment(text, style string) string {
	result := []string{}
	for _, line := range strings.Split(text, "\n") {
		if line == "" {
			result = append(result, line)
		} else {
			result = append(result, fmt.Sprintf(style, line))
		}
	}
	return strings.Join(result, "\n")
}

// removes the comment block in the given format containing the given token from the given text
func remove(text string, format Format, token string) string {
	commentWithToken := format.renderStart(token)
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
		if !inComment {
			result = append(result, line)
		}
	}
	return strings.Join(result, "\n")
}
