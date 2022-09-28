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
		tag := NewCommandIn(dir, "git", "tag", fmt.Sprintf("v%s", nextVersion.String()), "-a")
		tag.Stdin = os.Stdin
		Check(tag.Run())
	} else {
		Check(NewCommandIn(dir, "git", "tag", fmt.Sprintf("v%s", nextVersion.String())).Run())
	}

	if !dry {
		Confirm(fmt.Sprintf("Pressing [y] will push this (%s) release to GitHub. Are you sure?", dir))
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

func Confirm(message string, args ...interface{}) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("%s [y/n] ", fmt.Sprintf(message, args...))
		answer, err := reader.ReadString('\n')
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
