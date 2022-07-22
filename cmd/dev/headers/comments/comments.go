// helper functions for creating and removing comments from source code files in a variety of programming languages
package comments

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// signature for functions that create comments for different programming languages
type formatFunc func(text string) string

// a file format that we know about, represented as its file extension
type FileType string

// all file formats that we can create comments for
var formatFuncs = map[FileType]formatFunc{
	"cs":   prependDoubleSlash,
	"dart": prependDoubleSlash,
	"go":   prependDoubleSlash,
	"java": prependDoubleSlash,
	"js":   prependDoubleSlash,
	"md":   wrapInHtmlComment,
	"php":  prependDoubleSlash,
	"py":   prependPound,
	"rb":   prependPound,
	"rs":   prependDoubleSlash,
	"ts":   prependDoubleSlash,
	"vue":  wrapInHtmlComment,
}

// indicates whether it is possible to add comments to the file with the given name
func Supports(filename string) bool {
	filetype := GetFileType(filename)
	_, ok := formatFuncs[filetype]
	return ok
}

// provides the extension of the given filename
func GetFileType(filename string) FileType {
	ext := filepath.Ext(filename)
	if len(ext) == 0 {
		return ""
	}
	return FileType(ext[1:])
}

// creates a comment in the given comment style containing the given text
func renderComment(text, style string) string {
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

// provides a YML-style comment containing the given text
func prependPound(text string) string {
	return renderComment(text, "# %s")
}

// provides a Go-style comment containing the given text
func prependDoubleSlash(text string) string {
	return renderComment(text, "// %s")
}

func wrapInHtmlComment(text string) string {
	return renderComment(text, "<!-- %s -->")
}

// removes the comment block in the given format containing the given token from the given text
func remove(text string, format formatFunc, token string) string {
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
func FileContentWithoutHeader(path, token string) (string, error) {
	buffer, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("cannot open file %q: %w", path, err)
	}
	fileType := GetFileType(path)
	formatter := formatFuncs[fileType]
	text := string(buffer)
	return remove(text, formatter, token), nil
}

func WriteFileWithHeader(path, header string, body string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot write file %q: %w", path, err)
	}
	defer file.Close()
	filetype := GetFileType(path)
	format, ok := formatFuncs[filetype]
	if !ok {
		return os.WriteFile(path, []byte(body), 0744)
	}
	headerComment := format(header)
	newContent := fmt.Sprintf("%s\n\n%s", headerComment, body)
	count, err := file.WriteString(newContent)
	if err != nil {
		return fmt.Errorf("cannot write into file %q: %w", path, err)
	}
	if count != len(newContent) {
		return fmt.Errorf("did not write the full %d bytes of header into %q: %w", len(headerComment), path, err)
	}
	return nil
}
