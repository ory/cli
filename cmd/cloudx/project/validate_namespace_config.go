// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	cloud "github.com/ory/client-go"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/osx"
)

// stdinSource is the value of --file that makes the command read from stdin.
const stdinSource = "-"

func NewValidateNamespaceConfigCmd() *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use: "opl",
		Aliases: []string{
			"namespaces-config",
		},
		Args:  cobra.NoArgs,
		Short: "Validate the syntax of an Ory Permission Language file",
		Long: `Validate the syntax of an Ory Permission Language file without applying it.

Ory Network checks the file the same way ` + "`ory update opl`" + ` does before storing it,
but nothing is written to the project. The command exits with a non-zero status
when the file has syntax errors, which makes it usable as a CI/CD check.`,
		Example: `$ {{ .CommandPath }} --file /path/to/namespace_config.ts

The Ory Permission Language file is valid.

$ cat namespace_config.ts | {{ .CommandPath }} --file - --format json

{"errors":[{"end":{"Line":1,"column":36},"message":"expected 'permits' or 'related', got \"\"","start":{"Line":1,"column":36}}]}`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			h, err := client.NewCobraCommandHelper(cmd)
			if err != nil {
				return err
			}

			data, err := readOPLSource(cmd, file)
			if err != nil {
				return err
			}

			c, err := h.NewProjectAPIClient(ctx)
			if err != nil {
				return err
			}

			result, _, err := c.RelationshipAPI.CheckOplSyntax(ctx).Body(string(data)).Execute()
			if err != nil {
				return cmdx.PrintOpenAPIError(cmd, err)
			}

			errs := oplSyntaxErrors(result.Errors)
			switch outputFormat(cmd) {
			case cmdx.FormatDefault, cmdx.FormatTable, cmdx.FormatQuiet:
				for _, parseErr := range errs {
					_, _ = fmt.Fprintln(cmd.ErrOrStderr(), formatParseError(sourceName(file), parseErr))
				}
				if len(errs) == 0 {
					_, _ = cmdx.NewLoudOutPrinter(cmd).Println("The Ory Permission Language file is valid.")
				}
			default:
				cmdx.PrintTable(cmd, errs)
			}

			if len(errs) > 0 {
				return cmdx.FailSilently(cmd)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "",
		`Ory Permission Language file (namespace_config.ts, https://example.org/namespace_config.ts, "-" for stdin, ...) to validate`)
	_ = cmd.MarkFlagRequired("file")
	client.RegisterProjectFlag(cmd.Flags())
	client.RegisterWorkspaceFlag(cmd.Flags())

	return cmd
}

// readOPLSource reads the OPL file from stdin when source is "-", and from any
// source supported by osx otherwise.
func readOPLSource(cmd *cobra.Command, source string) ([]byte, error) {
	if source == stdinSource {
		return io.ReadAll(cmd.InOrStdin())
	}
	return osx.ReadFileFromAllSources(source)
}

func sourceName(source string) string {
	if source == stdinSource {
		return "stdin"
	}
	return source
}

// outputFormat mirrors how cmdx picks the output format, so that the
// human-readable rendering below is only used for the human-readable formats.
func outputFormat(cmd *cobra.Command) cmdx.Format {
	if quiet, err := cmd.Flags().GetBool(cmdx.FlagQuiet); err == nil && quiet {
		return cmdx.FormatQuiet
	}
	format, err := cmd.Flags().GetString(cmdx.FlagFormat)
	if err != nil {
		return cmdx.FormatDefault
	}
	// jsonpath and jsonpointer carry their argument in the flag value.
	name, _, _ := strings.Cut(format, "=")
	return cmdx.Format(name)
}

// formatParseError renders a parse error the way compilers do, so that editors
// and CI log parsers can pick up the location.
func formatParseError(source string, err cloud.ParseError) string {
	location := source
	if start := err.Start; start != nil {
		if start.Line != nil {
			location += ":" + strconv.FormatInt(*start.Line, 10)
		}
		if start.Column != nil {
			location += ":" + strconv.FormatInt(*start.Column, 10)
		}
	}
	return fmt.Sprintf("%s: %s", location, err.GetMessage())
}

type oplSyntaxErrors []cloud.ParseError

var _ cmdx.Table = (oplSyntaxErrors)(nil)

func (oplSyntaxErrors) Header() []string {
	return []string{"LINE", "COLUMN", "MESSAGE"}
}

func (e oplSyntaxErrors) Table() [][]string {
	rows := make([][]string, len(e))
	for i, parseErr := range e {
		line, column := cmdx.None, cmdx.None
		if start := parseErr.Start; start != nil {
			if start.Line != nil {
				line = strconv.FormatInt(*start.Line, 10)
			}
			if start.Column != nil {
				column = strconv.FormatInt(*start.Column, 10)
			}
		}
		rows[i] = []string{line, column, parseErr.GetMessage()}
	}
	return rows
}

func (e oplSyntaxErrors) Interface() interface{} {
	return cloud.CheckOplSyntaxResult{Errors: e}
}

func (e oplSyntaxErrors) Len() int {
	return len(e)
}
