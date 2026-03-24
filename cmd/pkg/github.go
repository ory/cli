// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"cmp"
	"os"
	"strings"
)

func GitHubSHA() string {
	gitHash := cmp.Or(
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

	tag := cmp.Or(
		os.Getenv("CIRCLE_TAG"),
		ghTag,
		strings.TrimSpace(CommandGetOutput("git", "tag", "--points-at", "HEAD")),
	)
	return tag
}
