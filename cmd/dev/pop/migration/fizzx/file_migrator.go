package fizzx

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

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
		d, err := ioutil.TempDir(os.TempDir(),
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
		content, err := pop.MigrationContent(mf, tx, f, true)
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
			return ioutil.WriteFile(path, contents, 0666)
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "Skipping file: %s\n", path)
			return nil
		}
	}
	_, _ = fmt.Fprintf(os.Stderr, "Wrote file: %s\n", path)
	return ioutil.WriteFile(path, contents, 0666)
}
