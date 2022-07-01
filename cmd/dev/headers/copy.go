// Tool for copying files inside https://github.com/ory/meta/blob/master/scripts/sync.sh
// and adding a link to the original as a header comment.

package headers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ory/cli/cmd/dev/headers/comments"
	"github.com/spf13/cobra"
)

// template for the header
const LINK_TEMPLATE = "AUTO-GENERATED, DO NOT EDIT! Please edit the original at %s."

// the token that identifies comments containing the license
const LINK_TOKEN = "Copyright ©"

// the root path for links to the original
const ROOT_PATH = "https://github.com/ory/meta/blob/master/"

// copies the source file (relative path) to the given absolute path
func CopyFile(src, dst string) error {
	contentBytes, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("cannot read file %q: %w", src, err)
	}
	dstStat, err := os.Lstat(dst)
	destPath := dst
	if err == nil && dstStat.IsDir() {
		destPath = filepath.Join(dst, filepath.Base(src))
	}
	headerText := fmt.Sprintf(LINK_TEMPLATE, ROOT_PATH+src)
	comments.WriteFileWithHeader(destPath, headerText, contentBytes)
	return nil
}

var copy = &cobra.Command{
	Use:   "cp",
	Short: "Behaves like cp but adds a header pointing to the original to copied files.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return CopyFile(args[1], args[2])
	},
}

func init() {
	Main.AddCommand(copy)
	copy.Flags().BoolVarP(&recursive, "recursive", "R", false, "Whether to copy files in subdirectories")
}

var recursive bool
