package config

// ModelConfig 
// use example：
//
//    config := config.ModelConfig{
//        Name:    "my-claude",
//        BaseURL: "https://openrouter.ai/api/v1",
//        APIKey:  "sk-or-v1-xxx",
//        Model:   "anthropic/claude-3.5-sonnet",
//    }

type ModelConfig struct{
	Name string `json:"name"`

	BaseURL string 	`json:"baseURL"`

	APIKey string `json:"apiKey"`

	Model string `json: "model"`
}