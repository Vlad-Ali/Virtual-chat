package ollama

type OllamaConfig struct {
	URL            string
	ModelName      string
	SystemPrompt   string `yaml:"system_prompt"`
	CallbackPrompt string `yaml:"callback_prompt"`
}
