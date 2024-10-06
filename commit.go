package main

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/fatih/color"
)

var errNoChanges = errors.New("no changes to commit")

type GenerateMessageOptions struct {
	APIKey  string `json:"-"` // Just in case exclude for security reasons
	Model   string
	Staged  bool
	Choices int
	Dir     string
	Hint    string
}

func GenerateMessage(opts GenerateMessageOptions) error {
	diff, err := gitDiff(opts.Staged, opts.Dir)

	// If no changes are staged, check unstaged changes instead. Most likely the user
	// just forgot to stage anything.
	if err != nil {
		if errors.Is(err, errNoChanges) && opts.Staged {
			diff, err = gitDiff(false, opts.Dir)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// Generate prompt
	msg := generatePrompt(diff, opts.Hint, opts.Choices)
	slog.With("options", opts).Debug(msg)

	// Generate commit message(s)
	ai := NewOpenAIGenerator(opts.APIKey, opts.Model)
	answer, err := ai.GenerateCommitMsg(msg)
	if err != nil {
		return err
	}

	// Print commit messages
	var printColor colorFn
	for i, response := range answer.Messages {
		printColor = determineColor(answer, response.DiffCoverage)
		printColor(fmt.Sprintf("\n%[1]s Message %[2]d (%.0[3]f%%) %[1]s\n\n", seperator, i+1, response.DiffCoverage*100))
		fmt.Printf("%s\n", response)
	}

	// Final commit message
	final := answer.FinalCommitMessage
	printColor = determineColor(answer, final.DiffCoverage)
	printColor(fmt.Sprintf("\n%[1]s FINAL (%.0[2]f%%) %[1]s\n\n", seperator, final.DiffCoverage*100))
	fmt.Printf("%s\n", final)

	return nil
}

type colorFn func(format string, a ...interface{})

func determineColor(answer Answer, confidence float64) colorFn {
	if confidence == answer.HighestConfidence() {
		return color.Green
	}
	if confidence == answer.LowestConfidence() {
		return color.Red
	}
	return color.Cyan
}
