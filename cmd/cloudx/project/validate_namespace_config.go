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

The file is checked by the Ory Network project's syntax check endpoint; nothing
is written to the project. The command exits with a non-zero status when the
file has syntax errors, which makes it usable as a CI/CD check before running
` + "`ory update opl`" + `.`,
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

			var errs oplSyntaxErrors
			// An empty file has no syntax errors, and the SDK refuses to send a
			// zero-length body ("invalid body type text/plain").
			if len(data) > 0 {
				result, _, err := c.RelationshipAPI.CheckOplSyntax(ctx).Body(string(data)).Execute()
				if err != nil {
					return cmdx.PrintOpenAPIError(cmd, err)
				}
				// GetErrors is nil-safe: result is nil if the API replies with an empty body.
				errs = oplSyntaxErrors(result.GetErrors())
			}

			switch outputFormat(cmd) {
			case cmdx.FormatJSON, cmdx.FormatJSONPretty, cmdx.FormatJSONPath, cmdx.FormatJSONPointer, cmdx.FormatYAML:
				cmdx.PrintTable(cmd, errs)
			default:
				// Everything else (including an unknown value, which cmdx itself
				// treats as the default format) is rendered for humans.
				source := sourceName(file)
				for _, parseErr := range errs {
					// Syntax errors are what the command is for, so they are also
					// printed with --quiet; only the success message is noise.
					_, _ = fmt.Fprintln(cmd.ErrOrStderr(), formatParseError(source, parseErr))
				}
				if len(errs) == 0 {
					_, _ = cmdx.NewLoudOutPrinter(cmd).Println("The Ory Permission Language file is valid.")
				}
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
	// The column is only appended together with the line, otherwise it would be
	// rendered in the position where readers expect the line number.
	if start := err.Start; start != nil && start.Line != nil {
		location += ":" + strconv.FormatInt(*start.Line, 10)
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
	if e == nil {
		// The field is omitted when nil, but CI/CD consumers should be able to
		// rely on "errors" always being present, e.g. `jq '.errors | length'`.
		e = oplSyntaxErrors{}
	}
	return cloud.CheckOplSyntaxResult{Errors: e}
}

func (e oplSyntaxErrors) Len() int {
	return len(e)
}
