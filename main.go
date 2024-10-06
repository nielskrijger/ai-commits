package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/golang-cz/devslog"
	"github.com/urfave/cli/v2"
)

var (
	staged = &cli.BoolFlag{
		Name:  "staged",
		Value: true,
		Usage: "generate a message from staged changes",
	}
	path = &cli.StringFlag{
		Name:  "dir",
		Value: ".",
		Usage: "the directory path to the git repository",
	}
	choices = &cli.IntFlag{
		Name:  "choices",
		Value: 8,
		Usage: "how many choices to generate",
	}
	model = &cli.StringFlag{
		Name:  "model",
		Value: "gpt-4o-mini",
		Usage: "which ChatGPT model to use (not implemented yet)",
	}
	hint = &cli.StringFlag{
		Name:  "hint",
		Usage: "provide a hint to the AI model to improve the quality",
	}
	debug = &cli.BoolFlag{
		Name:  "debug",
		Value: false,
		Usage: "when true will print debug information",
	}
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{staged, path, choices, model, debug, hint},
		Action: func(c *cli.Context) error {
			// Default to the current directory if no path is provided
			projectPath := c.Args().Get(0)
			if projectPath == "" {
				projectPath = "."
			}

			// Get the API key from the config file
			configReader := &FileConfigReader{
				Filename: "config.yml",
			}
			apiKey, err := configReader.ReadAPIKey()
			if err != nil {
				log.Fatal(err)
			}

			// Create a custom logger
			slogOpts := &slog.HandlerOptions{}
			logOpts := devslog.Options{HandlerOptions: slogOpts}
			if debug.Get(c) {
				slogOpts.Level = slog.LevelDebug
			}
			logger := slog.New(devslog.NewHandler(os.Stdout, &logOpts))

			// Replace the default logger with the custom logger
			slog.SetDefault(logger)

			// Generate the commit message option
			opts := GenerateMessageOptions{
				APIKey:  apiKey,
				Dir:     path.Get(c),
				Staged:  staged.Get(c),
				Choices: choices.Get(c),
				Model:   model.Get(c),
				Hint:    hint.Get(c),
			}
			err = GenerateMessage(opts)
			if err != nil {
				slog.With("options", opts).Error(err.Error())
				os.Exit(1)
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		slog.Error(err.Error())
	}
}
