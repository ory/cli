// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package relationtuples

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/client"
	ketoClient "github.com/ory/keto/cmd/client"
	"github.com/ory/keto/cmd/relationtuple"
)

const FlagAll = "all"

var ErrDeleteMissingAllFlag = fmt.Errorf("please select the tuples to delete or use `--all` to delete all tuples")

func NewListCmd() *cobra.Command {
	cmd := relationtuple.NewGetCmd()
	cmd.Short = "List relation tuples"
	cmd.Long = `List relation tuples matching the given partial tuple.
Returns paginated results.`
	wrapForOryCLI(cmd)

	return cmd
}

func NewDeleteCmd() *cobra.Command {
	cmd := relationtuple.NewDeleteAllCmd()
	wrapForOryCLI(cmd)

	cmd.Flags().Bool(FlagAll, false, "Delete all relation tuples")
	wrapForDelete(cmd)

	return cmd
}

func NewCreateCmd() *cobra.Command {
	cmd := relationtuple.NewCreateCmd()
	wrapForOryCLI(cmd)

	return cmd
}

func NewParseCmd() *cobra.Command {
	cmd := relationtuple.NewParseCmd()
	wrapForOryCLI(cmd)

	return cmd
}

func forwardConnectionInfo(cmd *cobra.Command) {
	originalRunE := cmd.RunE
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		h, err := client.NewCobraCommandHelper(cmd)
		if err != nil {
			return err
		}

		project, err := h.GetSelectedProject(ctx)
		if err != nil {
			return err
		}

		tokenSource, baseURL, err := h.ProjectAuthToken(ctx)
		if err != nil {
			return err
		}

		upstream := baseURL(project.Slug + ".projects")
		if upstream.Port() == "" {
			upstream.Host = upstream.Host + ":443"
		}

		token, err := tokenSource.Token()
		if err != nil {
			return err
		}
		_ = os.Setenv(ketoClient.EnvAuthToken, token.AccessToken)
		_ = os.Setenv(ketoClient.EnvReadRemote, upstream.Host)
		_ = os.Setenv(ketoClient.EnvWriteRemote, upstream.Host)

		return originalRunE(cmd, args)
	}
}

func wrapForDelete(cmd *cobra.Command) {
	originalRunE := cmd.RunE
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		all, err := cmd.Flags().GetBool(FlagAll)
		if err != nil {
			return err
		}
		if all {
			return originalRunE(cmd, args)
		}

		// At least one of the query flags must have been set
		queryFlags := []string{
			relationtuple.FlagNamespace,
			relationtuple.FlagObject,
			relationtuple.FlagRelation,
			relationtuple.FlagSubjectID,
			relationtuple.FlagSubjectSet,
		}
		for _, flag := range queryFlags {
			if cmd.Flags().Changed(flag) {
				return originalRunE(cmd, args)
			}
		}
		return ErrDeleteMissingAllFlag
	}
}

func hideKetoFlags(cmd *cobra.Command) {
	for _, flag := range []string{
		ketoClient.FlagReadRemote,
		ketoClient.FlagWriteRemote,
		ketoClient.FlagInsecureNoTransportSecurity,
		ketoClient.FlagInsecureSkipHostVerification,
	} {
		_ = cmd.Flags().MarkHidden(flag)
	}
}

// wrapForOryCLI wraps the Keto command to be used in the ORY CLI.
func wrapForOryCLI(cmd *cobra.Command) {
	cmd.Use = "relationships"
	cmd.Aliases = []string{"relation-tuples", "relationship", "relation-tuple"}
	client.RegisterProjectFlag(cmd.Flags())
	client.RegisterWorkspaceFlag(cmd.Flags())
	forwardConnectionInfo(cmd)
	hideKetoFlags(cmd)
}
