// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fizzx

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gobuffalo/fizz"

	"github.com/ory/x/logrusx"
	"github.com/ory/x/stringsx"

	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"
)

type DumpMigrator struct {
	Migrator
	Path string
}

func NewDumpMigrator(path string, dest string, shouldReplace, dumpSchema bool, c *pop.Connection, l *logrusx.Logger) (DumpMigrator, error) {
	fm := DumpMigrator{
		Migrator: NewMigrator(c, l),
		Path:     path,
	}

	if dumpSchema {
		d, err := os.MkdirTemp(os.TempDir(),
			fmt.Sprintf("schema-%s-*", c.Dialect.Name()))
		if err != nil {
			return fm, err
		}
		fm.SchemaPath = d
	}

	runner := func(mf Migration) ([]MigrationTuple, error) {
		f, err := os.Open(mf.Path)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		content, err := MigrationContent(mf, c, f, true)
		if err != nil {
			return nil, errors.Wrapf(err, "error processing %s", mf.Path)
		}

		if len(content) == 0 {
			return nil, nil
		}

		var tuples []MigrationTuple
		for k, statement := range content {
			if strings.TrimSpace(statement) == "" {
				continue
			}
			statement := strings.TrimSuffix(statement, ";") + ";"

			id := fmt.Sprintf("%s%06d", mf.Version, len(content)-1-k)
			if mf.Direction == "up" {
				id = fmt.Sprintf("%s%06d", mf.Version, k)
			}

			fn := fmt.Sprintf("%s_%s.%s.%s.sql", id, mf.Name, c.Dialect.Name(), mf.Direction)
			if err := writeFile(filepath.Join(dest, fn), []byte(statement), shouldReplace); err != nil {
				return nil, err
			}

			placeholder := fmt.Sprintf("%s_%s.%s.%s.sql", id, mf.Name, c.Dialect.Name(), "up")
			if mf.Direction == "up" {
				placeholder = fmt.Sprintf("%s_%s.%s.%s.sql", id, mf.Name, c.Dialect.Name(), "down")
			}

			placeholder = filepath.Join(dest, placeholder)
			if _, err := os.Stat(placeholder); os.IsNotExist(err) {
				l.WithField("file", placeholder).Info("Writing filler file.")
				if err := writeFile(placeholder, []byte{}, shouldReplace); err != nil {
					return nil, err
				}
			}

			tuples = append(tuples, MigrationTuple{
				ID:        id,
				Statement: statement + ";",
			})
		}

		return tuples, nil
	}

	err := fm.findMigrations(runner)
	if err != nil {
		return fm, err
	}

	return fm, nil
}

func (fm *DumpMigrator) findMigrations(runner func(mf Migration) ([]MigrationTuple, error)) error {
	dir := fm.Path
	if fi, err := os.Stat(dir); err != nil || !fi.IsDir() {
		// directory doesn't exist
		return nil
	}
	return filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			match, err := pop.ParseMigrationFilename(info.Name())
			if err != nil {
				return err
			}
			if match == nil {
				return nil
			}
			mf := Migration{
				Path:      p,
				Version:   match.Version,
				Name:      match.Name,
				DBType:    match.DBType,
				Direction: match.Direction,
				Type:      match.Type,
				Runner:    runner,
			}
			fm.Migrations[mf.Direction] = append(fm.Migrations[mf.Direction], mf)
		}
		return nil
	})
}

func writeFile(path string, contents []byte, replace bool) error {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		if replace {
			_, _ = fmt.Fprintf(os.Stderr, "Wrote file: %s\n", path)
			return os.WriteFile(path, contents, 0666)
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "Skipping file: %s\n", path)
			original, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if string(contents) != string(original) {
				_, _ = fmt.Fprintf(os.Stderr, `Migrations are not equal!

Expected:

%s

------------------------
Actual:

%s

`, original, contents)
				return errors.Errorf("migrations are not equal")
			}
			return nil
		}
	}
	_, _ = fmt.Fprintf(os.Stderr, "Wrote file: %s\n", path)
	return os.WriteFile(path, contents, 0666)
}

func MigrationContent(mf Migration, c *pop.Connection, r io.Reader, usingTemplate bool) ([]string, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, nil
	}

	content := ""
	if usingTemplate {
		t, err := template.New("migration").Parse(string(b))
		if err != nil {
			return nil, errors.Wrapf(err, "could not parse template %s", mf.Path)
		}
		var bb bytes.Buffer
		err = t.Execute(&bb, struct {
			IsSQLite     bool
			IsCockroach  bool
			IsMySQL      bool
			IsMariaDB    bool
			IsPostgreSQL bool
		}{
			IsSQLite:     c.Dialect.Name() == "sqlite3",
			IsCockroach:  c.Dialect.Name() == "cockroach",
			IsMySQL:      c.Dialect.Name() == "mysql",
			IsMariaDB:    c.Dialect.Name() == "mariadb",
			IsPostgreSQL: c.Dialect.Name() == "postgres",
		})
		if err != nil {
			return nil, errors.Wrapf(err, "could not execute migration template %s", mf.Path)
		}
		content = bb.String()
	} else {
		content = string(b)
	}

	if mf.Type == "fizz" {
		content, err = fizz.AString(content, c.Dialect.FizzTranslator())
		if err != nil {
			return nil, errors.Wrapf(err, "could not fizz the migration %s", mf.Path)
		}
	}

	content = strings.ReplaceAll(content, "COMMIT TRANSACTION;", "")
	content = strings.ReplaceAll(content, "BEGIN TRANSACTION;", "")

	return stringsx.Splitx(content, ";\n"), nil
}
