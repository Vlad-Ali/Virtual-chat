package ollama

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Vlad-Ali/Virtual-chat/internal/config/ollama"
	"github.com/Vlad-Ali/Virtual-chat/internal/domain/message"
	"github.com/ollama/ollama/api"
)

var (
	SystemRole    string = "system"
	UserRole      string = "user"
	AssistantRole string = "assistant"
)

type Client struct {
	client *api.Client
	config *ollama.OllamaConfig
	logger *slog.Logger
}

func NewClient(client *api.Client, config *ollama.OllamaConfig) *Client {
	return &Client{client: client, config: config, logger: slog.Default().With("module", "ollama_client")}
}

func (c *Client) GetAnswer(ctx context.Context, history []*message.Message) (string, error) {
	ollamaMessages := c.generateMessages(history)
	request := &api.ChatRequest{
		Model:    c.config.ModelName,
		Messages: ollamaMessages,
		Stream:   boolPtr(false),
		Options: map[string]interface{}{
			"num_predict":      80,
			"temperature":      0.75,
			"top_p":            0.9,
			"repeat_penalty":   1.1,
			"presence_penalty": 0.2,
		},
	}

	var response string
	c.logger.Debug("sending prompt to ollama", "prompt", ollamaMessages)
	err := c.client.Chat(ctx, request, func(resp api.ChatResponse) error {
		response = resp.Message.Content
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("error sending prompt to ollama: %w", err)
	}

	c.logger.Debug("received prompt response", "response", response)

	return response, nil
}

func (c *Client) generateMessages(history []*message.Message) []api.Message {
	ollamaMessages := make([]api.Message, 0)

	ollamaMessages = append(ollamaMessages, api.Message{
		Role:    SystemRole,
		Content: c.config.SystemPrompt,
	})

	for _, msg := range history {
		if msg.Type == message.CallBackType {
			continue
		}

		role := UserRole
		if msg.Type == message.ModelType {
			role = AssistantRole
		}

		ollamaMessages = append(ollamaMessages, api.Message{
			Role:    role,
			Content: msg.Msg,
		})
	}

	if len(history) > 0 && history[len(history)-1].Type == message.CallBackType {
		ollamaMessages = append(ollamaMessages, api.Message{
			Role:    UserRole,
			Content: c.config.CallbackPrompt,
		})
	}

	return ollamaMessages
}

func boolPtr(b bool) *bool {
	return &b
}
