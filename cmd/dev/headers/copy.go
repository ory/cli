// Tool for copying files inside https://github.com/ory/meta/blob/master/scripts/sync.sh
// and adding a link to the original as a header comment.

package headers

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"

	"github.com/ory/cli/cmd/dev/headers/comments"
	"github.com/spf13/cobra"
)

// template for the header
const LINK_TEMPLATE = "AUTO-GENERATED, DO NOT EDIT! Please edit the original at %s."

// the token that identifies comment blocks containing the license
const LINK_TOKEN = "Copyright ©"

// the root path for links to the original
const ROOT_PATH = "https://github.com/ory/meta/blob/master/"

// Copies the given source file (path must be relative to CWD) to the given absolute path.
// Behaves similar to the unix `cp` command.
func CopyFile(src, dst string) error {
	body, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("cannot read file %q: %w", src, err)
	}
	dstPath := copyFilesDstPath(src, dst)
	headerText := fmt.Sprintf(LINK_TEMPLATE, ROOT_PATH+src)
	comments.WriteFileWithHeader(dstPath, headerText, string(body))
	return nil
}

// Copies all files in the given `src` directory (path must be relative to CWD) to the given absolute path.
// Behaves similar to the unix `cp` command.
func CopyFiles(src, dst string) error {
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("cannot read directory %q: %w", path, err)
		}
		dstPath := dst + path[len(src):]
		if info.IsDir() {
			err := os.Mkdir(dstPath, 0744)
			if err == nil {
				return nil
			}
			// ignore folder already exists error
			pathErr := err.(*os.PathError)
			if pathErr.Err == syscall.EEXIST {
				return nil
			}
			return err
		}
		return CopyFile(path, dstPath)
	})
}

// Provides the full destination path for the cp operation of the given src file to the given dst destination.
// The dst value can be a path to a file or directory.
func copyFilesDstPath(src, dst string) string {
	dstStat, err := os.Lstat(dst)
	if err == nil && dstStat.IsDir() {
		return filepath.Join(dst, filepath.Base(src))
	} else {
		return dst
	}
}

var copy = &cobra.Command{
	Use:   "cp",
	Short: "Behaves like cp but adds a header pointing to the original to copied files.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if recursive {
			return CopyFiles(args[1], args[2])
		} else {
			return CopyFile(args[1], args[2])
		}
	},
}

func init() {
	Main.AddCommand(copy)
	copy.Flags().BoolVarP(&recursive, "recursive", "r", false, "Whether to copy files in subdirectories")
}

// contains the value of the "-r" CLI flag
var recursive bool
