package pkg

import (
	"os"
	"strings"
)

func CircleSHA1() string {
	gitHash := os.Getenv("CIRCLE_SHA1")
	if len(gitHash) == 0 {
		gitHash = strings.TrimSpace(CommandGetOutput("git", "rev-parse", "HEAD"))
		if len(gitHash) > 16 {
			gitHash = gitHash[:16]
		}
	}
	return gitHash
}

func CircleTag() string {
	circleTag := os.Getenv("CIRCLE_TAG")
	if len(circleTag) == 0 {
		circleTag = strings.TrimSpace(CommandGetOutput("git", "tag", "--points-at", "HEAD"))
	}
	return circleTag
}
