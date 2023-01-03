// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package comments

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// FileContentWithoutHeader provides the content of the file with the given path,
// without the comment block identified by the given token.
func FileContentWithoutHeader(path string, headerRegexp *regexp.Regexp) (string, error) {
	text, err := FileContent(path)
	if err != nil {
		return "", err
	}
	format, knowsFormat := commentFormats[GetFileType(path)]
	if !knowsFormat {
		return text, nil
	}
	_, content := format.SplitHeaderFromContent(text, headerRegexp)
	return content, nil
}

func FileContent(path string) (string, error) {
	buffer, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("cannot open file %q: %w", path, err)
	}
	text := string(buffer)
	return text, nil
}

func StripPrefixes(fileContent string, prefixes []string) (string, string) {
	prefix := ""
	for _, p := range prefixes {
		if len(p) > len(prefix) && strings.HasPrefix(fileContent, p) {
			prefix = p
		}
	}

	return prefix, strings.TrimPrefix(fileContent, prefix)
}

// WriteFileWithHeader creates a file at the given path containing the given file content (header + body).
// The header argument should contain only text. This method will transform it into the correct comment format.
func WriteFileWithHeader(path, header, body string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot write file %q: %w", path, err)
	}
	defer file.Close()
	format, knowsFormat := commentFormats[GetFileType(path)]
	if !knowsFormat {
		return os.WriteFile(path, []byte(body), 0744)
	}
	headerComment := format.renderBlock(header)
	bom, orig := StripPrefixes(body, []string{
		// see: https://en.wikipedia.org/wiki/Byte_order_mark#Byte_order_marks_by_encoding
		"\xef\xbb\xbf",   // UTF-8
		"\ufffe",         // UTF-16 (LE)
		"\ufeff",         // UTF-16 (BE)
		"\ufffe\x00\x00", // UTF-32 (LE)
		"\x00\x00\ufeff", // UTF-32 (BE)
	})
	content := fmt.Sprintf("%s%s\n\n%s", bom, headerComment, orig)
	count, err := file.WriteString(content)
	if err != nil {
		return fmt.Errorf("cannot write into file %q: %w", path, err)
	}
	if count != len(content) {
		return fmt.Errorf("did not write the full %d bytes of header into %q: %w", len(headerComment), path, err)
	}
	return nil
}
