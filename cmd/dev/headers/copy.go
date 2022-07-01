// Tool for copying files inside https://github.com/ory/meta/blob/master/scripts/sync.sh
// and adding a link to the original as a header comment.

package headers

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// template for the header
const LINK_TEMPLATE = "AUTO-GENERATED, DO NOT EDIT! Please edit the original at %s."

// the token that identifies comments containing the license
const LINK_TOKEN = "Copyright ©"

// the root path for links to the original
const ROOT_PATH = "https://github.com/ory/meta/blob/master/"

// copies the source file (relative path) to the given absolute path
func CopyFile(source, destPath string) error {
	contentBytes, err := os.ReadFile(source)
	if err != nil {
		return fmt.Errorf("cannot read file %q: %w", source, err)
	}
	filetype := FileExt(source)
	commentFunc, ok := formatFuncs[filetype]
	if !ok {
		// not a file that we can add comments to
		return os.WriteFile(destPath, contentBytes, 0744)
	}
	headerText := fmt.Sprintf(LINK_TEMPLATE, ROOT_PATH+source)
	headerComment := commentFunc(headerText)
	file, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("cannot write file %q: %w", destPath, err)
	}
	defer file.Close()
	count, err := file.WriteString(headerComment)
	if err != nil {
		return fmt.Errorf("cannot write into file %q: %w", destPath, err)
	}
	if count != len(headerComment) {
		return fmt.Errorf("did not write the full %d bytes of header into %q: %w", len(headerComment), destPath, err)
	}
	count, err = file.WriteString("\n\n")
	if err != nil {
		return fmt.Errorf("cannot write into file %q: %w", destPath, err)
	}
	if count != 2 {
		return fmt.Errorf("did not write the full %d bytes of header into %q: %w", len(headerComment), destPath, err)
	}
	count, err = file.Write(contentBytes)
	if err != nil {
		return fmt.Errorf("cannot write into file %q: %w", destPath, err)
	}
	if count != len(contentBytes) {
		return fmt.Errorf("did not write the full %d bytes of header into %q: %w", len(headerComment), destPath, err)
	}
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
	Main.AddCommand(copyright)
	copy.Flags().BoolVarP(&recursive, "recursive", "R", false, "Whether to copy files in subdirectories")
}

var recursive bool
