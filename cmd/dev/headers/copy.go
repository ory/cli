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
	var dstPath = determineDestPath(src, dst)
	headerText := fmt.Sprintf(LINK_TEMPLATE, ROOT_PATH+src)
	comments.WriteFileWithHeader(dstPath, headerText, string(contentBytes))
	return nil
}

func CopyFiles(src, dst string) error {
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("cannot read directory %q: %w", path, err)
		}
		dstPath := createDstPath(path, dst, src)
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

func createDstPath(path, dst, src string) string {
	return dst + path[len(src):]
}

// Determines the full destination path for the cp operation of the given src to the given dst.
// The dst value can be a full path to a file or a path to the directory to put the file in.
func determineDestPath(src, dst string) string {
	if isDir(dst) {
		return filepath.Join(dst, filepath.Base(src))
	} else {
		return dst
	}
}

// indicates whether the given file path points to a directory
func isDir(filepath string) bool {
	stat, err := os.Lstat(filepath)
	return err == nil && stat.IsDir()
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
