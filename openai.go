package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

type CommitMessage struct {
	Type         string  `json:"type"`
	Subject      string  `json:"subject"`
	Description  string  `json:"description"`
	DiffCoverage float64 `json:"diff_coverage"`
}

func (c CommitMessage) String() string {
	return fmt.Sprintf("%s: %s\n\n%s", c.Type, c.Subject, c.Description)
}

type FinalCommitMessage struct {
	Type         string   `json:"type"`
	Subject      string   `json:"subject"`
	DiffCoverage float64  `json:"diff_coverage"`
	Changes      []string `json:"changes"`
}

func (c FinalCommitMessage) String() string {
	msg := fmt.Sprintf("%s: %s\n\n", c.Type, c.Subject)
	for _, change := range c.Changes {
		// If change doesn't have the "- " suffix already add it
		if change[:2] != "- " {
			change = "- " + change
		}
		msg = msg + change + "\n"
	}
	return msg
}

type Answer struct {
	Messages           []CommitMessage    `json:"messages"`
	FinalCommitMessage FinalCommitMessage `json:"final_message"`
}

var result Answer

func (r Answer) HighestConfidence() float64 {
	highest := r.FinalCommitMessage.DiffCoverage
	for _, msg := range r.Messages {
		if msg.DiffCoverage > highest {
			highest = msg.DiffCoverage
		}
	}
	return highest
}

func (r Answer) LowestConfidence() float64 {
	lowest := r.FinalCommitMessage.DiffCoverage
	for _, msg := range r.Messages {
		if msg.DiffCoverage < lowest {
			lowest = msg.DiffCoverage
		}
	}
	return lowest
}

type OpenAIGenerator struct {
	apiKey string
	model  string
}

func NewOpenAIGenerator(apiKey, model string) *OpenAIGenerator {
	return &OpenAIGenerator{
		apiKey: apiKey,
		model:  model,
	}
}

func (o OpenAIGenerator) GenerateCommitMsg(prompt string) (Answer, error) {
	if prompt == "" {
		return Answer{}, errors.New("empty prompt")
	}

	// Generate JSON schema for the result
	schema, err := jsonschema.GenerateSchemaForType(result)
	if err != nil {
		return Answer{}, fmt.Errorf("generate JSON schema: %s", err)
	}

	// Talk to OpenAI
	client := openai.NewClient(o.apiKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: o.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},

			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
				JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
					Name:   "commits",
					Schema: schema,
					Strict: true,
				},
			},
		},
	)

	// Unmarshal the response
	err = schema.Unmarshal(resp.Choices[0].Message.Content, &result)
	if err != nil {
		return Answer{}, fmt.Errorf("unmarshal ChatGPT JSON response: %s", err)
	}

	return result, err
}
