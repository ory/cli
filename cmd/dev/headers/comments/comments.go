// helper functions for creating and removing comments from source code files in a variety of programming languages
package comments

import (
	"fmt"
	"path/filepath"
	"strings"
)

// signature of functions to create comments for different programming languages
type FormatFunc func(text string) string

// a file format that we know about, represented as its file extension
type FileType string

// all file formats that we can create comments for
var FormatFuncs = map[FileType]FormatFunc{
	"cs":   PrependDoubleSlash,
	"dart": PrependDoubleSlash,
	"go":   PrependDoubleSlash,
	"java": PrependDoubleSlash,
	"js":   PrependDoubleSlash,
	"md":   WrapInHtmlComment,
	"php":  PrependDoubleSlash,
	"py":   PrependPound,
	"rb":   PrependPound,
	"rs":   PrependDoubleSlash,
	"ts":   PrependDoubleSlash,
	"vue":  WrapInHtmlComment,
}

// provides the extension of the given filename
func GetFileType(filename string) FileType {
	ext := filepath.Ext(filename)
	if len(ext) == 0 {
		return ""
	}
	return FileType(ext[1:])
}

// PrependPound provides a YML comment containing the given text.
func PrependPound(text string) string {
	result := []string{}
	for _, line := range strings.Split(text, "\n") {
		if line == "" {
			result = append(result, line)
		} else {
			result = append(result, fmt.Sprintf("# %s", line))
		}
	}
	return strings.Join(result, "\n")
}

// PrependDoubleSlash provides a Go comment containing the given text.
func PrependDoubleSlash(text string) string {
	result := []string{}
	for _, line := range strings.Split(text, "\n") {
		if line == "" {
			result = append(result, line)
		} else {
			result = append(result, fmt.Sprintf("// %s", line))
		}
	}
	return strings.Join(result, "\n")
}

func WrapInHtmlComment(text string) string {
	result := []string{}
	for _, line := range strings.Split(text, "\n") {
		if line == "" {
			result = append(result, line)
		} else {
			result = append(result, fmt.Sprintf("<!-- %s -->", line))
		}
	}
	return strings.Join(result, "\n")
}

// removes the comment in the given format containing the given token from the given text
func Remove(text string, format FormatFunc, token string) string {
	commentWithToken := format(token)
	inComment := false
	result := []string{}
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, commentWithToken) {
			inComment = true
		}
		if line == "" && inComment {
			inComment = false
			continue
		}
		if !inComment {
			result = append(result, line)
		}
	}
	return strings.Join(result, "\n")
}

// the missing "contains" function in Go, indicates whether the given list of FileTypes contains the given FileType
func ContainsFileType(fileTypes []FileType, fileType FileType) bool {
	for _, ft := range fileTypes {
		if ft == fileType {
			return true
		}
	}
	return false
}

// provides the content of the file with the given path, without the header identified by the given token
// func FileContentWithoutHeader(path, token string) (string, FormatFunc, error) {}

// func WriteFileWithHeader(path, header string, format FormatFunc, body string) error {}
