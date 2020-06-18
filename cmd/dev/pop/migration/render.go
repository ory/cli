package migration

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/gobuffalo/pop/v5"
	"github.com/ory/x/flagx"
	"github.com/ory/x/randx"
	"github.com/ory/x/sqlcon/dockertest"
	"github.com/spf13/cobra"

	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/avast/retry-go"

	"github.com/ory/cli/cmd/dev/pop/migration/fizzx"
	"github.com/ory/cli/cmd/pkg"
)

var render = &cobra.Command{
	Use:   "render [path/to/fizz-templates] [path/to/output]",
	Short: "Renders all fizz templates to their SQL counterparts",
	Long: `This command takes fizz migration templates and renders them as SQL.

It currently supports MySQL, SQLite, PostgreSQL, and CockroachDB (SQL). To use this tool you need Docker installed.
`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		defer dockertest.KillAllTestDatabases()

		var l sync.Mutex
		dsns := map[string]string{
			"sqlite": "sqlite3://" + filepath.Join(os.TempDir(), randx.MustString(12, randx.AlphaNum)) + ".sql?mode=memory&_fk=true"}

		dockertest.Parallel([]func(){
			func() {
				u, err := dockertest.RunPostgreSQL()
				pkg.Check(err)
				l.Lock()
				dsns["postgres"] = u
				l.Unlock()
			},
			func() {
				u, err := dockertest.RunMySQL()
				pkg.Check(err)
				l.Lock()
				dsns["mysql"] = u
				l.Unlock()
			},
			func() {
				u, err := dockertest.RunCockroachDB()
				pkg.Check(err)
				l.Lock()
				dsns["cockroach"] = u
				l.Unlock()
			},
		})

		pkg.Check(os.MkdirAll(args[1], 0777))

		// Ensure a connection exists and works before running the translators.
		for _, dsn := range dsns {
			c, err := pop.NewConnection(&pop.ConnectionDetails{URL: dsn})
			pkg.Check(err)

			pkg.Check(retry.Do(func() error {
				if err := c.Open(); err != nil {
					return err
				}
				return c.RawQuery("SELECT 1").Exec()
			}))

			m, err := fizzx.NewFileMigrator(args[0], args[1], flagx.MustGetBool(cmd, "replace"), c)
			pkg.Check(err)

			pkg.Check(m.Up())
			pkg.Check(m.Down(-1))
			pkg.Check(c.Close())
		}

		return nil
	},
}

func init() {
	Main.AddCommand(render)

	render.Flags().BoolP("replace", "r", false, "Replaces existing files if set.")
}
