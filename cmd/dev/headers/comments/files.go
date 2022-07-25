package comments

import (
	"fmt"
	"os"
)

// FileContentWithoutHeader provides the content of the file with the given path,
// without the comment block identified by the given token.
func FileContentWithoutHeader(path, token string) (string, error) {
	buffer, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("cannot open file %q: %w", path, err)
	}
	text := string(buffer)
	format, knowsFormat := commentFormats[GetFileType(path)]
	if !knowsFormat {
		return text, nil
	}
	return remove(text, format, token), nil
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
	headerComment := format.render(header)
	content := fmt.Sprintf("%s\n\n%s", headerComment, body)
	count, err := file.WriteString(content)
	if err != nil {
		return fmt.Errorf("cannot write into file %q: %w", path, err)
	}
	if count != len(content) {
		return fmt.Errorf("did not write the full %d bytes of header into %q: %w", len(headerComment), path, err)
	}
	return nil
}
