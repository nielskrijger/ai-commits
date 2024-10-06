package main

import (
	"fmt"
	"strings"
)

var (
	bannedWords = []string{
		"better",
		"clarity",
		"clean",
		"cleaned",
		"cleaner",
		"cleaning",
		"cohesive",
		"combine",
		"comprehensive",
		"consistent",
		"consolidate",
		"describe",
		"developed",
		"elegantly",
		"enhance",
		"enhance",
		"essential",
		"establish",
		"flexible",
		"improve",
		"improve",
		"integrate",
		"introduce",
		"maintainability",
		"optimize",
		"populate",
		"readability",
		"refactored",
		"refine",
		"reformat",
		"reorganize",
		"reusable",
		"robust",
		"simplify",
		"streamline",
		"streamline",
		"tidy",
		"tidying",
		"unified",
	}
)

const (
	seperator = "----"

	prompt = `
Generate [CHOICES] commits messages based on the git diff changes below.

[HINT]

Each commit contains the following fields:
- "type": a string representing the type of commit using conventional commits.
- "diff_coverage": a number between 0 and 1 representing the 0-100% coverage of the git diff.
  - This is NOT how confident is about the commit message.
  - It is OK to have a low coverage if the commit message only covers a small part of the diff.
  - If everything is 100% you did something wrong.
- "subject": a short (60-80 characters) subject in the imperative mood.
  - Start with a lowercase letter.
  - Do not include the conventional commit type in the subject.
  - Imagine saying "If applied, this commit will [SUBJECT]" where [SUBJECT] is the subject.
- "description": an optional short explanation of max 250 characters.
  - Remove the explanation if it does not improve the "coverage" score.

Finally, generate a final commit based on the best commit messages.
- "changes": order by the highest coverage change to the lowest.
  - Exclude the "why" or "for" explanation if most developers would know the reason.

You MUST follow these guidelines;
- Never list any filename like "file.go".
- Never use the following words: [BANNED_WORDS].
- Ignore any diff lines that contain error handling changes.
- Ignore any diff lines that contain spelling changes or tiny refactors.

Everything after below is the git diff:

[GIT_DIFF]
`
)

// generatePrompt creates the full prompt message for ChatGPT.
func generatePrompt(diff, hint string, choices int) string {
	msg := strings.ReplaceAll(prompt, "[CHOICES]", fmt.Sprintf("%d", choices))
	msg = strings.ReplaceAll(msg, "[GIT_DIFF]", diff)
	msg = strings.ReplaceAll(msg, "[HINT]", fmt.Sprintf("Use the following HINT to generate the messages: %s", hint))
	msg = strings.ReplaceAll(msg, "[BANNED_WORDS]", strings.Join(bannedWords, ", "))

	return msg
}
