package pat

import (
	"errors"
	"fmt"
	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
	"github.com/spf13/cobra"
)

func NewCreatePATCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pat <project-id>",
		Args:  cobra.ExactArgs(1),
		Short: fmt.Sprintf("Create a new Ory Cloud Project"),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := client.NewCommandHelper(cmd)
			if err != nil {
				return err
			}

			name := flagx.MustGetString(cmd, "name")
			if len(name) == 0 && flagx.MustGetBool(cmd, cmdx.FlagQuiet) {
				return errors.New("you must specify the --name flag when using --quiet")
			}

			stdin := h.Stdin
			for name == "" {
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Enter a name for your personal access token: ")
				name, err = stdin.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read from stdin: %w", err)
				}
			}

			pat, err := h.CreatePAT(args[0], name)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintln(h.VerboseErrWriter, "PAT created successfully!")
			cmdx.PrintRow(cmd, pat)
			return nil
		},
	}

	cmd.Flags().StringP("name", "n", "", "The name of the project, required when quiet mode is used")
	cmdx.RegisterFormatFlags(cmd.Flags())
	return cmd
}
