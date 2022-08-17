// Tool for copying files inside https://github.com/ory/meta/blob/master/scripts/sync.sh
// and adding a link to the original as a header comment.

package headers

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/dev/headers/comments"
)

// template for the header
const COPY_HEADER_TEMPLATE = "AUTO-GENERATED, DO NOT EDIT!\nPlease edit the original at %s"

// the root path for links to the original
// NOTE: might have to convert to a CLI switch
const ROOT_PATH = "https://github.com/ory/meta/blob/master/"

// Header-aware equivalent of the Unix `cp` command.
// Copies the given source file (path must be relative to CWD) to the given absolute path
// and prepends the COPY_HEADER_TEMPLATE to the content.
func CopyFile(src, dst string) error {
	if strings.HasSuffix(dst, "/") {
		return fmt.Errorf("cannot create file %q", dst)
	}
	body, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("cannot read file %q: %w", src, err)
	}
	dstPath := dst
	dstStat, err := os.Lstat(dst)
	if err == nil && dstStat.IsDir() {
		dstPath = filepath.Join(dst, filepath.Base(src))
	}
	headerText := fmt.Sprintf(COPY_HEADER_TEMPLATE, ROOT_PATH+src)
	comments.WriteFileWithHeader(dstPath, headerText, string(body))
	return nil
}

// Header-aware equivalent of the Unix `cp -r` command.
// Copies all files in the given `src` directory (path must be relative to CWD) to the given absolute path
// and prepends the COPY_HEADER_TEMPLATE to the content.
func CopyFiles(src, dst string) error {
	srcStat, err := os.Lstat(src)
	if err != nil {
		return err
	}
	if !srcStat.IsDir() {
		return CopyFile(src, dst)
	}
	hasDst, err := folderExists(dst)
	if err != nil {
		return fmt.Errorf("cannot determine if folder %q exists: %w", dst, err)
	}
	extraPath := ""
	if hasDst {
		srcLast := filepath.Base(src)
		os.MkdirAll(filepath.Join(dst, srcLast), 0744)
		extraPath = srcLast
	}
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("cannot read directory %q: %w", path, err)
		}
		dstPath := filepath.Join(dst, extraPath, path[len(src):])
		if info.IsDir() {
			err := os.MkdirAll(dstPath, 0744)
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

// indicates whether the folder with the given path exists
func folderExists(path string) (bool, error) {
	dstStat, err := os.Lstat(path)
	if err != nil {
		pathErr := err.(*os.PathError)
		if pathErr.Err == syscall.ENOENT {
			return false, nil
		}
		return false, err
	}
	return dstStat.IsDir(), nil
}

var copy = &cobra.Command{
	Use:   "cp",
	Short: "Behaves like cp but adds a header pointing to the original to copied files.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if recursive {
			return CopyFiles(args[0], args[1])
		} else {
			return CopyFile(args[0], args[1])
		}
	},
}

func init() {
	Main.AddCommand(copy)
	copy.Flags().BoolVarP(&recursive, "recursive", "r", false, "Whether to copy files in subdirectories")
}

// contains the value of the "-r" CLI flag
var recursive bool
