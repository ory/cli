// Tool for adding a license header to all supported files.

package headers

import (
	"fmt"
	"io/fs"
	"os"
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
		filetype := comments.GetFileType(path)
		if comments.ContainsFileType(noLicenseHeadersFor, filetype) {
			return nil
		}
		commentFunc, ok := comments.FormatFuncs[filetype]
		if !ok {
			// not a file that we can add comments to --> nothing to do here
			return nil
		}
		file, err := os.OpenFile(path, os.O_RDWR, 0744)
		if err != nil {
			return fmt.Errorf("cannot open file %q: %w", path, err)
		}
		defer file.Close()
		buffer := make([]byte, info.Size())
		count, err := file.Read(buffer)
		if err != nil {
			return fmt.Errorf("cannot read file %q: %w", path, err)
		}
		if int64(count) != info.Size() {
			return fmt.Errorf("did not read the entire %d bytes of file %q but only %d", info.Size(), path, count)
		}
		pos, err := file.Seek(0, 0)
		if err != nil {
			return fmt.Errorf("cannot seek to beginning of file %q: %w", path, err)
		}
		if pos != 0 {
			return fmt.Errorf("didn't end up at the beginning of file %q after seeking but at %d: %w", path, pos, err)
		}
		err = file.Truncate(0)
		if err != nil {
			return fmt.Errorf("cannot truncate file %q: %w", path, err)
		}
		fileContent := string(buffer)
		fileContentNoHeader := comments.Remove(fileContent, commentFunc, LICENSE_TOKEN)

		newHeader := commentFunc(licenseText)
		fileContentNewHeader := fmt.Sprintf("%s\n\n%s", newHeader, fileContentNoHeader)
		count, err = file.WriteString(fileContentNewHeader)
		if err != nil {
			return fmt.Errorf("cannot write file %q: %w", path, err)
		}
		if count != len(fileContentNewHeader) {
			return fmt.Errorf("did not write the entire %d bytes into %q but only %d, file is corrupted now", len(fileContentNewHeader), path, count)
		}
		return nil
	})
	return nil
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
