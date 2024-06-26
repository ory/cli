// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package workspace

import (
	cloud "github.com/ory/client-go"
	"github.com/ory/x/cmdx"
)

type (
	outputWorkspace  cloud.Workspace
	outputWorkspaces []cloud.Workspace
)

var (
	workspaceHeader               = []string{"ID", "NAME", "SUBSCRIPTION PLAN"}
	_               cmdx.TableRow = (*outputWorkspace)(nil)
	_               cmdx.Table    = (*outputWorkspaces)(nil)
)

func (*outputWorkspace) Header() []string {
	return workspaceHeader
}

func (w *outputWorkspace) Columns() []string {
	subPlan := cmdx.None
	if w.SubscriptionPlan.Get() != nil {
		subPlan = *w.SubscriptionPlan.Get()
	}
	return []string{w.Id, w.Name, subPlan}
}

func (w *outputWorkspace) Interface() interface{} {
	return w
}

func (o outputWorkspaces) Header() []string {
	return workspaceHeader
}

func (o outputWorkspaces) Table() [][]string {
	res := make([][]string, len(o))
	for k, v := range o {
		res[k] = (*outputWorkspace)(&v).Columns()
	}
	return res
}

func (o outputWorkspaces) Interface() interface{} {
	return o
}

func (o outputWorkspaces) Len() int {
	return len(o)
}
