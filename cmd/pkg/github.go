package pkg

import (
	"os"
	"strings"
)

func GitHubSHA() string {
	gitHash := os.Getenv("GITHUB_SHA")
	if len(gitHash) == 0 {
		gitHash = strings.TrimSpace(CommandGetOutput("git", "rev-parse", "HEAD"))
		if len(gitHash) > 16 {
			gitHash = gitHash[:16]
		}
	}
	return gitHash
}

func GitHubTag() string {
	var ghTag string
	if os.Getenv("GITHUB_REF_TYPE") == "tag" {
		ghTag = os.Getenv("GITHUB_REF_NAME")
		if len(ghTag) == 0 {
			ghTag = strings.TrimSpace(CommandGetOutput("git", "tag", "--points-at", "HEAD"))
		}
	}
	return ghTag
}
