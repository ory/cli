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

	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"
)

type DumpMigrator struct {
	pop.Migrator
	Path string
}

func NewDumpMigrator(path string, dest string, shouldReplace, dumpSchema bool, c *pop.Connection) (DumpMigrator, error) {
	fm := DumpMigrator{
		Migrator: pop.NewMigrator(c),
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

	runner := func(mf pop.Migration, tx *pop.Connection) error {
		f, err := os.Open(mf.Path)
		if err != nil {
			return err
		}
		defer f.Close()
		content, err := MigrationContent(mf, tx, f, true)
		if err != nil {
			return errors.Wrapf(err, "error processing %s", mf.Path)
		}
		if content == "" {
			return nil
		}

		_, fn := filepath.Split(mf.Path)
		fn = strings.Replace(fn, ".up.fizz", fmt.Sprintf(".%s.up.sql", tx.Dialect.Name()), -1)
		fn = strings.Replace(fn, ".down.fizz", fmt.Sprintf(".%s.down.sql", tx.Dialect.Name()), -1)
		if err := writeFile(filepath.Join(dest, fn), []byte(content), shouldReplace); err != nil {
			return err
		}

		err = tx.RawQuery(content).Exec()
		if err != nil {
			return errors.Wrapf(err, "error executing %s, sql: %s", mf.Path, content)
		}
		return nil
	}

	err := fm.findMigrations(runner)
	if err != nil {
		return fm, err
	}

	return fm, nil
}

func (fm *DumpMigrator) findMigrations(runner func(mf pop.Migration, tx *pop.Connection) error) error {
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
			mf := pop.Migration{
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

func MigrationContent(mf pop.Migration, c *pop.Connection, r io.Reader, usingTemplate bool) (string, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return "", nil
	}

	content := ""
	if usingTemplate {
		t, err := template.New("migration").Parse(string(b))
		if err != nil {
			return "", errors.Wrapf(err, "could not parse template %s", mf.Path)
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
			return "", errors.Wrapf(err, "could not execute migration template %s", mf.Path)
		}
		content = bb.String()
	} else {
		content = string(b)
	}

	if mf.Type == "fizz" {
		content, err = fizz.AString(content, c.Dialect.FizzTranslator())
		if err != nil {
			return "", errors.Wrapf(err, "could not fizz the migration %s", mf.Path)
		}
	}

	return content, nil
}
