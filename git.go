// Based on https://github.com/tskoyo/ai-git-commit

package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

var (
	execCommand         = exec.Command
	gitDiffFilePatterns = []string{":!go.sum"}
)

// gitDiff returns the diff of the git repository at the given directory path. If
// staged is true, it returns the staged changes. Otherwise, it returns the
// staged+unstaged changes. Files that are untracked are not included in either
// case.
func gitDiff(staged bool, dirPath string) (string, error) {
	err := isGitRepo(dirPath)
	if err != nil {
		return "", err
	}

	diff, err := getDiff(staged)
	if err != nil {
		return "", err
	}

	formattedDiff := formatDiff(diff)
	if formattedDiff == "" {
		return "", errNoChanges
	}

	return formattedDiff, nil
}

func isGitRepo(path string) error {
	cmd := execCommand("git", "rev-parse", "--is-inside-work-tree", "--git-dir", path)
	if err := cmd.Run(); err != nil {
		if path == "." {
			return errors.New("current directory is not a git repository")
		} else {
			return fmt.Errorf("%q is not a git repository", path)
		}
	}
	return nil
}

func getDiff(staged bool) (string, error) {
	args := []string{"diff"}
	if staged {
		args = append(args, "--staged")
	}
	args = append(args, "--")
	args = append(args, gitDiffFilePatterns...)

	cmd := execCommand("git", args...)

	// Run the command and capture output and error
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func formatDiff(diff string) string {
	lines := strings.Split(diff, "\n")

	var formattedLines []string

	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "diff") || strings.HasPrefix(line, "index") {
			continue
		}
		formattedLines = append(formattedLines, strings.TrimSpace(line))
	}

	return strings.Join(formattedLines, "\n")
}
