// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"os"
	"strings"

	"github.com/ory/x/stringsx"
)

func GitHubSHA() string {
	gitHash := stringsx.Coalesce(
		os.Getenv("CIRCLE_SHA1"),
		os.Getenv("GITHUB_SHA"),
		strings.TrimSpace(CommandGetOutput("git", "rev-parse", "HEAD")),
	)
	return gitHash
}

func GitHubTag() string {
	var ghTag string
	if os.Getenv("GITHUB_REF_TYPE") == "tag" {
		ghTag = os.Getenv("GITHUB_REF_NAME")
	}

	tag := stringsx.Coalesce(
		os.Getenv("CIRCLE_TAG"),
		ghTag,
		strings.TrimSpace(CommandGetOutput("git", "tag", "--points-at", "HEAD")),
	)
	return tag
}
