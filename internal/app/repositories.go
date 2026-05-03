package app

import (
	"github.com/Vlad-Ali/Virtual-chat/internal/domain/session"
	"github.com/Vlad-Ali/Virtual-chat/internal/infrastructure/ollama"
	"github.com/ollama/ollama/api"
)

type Repositories struct {
	Provider session.Provider
}

func NewRepositories(cfg *Config, client *api.Client) *Repositories {
	provider := ollama.NewClient(client, &cfg.OllamaConfig)
	return &Repositories{Provider: provider}
}
