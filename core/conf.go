package core

import (
	"gopkg.in/yaml.v3"
	"log"
	"server/config"
	"server/utils"
)

// InitConf 初始化 config 配置
func InitConf() *config.Config {
	c := &config.Config{}
	yamlConfig, err := utils.LoadYAML()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v\n", err)
		return nil
	}
	if err := yaml.Unmarshal(yamlConfig, c); err != nil {
		log.Fatalf("Failed to unmarshal configuration: %v\n", err)
		return nil
	}
	return c
}
