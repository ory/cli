package migration

import (
	"fmt"
	"os"
	"time"

	"github.com/ory/x/flagx"
	"github.com/ory/x/stringslice"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:  "create [name]",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		prefix := time.Now().Format("20060102150405")
		d := flagx.MustGetString(cmd, "dialect")

		suffix := ".fizz"
		if len(d) > 0 {
			if stringslice.Has(supportedDialects, d) {
				return fmt.Errorf(`expected dialect to be one of %v but got: %s`, supportedDialects, d)
			}
			suffix = fmt.Sprintf(".%s.sql", d)
		}

		f, err := os.Create(fmt.Sprintf("%s_%s%s", prefix, args[0], suffix))
		if err != nil {
			return err
		}

		return f.Close()
	},
}

var supportedDialects = []string{
	"sqlite",
	"cockroach",
	"mysql",
	"postgres",
}

func init() {
	Main.AddCommand(createCmd)

	createCmd.Flags().StringP("dialect", "d", "", fmt.Sprintf(`Choose a dialect from: %v`, supportedDialects))
}
