package accountexperience

import (
	"github.com/ory/cli/cmd/cloudx/client"
	br "github.com/pkg/browser"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

const project = "project"

func RegisterProjectFlag(f *flag.FlagSet) {
	f.String(project, "", "The project to use")
}

func NewAccountExperienceOpenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account-experience",
		Short: "Open Ory Account Experience Pages",
	}
	cmd.AddCommand(NewLoginCmd())

	return cmd
}
func NewLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "login",
		Aliases: []string{"registration", "settings"},
		Short:   "Open Ory Account Experience Pages",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, p, err := client.Client(cmd)
			if err != nil {
				return err
			}

			url := "https://" + p.GetSlug() + ".projects.oryapis.com/ui/" + cmd.CalledAs()

			err = br.OpenURL(url)
			if err != nil {
				return err
			}

			return nil
		},
	}
	return cmd
}
