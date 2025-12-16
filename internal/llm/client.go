package llm

import (
	"context"
	"fmt"
	"mime/multipart"

	"voiceline-audio-backend/internal/common"

	openai "github.com/sashabaranov/go-openai"
)

type Client struct {
	client *openai.Client
}

func NewClient(apiKey string) *Client {
	return &Client{client: openai.NewClient(apiKey)}
}

func (c *Client) TranscribeAudio(ctx context.Context, file multipart.File, filename string) (string, error) {
	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: filename,
		Reader:   file,
	}

	resp, err := c.client.CreateTranscription(ctx, req)
	if err != nil {
		return "", common.NewInternalServerError("transcription failed", err)
	}

	if resp.Text == "" {
		return "", common.NewInternalServerError("got empty transcription", nil)
	}

	return resp.Text, nil
}

func (c *Client) ProcessTranscription(ctx context.Context, transcription string) (*Result, error) {
	prompt := fmt.Sprintf(`Give me a summary and action items from this.
Return JSON: {"summary": "...", "action_items": [...]}

%s`, transcription)

	req := openai.ChatCompletionRequest{
		Model: openai.GPT4o,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "Extract summary and action items as JSON.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.3,
		MaxTokens:   1000,
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		},
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, common.NewInternalServerError("gpt call failed", err)
	}

	if len(resp.Choices) == 0 {
		return nil, common.NewInternalServerError("no response from gpt", nil)
	}

	content := resp.Choices[0].Message.Content

	result, err := ParseGPTResponse(content)
	if err != nil {
		return nil, common.NewInternalServerError("couldn't parse response", err)
	}

	return result, nil
}
