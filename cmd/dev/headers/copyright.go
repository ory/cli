// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

// Tool for adding copyright headers to files.

package headers

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	ignore "github.com/sabhiram/go-gitignore"
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/dev/headers/comments"
)

// HEADER_TEMPLATE defines the full header text.
const HEADER_TEMPLATE = "Copyright © %d Ory Corp\nSPDX-License-Identifier: Apache-2.0"

// HEADER_TOKEN defines a text snippet to recognize an existing copyright header in a file.
const HEADER_TOKEN = "Copyright ©"

// file types that we don't want to add copyright headers to
var noHeadersFor = []comments.FileType{"md", "yml", "yaml"}

// folders that are excluded by default
var defaultExcludedFolders = []string{"dist", "node_modules", "vendor"}

// AddHeaders adds or updates the Ory copyright header in all applicable files within the given directory.
func AddHeaders(dir string, year int, exclude []string) error {
	headerText := fmt.Sprintf(HEADER_TEMPLATE, year)
	gitIgnore, _ := ignore.CompileIgnoreFile(filepath.Join(dir, ".gitignore"))
	prettierIgnore, _ := ignore.CompileIgnoreFile(filepath.Join(dir, ".prettierignore"))
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
		if prettierIgnore != nil && prettierIgnore.MatchesPath(info.Name()) {
			return nil
		}
		if !comments.SupportsFile(path) {
			return nil
		}
		if !fileTypeNeedsCopyrightHeader(path) {
			return nil
		}
		if isInFolders(path, defaultExcludedFolders) {
			return nil
		}
		if isInFolders(path, exclude) {
			return nil
		}
		contentNoHeader, err := comments.FileContentWithoutHeader(path, HEADER_TOKEN)
		if err != nil {
			return err
		}
		return comments.WriteFileWithHeader(path, headerText, contentNoHeader)
	})
}

// isInFolders indicates whether the given path exists within the given list of folders
func isInFolders(path string, exclude []string) bool {
	for _, e := range exclude {
		if strings.HasPrefix(path, e) {
			return true
		}
	}
	return false
}

// indicates whether this tool should add a copyright header to the given file
func fileTypeNeedsCopyrightHeader(path string) bool {
	return !comments.ContainsFileType(noHeadersFor, comments.GetFileType(path))
}

var copyright = &cobra.Command{
	Use:   "copyright",
	Short: "Adds the copyright header to all files in the current directory",
	Long: `Adds the copyright header to all files that need one in the current directory.

Does not add the header to files listed in .gitignore and .prettierignore.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		year, _, _ := time.Now().Date()
		return AddHeaders(".", year, exclude)
	},
}

func init() {
	Main.AddCommand(copyright)
	copyright.Flags().StringSliceVarP(&exclude, "exclude", "e", []string{}, "folders to exclude, provide comma-separated values or multiple instances of this flag")
}

// contains the folders to exclude
var exclude []string
