// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import cloud "github.com/ory/client-go"

type MockCommandHelper struct {
	Project *cloud.Project
}

var _ Command = new(MockCommandHelper)

func (n *MockCommandHelper) GetProject(projectId string) (*cloud.Project, error) {
	return n.Project, nil
}
