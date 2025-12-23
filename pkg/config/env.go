package config

import "os"

const (
	EnvAPIKey  = "OPENAI_API_KEY"
	EnvBaseURL = "OPENAI_BASE_URL"
	EnvModel   = "OPENAI_MODEL"
)

func DefaultConfig() ModelConfig{
	return ModelConfig{
		Name: "default",
		APIKey: os.Getenv(EnvAPIKey),
		BaseURL: getEnvOrDefault(EnvBaseURL, "https://api.openai.com/v1"),
        Model:   getEnvOrDefault(EnvModel, "gpt-4o-mini"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
			return v
	}
	return defaultValue
}

func LoadWithEnv(filename string) (*ModelRegistry, error) {
	registry, err := LoadFromFile(filename)
	if err != nil {
			return nil, err
	}

	apiKey := os.Getenv(EnvAPIKey)
	for _, name := range registry.List() {
			cfg, _ := registry.Get(name)
			if cfg.APIKey == "" {
					cfg.APIKey = apiKey
					registry.Registry(cfg)
			}
	}

	return registry, nil
}