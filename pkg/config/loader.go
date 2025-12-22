package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func LoadFromFile(filename string) (*ModelRegistry, error){
	data, err := os.ReadFile(filename)
	if err != nil{
		return nil, fmt.Errorf("读取配置文件 '%s' 失败：%w", filename, err)
	}

	var configs []ModelConfig

	if err := json.Unmarshal(data, &configs); err != nil{
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}

	registry := NewModelRegistry()
	for _, config := range configs{
		registry.Registry(config)
	}

	return registry, err
}

func SaveToFile(registry *ModelRegistry, filename string) error{
	var configs []ModelConfig
	for _, name := range registry.List(){
		config, _ := registry.Get(name)
		configs = append(configs, config)
	}

	data, err := json.MarshalIndent(configs, "", " ")
	if err != nil{
		return fmt.Errorf("序列化配置失败：%W", err)
	}
	
	if err := os.WriteFile(filename, data, 0644); err != nil{
		return fmt.Errorf("写入文件 '%S' 失败：%w", filename, err)
	}

	return nil
}

func LoadOrCreate(filename string) *ModelRegistry{
	registry, err := LoadFromFile(filename)
	if err != nil {
		return NewModelRegistry()
	}

	return registry
}