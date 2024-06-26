// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package project_test

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ory/cli/cmd/cloudx/testhelpers"
)

type execFunc func(stdin io.Reader, args ...string) (string, string, error)

type RunWithProjectOption uint8

const (
	WithDefaultProject RunWithProjectOption = 1 << iota
	WithPositionalProject
	WithFlagProject
)

func runWithProject(t *testing.T, test func(t *testing.T, exec execFunc, projectID string), opts ...RunWithProjectOption) {
	for _, v := range opts {
		switch v {
		case WithDefaultProject:
			t.Run("project via configured default", func(t *testing.T) {
				testhelpers.SetDefaultProject(t, defaultConfig, extraProject.Id)

				test(t, func(stdin io.Reader, args ...string) (string, string, error) {
					return defaultCmd.Exec(stdin, args...)
				}, extraProject.Id)

				// make sure, the default wasn't changed implicitly
				assert.Equal(t, extraProject.Id, testhelpers.GetDefaultProjectID(t, defaultConfig))
			})
		case WithPositionalProject:
			t.Run("explicit project via positional argument", func(t *testing.T) {
				testhelpers.SetDefaultProject(t, defaultConfig, defaultProject.Id)

				test(t, func(stdin io.Reader, args ...string) (string, string, error) {
					return defaultCmd.Exec(stdin, append(args, extraProject.Id)...)
				}, extraProject.Id)

				// make sure, the default wasn't changed implicitly
				assert.Equal(t, defaultProject.Id, testhelpers.GetDefaultProjectID(t, defaultConfig))
			})
		case WithFlagProject:
			t.Run("explicit project via `--project` flag", func(t *testing.T) {
				testhelpers.SetDefaultProject(t, defaultConfig, defaultProject.Id)

				test(t, func(stdin io.Reader, args ...string) (string, string, error) {
					return defaultCmd.Exec(stdin, append(args, "--project", extraProject.Id)...)
				}, extraProject.Id)

				// make sure, the default wasn't changed implicitly
				assert.Equal(t, defaultProject.Id, testhelpers.GetDefaultProjectID(t, defaultConfig))
			})
		}
	}
}
