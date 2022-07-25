package comments

import (
	"fmt"
	"os"
)

// provides the content of the file with the given path
// without the comment block identified by the given token
func FileContentWithoutHeader(path, token string) (string, error) {
	buffer, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("cannot open file %q: %w", path, err)
	}
	text := string(buffer)
	fileType := GetFileType(path)
	format, knowsFormat := commentFormats[fileType]
	if !knowsFormat {
		return text, nil
	}
	return remove(text, format, token), nil
}

// WriteFileWithHeader creates a file at the given path containing the given file content (header comment + body).
// the header comment is supposed to contain only text. This method will transform it into a comment.
func WriteFileWithHeader(path, header, body string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot write file %q: %w", path, err)
	}
	defer file.Close()
	filetype := GetFileType(path)
	format, ok := commentFormats[filetype]
	if !ok {
		return os.WriteFile(path, []byte(body), 0744)
	}
	headerComment := format.render(header)
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
