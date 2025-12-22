package config
import (
	"fmt"

	"nvgo/pkg/agent"	
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type ModelRegistry struct{
	configs map[string]ModelConfig
}

func NewModelRegistry() *ModelRegistry{
	return &ModelRegistry{
			configs: make(map[string]ModelConfig),
	}
}

func(r *ModelRegistry) Registry(config ModelConfig){
	r.configs[config.Name] = config
}

func(r *ModelRegistry) Get(name string) (ModelConfig, bool){
	config, ok := r.configs[name]
	return config, ok
}

func (r *ModelRegistry) List() []string {
	names := make([]string, 0, len(r.configs))
	for name := range r.configs{
		names = append(names, name)
	}
	return names
}

func (r *ModelRegistry) CreateAgent(
	configName string,
	agentName string,
	instruction string,
)(*agent.Agent, error){
	config, ok := r.Get(configName)
	if !ok{
		return nil, fmt.Errorf("配置 '%s' 不存在，请先使用 Register() 注册", configName)
	}

	client := openai.NewClient(
		option.WithAPIKey(config.APIKey),
		option.WithBaseURL(config.BaseURL),
	)

	baseAgent := agent.New(configName).WithInstructions(config.Model).WithModel(config.Model).WithClient(client)

	return baseAgent, nil
}

func (r *ModelRegistry) CreateAgentWithOptions(configName string, setupFunc func(*agent.Agent) *agent.Agent,
)( *agent.Agent, error){
	config, ok := r.Get(configName)
	if !ok {
		return nil, fmt.Errorf("配置 '%s' 不存在", configName)
	}

	client := openai.NewClient(
		option.WithAPIKey(config.APIKey),
		option.WithBaseURL(config.BaseURL),
	)

	baseAgent := agent.New(config.Name).WithModel(config.Model).WithClient(client)

	if setupFunc != nil{
		baseAgent = setupFunc(baseAgent)
	}

	return baseAgent, nil
}

func (r *ModelRegistry) Has(name string) bool{
	_, ok := r.configs[name]
	return ok
}

func(r *ModelRegistry) Delete(name string){
	delete(r.configs, name)
}

func(r *ModelRegistry) Count() int {
	return len(r.configs)
}