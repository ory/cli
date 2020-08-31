package migration

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/ory/x/flagx"
	"github.com/ory/x/stringslice"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:  "create [destination] [name]",
	Args: cobra.ExactArgs(2),
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

		for _, fn := range []string{
			fmt.Sprintf("%s_%s.up%s", prefix, args[1], suffix),
			fmt.Sprintf("%s_%s.down%s", prefix, args[1], suffix),
		} {
			if err := ioutil.WriteFile(filepath.Join(args[0], fn), []byte{}, 0666); err != nil {
				return err
			}
		}
		return nil
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
