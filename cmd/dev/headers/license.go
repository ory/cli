// Tool for adding a license header to all supported files.

package headers

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"time"

	goGitIgnore "github.com/sabhiram/go-gitignore"
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/dev/headers/comments"
)

// LICENSE defines the full license text.
const LICENSE_TEMPLATE = "Copyright © %d Ory Corp"

// LICENSE_TOKEN defines the token that identifies comments containing the license.
const LICENSE_TOKEN = "Copyright ©"

// file types that we don't want to add license headers to
var noLicenseHeadersFor = []comments.FileType{"md", "yml", "yaml"}

// addLicenses adds or updates the Ory license header in all files within the given directory.
func AddLicenses(dir string, year int) error {
	licenseText := fmt.Sprintf(LICENSE_TEMPLATE, year)
	gitIgnore, _ := goGitIgnore.CompileIgnoreFile(filepath.Join(dir, ".gitignore"))
	return filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("cannot read directory %q: %w", path, err)
		}
		if info.IsDir() {
			return nil
		}
		if gitIgnore != nil && gitIgnore.MatchesPath(info.Name()) {
			return nil
		}
		if !comments.SupportsFile(path) {
			return nil
		}
		if !shouldAddLicense(path) {
			return nil
		}
		contentNoHeader, err := comments.FileContentWithoutHeader(path, LICENSE_TOKEN)
		if err != nil {
			return err
		}
		return comments.WriteFileWithHeader(path, licenseText, contentNoHeader)
	})
}

// indicates whether this tool is configured to add a license header to the file with the given path
func shouldAddLicense(path string) bool {
	return !comments.ContainsFileType(noLicenseHeadersFor, comments.GetFileType(path))
}

var copyright = &cobra.Command{
	Use:   "license",
	Short: "Adds the license header to all known files in the current directory",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		year, _, _ := time.Now().Date()
		return AddLicenses(args[0], year)
	},
}

func init() {
	Main.AddCommand(copyright)
}
