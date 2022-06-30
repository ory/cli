package headers

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// LICENSE defines the full license text.
const LICENSE_TEMPLATE = "Copyright © %d Ory Corp Inc."

// LICENSE_TOKEN defines the token that identifies comments containing the license.
const LICENSE_TOKEN = "Copyright ©"

// all file formats that get licenses
var formatFuncs = map[string]FormatFunc{
	"go":   PrependDoubleSlash,
	"java": PrependDoubleSlash,
	"js":   PrependDoubleSlash,
	"php":  PrependDoubleSlash,
	"py":   PrependPound,
	"rb":   PrependPound,
	"rs":   PrependDoubleSlash,
	"ts":   PrependDoubleSlash,
}

// addLicenses adds or updates the Ory license header in all files within the given directory.
func AddLicenses(dir string, year int) error {
	licenseText := fmt.Sprintf(LICENSE_TEMPLATE, year)
	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("cannot read directory %s: %w", path, err)
		}
		if info.IsDir() {
			// we'll traverse subdirectories through filepath.Walk automatically
			return nil
		}
		filetype := FileExt(path)
		commentFunc, ok := formatFuncs[filetype]
		if !ok {
			// not a file that we can add comments to --> nothing to do here
			return nil
		}
		file, err := os.OpenFile(path, os.O_RDWR, 0744)
		if err != nil {
			return fmt.Errorf("cannot open file %s: %w", path, err)
		}
		defer file.Close()
		buffer := make([]byte, info.Size())
		count, err := file.Read(buffer)
		if err != nil {
			return fmt.Errorf("cannot read file %s: %w", path, err)
		}
		if int64(count) != info.Size() {
			return fmt.Errorf("did not read the entire %d bytes of file %s but only %d", info.Size(), path, count)
		}
		fileContent := string(buffer)
		fileContentNoHeader := Remove(fileContent, commentFunc, LICENSE_TOKEN)
		newHeader := commentFunc(licenseText)
		fileContentNewHeader := fmt.Sprintf("%s\n\n%s", newHeader, fileContentNoHeader)
		pos, err := file.Seek(0, 0)
		if err != nil {
			return fmt.Errorf("cannot seek to beginning of file %s: %w", path, err)
		}
		if pos != 0 {
			return fmt.Errorf("didn't end up at the beginning of file %s after seeking but at %d: %w", path, pos, err)
		}
		err = file.Truncate(0)
		if err != nil {
			return fmt.Errorf("cannot truncate file %s: %w", path, err)
		}
		count, err = file.WriteString(fileContentNewHeader)
		if err != nil {
			return fmt.Errorf("cannot write file %s: %w", path, err)
		}
		if count != len(fileContentNewHeader) {
			return fmt.Errorf("did not write the entire %d bytes into %s but only %d, file is corrupted now", len(fileContentNewHeader), path, count)
		}
		return nil
	})
	return nil
}

// FormatFunc defines the signature of functions to create comments for different programming languages.
type FormatFunc func(text string) string

// FileExt provides the extension of the given file
func FileExt(filename string) string {
	ext := filepath.Ext(filename)
	if len(ext) == 0 {
		return ""
	}
	return ext[1:]
}

func Remove(text string, commentFunc FormatFunc, token string) string {
	commentWithToken := commentFunc(token)
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

var copyright = &cobra.Command{
	Use:   "copyright",
	Short: "Adds the copyright header to all known files in the current directory",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		year, _, _ := time.Now().Date()
		return AddLicenses(args[0], year)
	},
}

func init() {
	Main.AddCommand(copyright)
}
