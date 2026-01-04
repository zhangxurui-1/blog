package utils

import (
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"server/global"
)

// 读取或保存 yaml 文件

// 因为程序入口是 main.go，所以这里直接写 server 文件夹下的相对路径即可
const configFile = "config.yaml"

// LoadYAML 从 yaml 文件读取
func LoadYAML() ([]byte, error) {
	return os.ReadFile(configFile)
}

// SaveYAML 将配置保存至 yaml 文件
func SaveYAML() error {
	bytes, err := yaml.Marshal(global.Config)
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, bytes, fs.ModePerm)
}
