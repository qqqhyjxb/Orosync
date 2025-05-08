package config

import (
	"Orosync/internal/raft"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 结构体包含程序的配置选项
type Config struct {
	IsTestMode bool           `yaml:"is_test_mode"` // 是否是测试模式
	Env        string         `yaml:"env"`          // 环境设置，如开发、生产等
	Port       int            `yaml:"port"`         // 服务器端口
	Logs       *raft.LogsInfo `yaml:"logs"`         // 日志信息
}

// InitConfig 从 YAML 配置文件加载配置
func InitConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	config := &Config{}
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return config, nil
}
