package fizzx

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/gobuffalo/fizz"
	"github.com/gobuffalo/pop/v5"

	"github.com/ory/x/logrusx"
	"github.com/ory/x/stringslice"

	"github.com/pkg/errors"
)

var mrx = regexp.MustCompile(`^(\d+)_([^.]+)(\.[a-z0-9]+)?\.(up|down)\.(sql|fizz)$`)

type MigrationTuple struct {
	ID        string
	Statement string
}

// NewMigrator returns a new "blank" migrator. It is recommended
// to use something like MigrationBox or FileMigrator. A "blank"
// Migrator should only be used as the basis for a new type of
// migration system.
func NewMigrator(c *pop.Connection, l *logrusx.Logger) Migrator {
	return Migrator{
		Connection: c,
		l:          l.WithField("package", "github.com/ory/popx").WithField("component", "migrator").WithField("database", c.Dialect.Name()),
		Migrations: map[string]Migrations{
			"up":   {},
			"down": {},
		},
	}
}

// Migrator forms the basis of all migrations systems.
// It does the actual heavy lifting of running migrations.
// When building a new migration system, you should embed this
// type into your migrator.
type Migrator struct {
	l          *logrusx.Logger
	Connection *pop.Connection
	SchemaPath string
	Migrations map[string]Migrations
}

func (m Migrator) migrationIsCompatible(dialect string, mi Migration) bool {
	if mi.DBType == "all" || mi.DBType == dialect {
		return true
	}
	return false
}

// Up runs pending "up" migrations and applies them to the database.
func (m Migrator) Up() error {
	_, err := m.UpTo(0)
	return err
}

// UpTo runs up to step "up" migrations and applies them to the database.
// If step <= 0 all pending migrations are run.
func (m Migrator) UpTo(step int) (applied int, err error) {
	c := m.Connection
	err = m.exec(func() error {
		mtn := c.MigrationTableName()
		mfs := m.Migrations["up"]
		mfs.Filter(func(mf Migration) bool {
			return m.migrationIsCompatible(c.Dialect.Name(), mf)
		})
		sort.Sort(mfs)

		alreadyRan := make([]string, 0, len(mfs))

		for _, mi := range mfs {
			if stringslice.Has(alreadyRan, mi.Version) {
				continue
			}

			tuples, err := mi.Run()
			if err != nil {
				return err
			}

			alreadyRan = append(alreadyRan, mi.Version)

			for _, tuple := range tuples {
				m.l.WithField("sql", tuple.Statement).WithField("id", tuple.ID).Debug("Trying to execute SQL up migration.")

				exists, err := c.Where("version = ?", tuple.ID).Exists(mtn)
				if err != nil {
					return errors.Wrapf(err, "problem checking for migration version %s", tuple.ID)
				}

				if exists {
					continue
				}

				if err := c.Transaction(func(tx *pop.Connection) error {
					if tuple.Statement != "" {
						if err := tx.RawQuery(tuple.Statement).Exec(); err != nil {
							c.Select("SHOW")
							return errors.Wrapf(err, "unable to execute migration %s", tuple.ID)
						}
					}

					return errors.Wrapf(
						tx.RawQuery(fmt.Sprintf("INSERT INTO %s (version) VALUES (?)", mtn), tuple.ID).Exec(),
						"problem deleting migration version %s", tuple.ID)
				}); err != nil {
					return err
				}
			}

			m.l.Debugf("> %s", mi.Name)
			applied++
			if step > 0 && applied >= step {
				break
			}
		}
		if applied == 0 {
			m.l.Debugf("Migrations already up to date, nothing to apply")
		} else {
			m.l.Debugf("Successfully applied %d migrations.", applied)
		}
		return nil
	})
	return
}

// Down runs pending "down" migrations and rolls back the
// database by the specified number of steps.
func (m Migrator) Down(step int) error {
	c := m.Connection
	return m.exec(func() error {
		mtn := c.MigrationTableName()
		count, err := c.Count(mtn)
		if err != nil {
			return errors.Wrap(err, "migration down: unable count existing migration")
		}
		mfs := m.Migrations["down"]
		mfs.Filter(func(mf Migration) bool {
			return m.migrationIsCompatible(c.Dialect.Name(), mf)
		})
		sort.Sort(sort.Reverse(mfs))
		// skip all ran migration
		if len(mfs) > count {
			mfs = mfs[len(mfs)-count:]
		}
		// run only required steps
		if step > 0 && len(mfs) >= step {
			mfs = mfs[:step]
		}

		alreadyRan := make([]string, 0, len(mfs))

		for _, mi := range mfs {
			if stringslice.Has(alreadyRan, mi.Version) {
				continue
			}

			tuples, err := mi.Run()
			if err != nil {
				return err
			}

			alreadyRan = append(alreadyRan, mi.Version)

			for _, tuple := range tuples {
				m.l.WithField("sql", tuple.Statement).WithField("id", tuple.ID).Debug("Trying to execute SQL down migration.")

				exists, err := c.Where("version = ?", tuple.ID).Exists(mtn)
				if err != nil {
					return errors.Wrapf(err, "problem checking for migration version %s", tuple.ID)
				}

				if !exists {
					m.l.WithField("migration", tuple.ID).WithField("file", mi.Name).Warn("Migration was not found, but this is ok as down migrations might need more statements.")
				}

				if err := c.Transaction(func(tx *pop.Connection) error {
					if tuple.Statement != "" {
						if err := tx.RawQuery(tuple.Statement).Exec(); err != nil {
							return errors.Wrapf(err, "unable to execute migration %s", tuple.ID)
						}
					}

					return errors.Wrapf(
						tx.RawQuery(fmt.Sprintf("DELETE FROM %s WHERE version = ?", mtn), tuple.ID).Exec(),
						"problem deleting migration version %s", tuple.ID)
				}); err != nil {
					return err
				}
			}

			m.l.Debugf("< %s", mi.Name)
		}
		return nil
	})
}

