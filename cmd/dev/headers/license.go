// Tool for adding a license header to all supported files.

package headers

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/ory/cli/cmd/dev/headers/comments"
	gitIgnore "github.com/sabhiram/go-gitignore"
	"github.com/spf13/cobra"
)

// LICENSE defines the full license text.
const LICENSE_TEMPLATE = "Copyright © %d Ory Corp Inc."

// LICENSE_TOKEN defines the token that identifies comments containing the license.
const LICENSE_TOKEN = "Copyright ©"

// file types that we don't want to add license headers to
var noLicenseHeadersFor = []comments.FileType{"md"}

// addLicenses adds or updates the Ory license header in all files within the given directory.
func AddLicenses(dir string, year int) error {
	licenseText := fmt.Sprintf(LICENSE_TEMPLATE, year)
	ignore, _ := gitIgnore.CompileIgnoreFile(filepath.Join(dir, ".gitignore"))
	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("cannot read directory %q: %w", path, err)
		}
		if info.IsDir() {
			// we'll traverse subdirectories through filepath.Walk automatically
			return nil
		}
		if ignore != nil && ignore.MatchesPath(info.Name()) {
			// file is git-ignored
			return nil
		}
		if !comments.Supports(path) {
			// we don't know how to write comments for this file
			return nil
		}
		if !shouldAddLicense(path) {
			// this tool is configured to not add licenses for this file type
			return nil
		}
		contentNoHeader, err := comments.FileContentWithoutHeader(path, LICENSE_TOKEN)
		if err != nil {
			return err
		}
		return comments.WriteFileWithHeader(path, licenseText, []byte(contentNoHeader))
	})
	return nil
}

// indicates whether we should add a license header to the file with the given path
func shouldAddLicense(path string) bool {
	return !comments.ContainsFileType(noLicenseHeadersFor, comments.GetFileType(path))
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
