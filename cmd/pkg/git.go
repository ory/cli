// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Masterminds/semver/v3"
)

const GitCommitMessagePreviousVersion = "Bumps from"

// stdin is shared across all prompts. Creating a fresh bufio.Reader per prompt
// discards read-ahead input, which breaks piped stdin and typed-ahead answers.
var stdin = bufio.NewReader(os.Stdin)

func NewCommand(name string, args ...string) *exec.Cmd {
	_, _ = fmt.Fprintf(os.Stderr, "$ %s %s\n", name, strings.Join(args, " "))
	ec := exec.Command(name, args...)
	ec.Stdout = os.Stdout
	ec.Stderr = os.Stderr
	// ec.Stdin = os.Stdin
	return ec
}

func NewCommandIn(dir, name string, args ...string) *exec.Cmd {
	ec := NewCommand(name, args...)
	ec.Dir = dir
	return ec
}

func GitTagRelease(dir string, annotate, dry bool, nextVersion semver.Version, previousVersion *semver.Version) {
	Check(NewCommandIn(dir, "git", "add", "-A").Run())

	gitArgs := []string{"commit", "-a", "--allow-empty", "-m",
		fmt.Sprintf("autogen: pin v%s release commit", nextVersion.String())}
	if previousVersion != nil {
		gitArgs = append(gitArgs, "-m", fmt.Sprintf("%s v%s", GitCommitMessagePreviousVersion, previousVersion.String()))
	}

	Check(NewCommandIn(dir, "git", gitArgs...).Run())

	if annotate {
		message := promptTagMessage(nextVersion)
		Check(NewCommandIn(dir, "git", "tag", fmt.Sprintf("v%s", nextVersion.String()), "-a", "-m", message).Run())
	} else {
		Check(NewCommandIn(dir, "git", "tag", fmt.Sprintf("v%s", nextVersion.String())).Run())
	}

	if !dry {
		Confirm("Pressing [y] will push this (%s) release to GitHub. Are you sure?", dir)
		Check(NewCommandIn(dir, "git", "push").Run())
		Check(NewCommandIn(dir, "git", "push", "--tags").Run())
	}
}

func GitClone(repo string) string {
	dest, err := os.MkdirTemp(os.TempDir(), "ory-release-*")
	Check(err)
	Check(NewCommand("git", "clone", repo, dest).Run())
	return dest
}

// promptTagMessage reads the annotated tag message from stdin. Opening
// $EDITOR instead is not safe here: handing the inherited terminal to vim
// mid-run can leave the terminal in a broken state and subsequent stdin
// reads fail with EOF.
func promptTagMessage(version semver.Version) string {
	fmt.Printf("Enter the tag message for v%s. Finish with an empty line:\n> ", version.String())
	var lines []string
	for {
		line, err := stdin.ReadString('\n')
		Check(err)

		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			if len(lines) > 0 {
				return strings.Join(lines, "\n")
			}
			fmt.Print("The tag message must not be empty.\n> ")
			continue
		}
		lines = append(lines, line)
		fmt.Print("> ")
	}
}

func Confirm(message string, args ...interface{}) {
	for {
		fmt.Printf("%s [y/n] ", fmt.Sprintf(message, args...))
		answer, err := stdin.ReadString('\n')
		Check(err)

		answer = strings.TrimSpace(answer)
		if answer == "n" {
			Fatalf("Aborting because your answer was: %s", answer)
		} else if answer != "y" {
			continue
		} else {
			// answer is 'y'
			break
		}
	}
}

func GitListTags() string {
	return CommandGetOutput("git", "tag", "--sort=creatordate")
}

func CommandGetOutput(name string, args ...string) string {
	var b bytes.Buffer
	cmd := NewCommand(name, args...)
	cmd.Stdout = &b
	Check(cmd.Run())
	return b.String()
}

func GitGetCurrentTag() string {
	var b bytes.Buffer
	cmd := NewCommand("git", "describe", "--abbrev=0", "--tags")
	cmd.Stdout = &b
	if cmd.Run() != nil {
		return "v0.0.0"
	}
	return strings.TrimSpace(b.String())
}

func BashPipe(in string) (result []string) {
	for _, part := range strings.Split(in, "\n") {
		if len(strings.TrimSpace(part)) > 0 {
			result = append(result, part)
		}
	}

	return
}
