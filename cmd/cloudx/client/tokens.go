// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"os"
)

// GetProjectAPIKeyFromEnvironment returns the project API key from the environment variable.
func GetProjectAPIKeyFromEnvironment() string {
	return os.Getenv("ORY_API_KEY")
}
