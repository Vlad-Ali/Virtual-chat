package app

import (
	"fmt"
	"os"

	"github.com/Vlad-Ali/Virtual-chat/internal/config/cors"
	"github.com/Vlad-Ali/Virtual-chat/internal/config/http"
	"github.com/Vlad-Ali/Virtual-chat/internal/config/ollama"
	"github.com/Vlad-Ali/Virtual-chat/internal/config/websocket"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	HTTPConfig      httpconfig.HTTPConfig           `yaml:"http"`
	WebSocketConfig websocketconfig.WebSocketConfig `yaml:"ws"`
	OllamaConfig    ollama.OllamaConfig             `yaml:"ollama"`
	CorsConfig      cors.CorsConfig
}

func LoadConfig(configPath string) (*Config, error) {
	_ = godotenv.Load()
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config: %v", err)
	}

	var config Config
	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config %v", err)
	}

	config.HTTPConfig.Address = os.Getenv("HTTP_ADDRESS")
	config.OllamaConfig.URL = os.Getenv("OLLAMA_URL")
	config.OllamaConfig.ModelName = os.Getenv("OLLAMA_MODEL_NAME")
	config.CorsConfig.AllowedOrigins = []string{os.Getenv("FRONTEND_URL")}

	return &config, nil
}
