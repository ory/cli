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

// FormatFunc defines the signature of functions to create comments for different programming languages.
type FormatFunc func(text string) string

// formatFuncs lists all formatFuncs known to this tool.
var formatFuncs = map[string]FormatFunc{
	"yml":  ymlComment,
	"yaml": ymlComment,
}

// addLicenses adds or updates the Ory license header
// in all files of the current directory or its subdirectories.
func addLicenses(dir string, year int) error {
	licenseText := fmt.Sprintf(LICENSE_TEMPLATE, year)
	filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			// filepath.Walk traverses subdirectories
			return nil
		}
		filetype := filepath.Ext(path)
		commentFunc, ok := formatFuncs[filetype]
		if !ok {
			// not a file that we can add comments to --> nothing to do here
			return nil
		}
		file, err := os.OpenFile(path, os.O_RDWR, 0755)
		buffer := make([]byte, info.Size())
		count, err := file.Read(buffer)
		if err != nil {
			return fmt.Errorf("cannot read file %s: %w", path, err)
		}
		if count != int(info.Size()) {
			return fmt.Errorf("did not read the entire file")
		}
		fileContent := string(buffer)
		fileContentWithoutHeader := removeHeader(fileContent, commentFunc, LICENSE_TOKEN)
		header := commentFunc(licenseText)
		fileContentWithNewHeader := fmt.Sprintf("%s\n%s", header, fileContentWithoutHeader)
		file.Truncate(0)
		file.WriteString(fileContentWithNewHeader)
		return file.Close()
	})
	return nil
}

func removeHeader(text string, commentFunc FormatFunc, token string) (result string) {
	commentWithToken := commentFunc(token)
	inComment := false
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, commentWithToken) {
			inComment = true
		}
		if line == "" {
			inComment = false
		}
		if !inComment {
			result += line
			result += "\n"
		}
	}
	return result
}

// addHeaderToFile adds the given header to the given file
func addHeaderToFile(header string, file string) error {
	return nil
}

// licenseText provides the complete text for the license.
func licenseText(text string, year int) {
	fmt.Sprintf(LICENSE_TEMPLATE, year)
}

// comment provides a valid multi-line comment for the given filetype containing the given text.
func comment(text, filetype string, year int) (string, error) {
	return "", nil
}

// commentForYML provides a YML comment containing the given text
func ymlComment(text string) (result string) {
	for _, line := range strings.Split(text, "\n") {
		result += fmt.Sprintf("# %s\n", line)
	}
	return result
}

var copyright = &cobra.Command{
	Use:   "copyright",
	Short: "Adds the copyright header to all known files in the current directory",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		year, _, _ := time.Now().Date()
		return addLicenses(args[0], year)
	},
}

func init() {
	Main.AddCommand(copyright)
}
