package migration

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/ory/x/randx"
	"github.com/ory/x/stringslice"

	"github.com/go-sql-driver/mysql"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/pop/v5/logging"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/x/flagx"
	"github.com/ory/x/sqlcon/dockertest"
	"github.com/spf13/cobra"

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

		// Disable log outputs
		pop.SetLogger(func(lvl logging.Level, s string, args ...interface{}) {})
		_ = mysql.SetLogger(log.New(ioutil.Discard, "", 0))

		var l sync.Mutex
		dialects := flagx.MustGetStringSlice(cmd, "dialects")
		dsns := map[string]string{}

		if stringslice.Has(dialects, "sqlite") {
			dsns["sqlite"] = "sqlite3://" + filepath.Join(os.TempDir(), randx.MustString(12, randx.AlphaNum)) + ".sql?mode=memory&_fk=true"
		}

		dockertest.Parallel([]func(){
			func() {
				if stringslice.Has(dialects, "postgres") {
					u, err := dockertest.RunPostgreSQL()
					pkg.Check(err)
					l.Lock()
					dsns["postgres"] = u
					l.Unlock()
				}
			},
			func() {
				if stringslice.Has(dialects, "mysql") {
					u, err := dockertest.RunMySQL()
					pkg.Check(err)
					l.Lock()
					dsns["mysql"] = u
					l.Unlock()
				}
			},
			func() {
				if stringslice.Has(dialects, "cockroach") {
					u, err := dockertest.RunCockroachDB()
					pkg.Check(err)
					l.Lock()
					dsns["cockroach"] = u
					l.Unlock()
				}
			},
		})

		if len(dsns) == 0 {
			panic(fmt.Sprintf("Expected at least one dialect out of [sqlite, mysql, postgres, cockroach], but got %v", dialects))
		}

		pkg.Check(os.MkdirAll(args[1], 0777))

		dump := flagx.MustGetBool(cmd, "dump")
		replace := flagx.MustGetBool(cmd, "replace")

		var wg sync.WaitGroup
		runner := func(name, dsn string) {
			defer wg.Done()
			c, err := pop.NewConnection(&pop.ConnectionDetails{URL: dsn})
			pkg.Check(err)

			pkg.Check(retry.Do(func() error {
				if err := c.Open(); err != nil {
					return err
				}
				return c.RawQuery("SELECT 1").Exec()
			}))

			m, err := fizzx.NewDumpMigrator(args[0], args[1], replace, dump, c)
			pkg.Check(err)

			pkg.Check(m.Up())

			if dump {
				_ = m.DumpMigrationSchema()
				_, _ = fmt.Fprintf(os.Stderr, "Dumped %s schema to: %s\n", name, m.SchemaPath)
			}

			pkg.Check(m.Down(-1))
			pkg.Check(c.Close())
		}

		wg.Add(len(dsns))
		// Ensure a connection exists and works before running the translators.
		for name, dsn := range dsns {
			go runner(name, dsn)
		}

		wg.Wait()
		return nil
	},
}

func init() {
	Main.AddCommand(render)

	render.Flags().BoolP("replace", "r", false, "Replaces existing files if set.")
	render.Flags().BoolP("dump", "d", false, "If set dumps the schema to a temporary location.")
	render.Flags().StringSlice("dialects", []string{"sqlite", "mysql", "postgres", "cockroach"}, "Select dialects to render. Comma separated list out of sqlite,mysql,postgres,cockroach.")
}
