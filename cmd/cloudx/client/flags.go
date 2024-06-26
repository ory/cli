// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"github.com/spf13/pflag"
)

const (
	FlagWorkspace = "workspace"
	FlagProject   = "project"
	FlagYes       = "yes"
)

func RegisterWorkspaceFlag(f *pflag.FlagSet) {
	f.String(FlagWorkspace, "", "The workspace to use, either workspace ID or a (partial) name.")
}

func RegisterProjectFlag(f *pflag.FlagSet) {
	f.String(FlagProject, "", "The project to use, either project ID or a (partial) slug.")
}

func RegisterYesFlag(f *pflag.FlagSet) {
	f.BoolP(FlagYes, FlagYes[:1], false, "Confirm all dialogs with yes.")
}