// Reset the database by running the down migrations followed by the up migrations.
func (m Migrator) Reset() error {
	err := m.Down(-1)
	if err != nil {
		return err
	}
	return m.Up()
}

// CreateSchemaMigrations sets up a table to track migrations. This is an idempotent
// operation.
func CreateSchemaMigrations(c *pop.Connection) error {
	mtn := c.MigrationTableName()
	err := c.Open()
	if err != nil {
		return errors.Wrap(err, "could not open connection")
	}
	_, err = c.Store.Exec(fmt.Sprintf("select * from %s", mtn))
	if err == nil {
		return nil
	}

	return c.Transaction(func(tx *pop.Connection) error {
		schemaMigrations := newSchemaMigrations(mtn)
		smSQL, err := c.Dialect.FizzTranslator().CreateTable(schemaMigrations)
		if err != nil {
			return errors.Wrap(err, "could not build SQL for schema migration table")
		}
		err = tx.RawQuery(smSQL).Exec()
		if err != nil {
			return errors.Wrap(err, smSQL)
		}
		return nil
	})
}

// CreateSchemaMigrations sets up a table to track migrations. This is an idempotent
// operation.
func (m Migrator) CreateSchemaMigrations() error {
	return CreateSchemaMigrations(m.Connection)
}

// Status prints out the status of applied/pending migrations.
func (m Migrator) Status(out io.Writer) error {
	err := m.CreateSchemaMigrations()
	if err != nil {
		return err
	}
	w := tabwriter.NewWriter(out, 0, 0, 3, ' ', tabwriter.TabIndent)
	_, _ = fmt.Fprintln(w, "Version\tName\tStatus\t")
	for _, mf := range m.Migrations["up"] {
		exists, err := m.Connection.Where("version = ?", mf.Version).Exists(m.Connection.MigrationTableName())
		if err != nil {
			return errors.Wrapf(err, "problem with migration")
		}
		state := "Pending"
		if exists {
			state = "Applied"
		}
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t\n", mf.Version, mf.Name, state)
	}
	return w.Flush()
}

// DumpMigrationSchema will generate a file of the current database schema
// based on the value of Migrator.SchemaPath
func (m Migrator) DumpMigrationSchema() error {
	if m.SchemaPath == "" {
		return nil
	}
	c := m.Connection
	schema := filepath.Join(m.SchemaPath, "schema.sql")
	f, err := os.Create(schema)
	if err != nil {
		return err
	}
	err = c.Dialect.DumpSchema(f)
	if err != nil {
		m.l.WithError(err).Error("Unable to dump schema.")
		os.RemoveAll(schema)
		return err
	}
	m.l.Infof("Dumped migration schema to: %s", schema)
	return nil
}

func (m Migrator) exec(fn func() error) error {
	now := time.Now()
	defer func() {
		err := m.DumpMigrationSchema()
		if err != nil {
			m.l.WithError(err).Error("Migrator: unable to dump schema")
		}
	}()
	defer m.printTimer(now)

	err := m.CreateSchemaMigrations()
	if err != nil {
		return errors.Wrap(err, "migrator: problem creating schema migrations")
	}
	return fn()
}

func (m Migrator) printTimer(timerStart time.Time) {
	diff := time.Since(timerStart).Seconds()
	if diff > 60 {
		m.l.Debugf("%.4f minutes", diff/60)
	} else {
		m.l.Debugf("%.4f seconds", diff)
	}
}

func newSchemaMigrations(name string) fizz.Table {
	return fizz.Table{
		Name: name,
		Columns: []fizz.Column{
			{
				Name:    "version",
				ColType: "string",
				Options: map[string]interface{}{
					"size": 48, // leave some extra room
				},
			},
		},
		Indexes: []fizz.Index{},
	}
}
