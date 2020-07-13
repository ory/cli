package release

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ory/x/flagx"
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/pkg"
)

var compile = &cobra.Command{
	Use:   "compile",
	Args:  cobra.ExactArgs(0),
	Short: "Compiles the current project using oryd/xgoreleaser a new release",
	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		pkg.Check(err)

		pkg.Check(pkg.NewCommand("docker", "run", "--mount",
			fmt.Sprintf(`type=bind,source=%s,target=/project`, wd),
			"oryd/xgoreleaser:"+flagx.MustGetString(cmd, "tag"),
			"--timeout", "60m",
			"--skip-publish", "--snapshot", "--rm-dist", "--parallelism",
			strconv.Itoa(flagx.MustGetInt(cmd, "parallelism"))).Run())
	},
}

func init() {
	Main.AddCommand(compile)
	compile.Flags().StringP("tag", "t", "1.14.4-0.139.0", "Set the xgoreleaser version tag.")
	compile.Flags().IntP("parallelism", "p", 4, "Build parallelism.")
}
