package headers

import "github.com/spf13/cobra"

var Main = &cobra.Command{
	Use:   "headers",
	Short: "Adds language-specific headers to files",
}
