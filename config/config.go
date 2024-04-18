package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type ConfigType struct {
	App struct {
		Name string `yaml:"name"`
		Port string `yaml:"port"`
	} `yaml:"app"`
	Redis struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Db       string `yaml:"db"`
	} `yaml:"redis"`
	Event struct {
		DataDir          string `yaml:"data_dir"`
		Encode           string `yaml:"encode"`
		DataFileMaxLines int32  `yaml:"data_file_max_lines"`
	} `yaml:"event"`
}

var Data ConfigType

// InitConfig 初始化配置
func InitConfig() {
	filePath := "config/config.yaml"
	// 读取YAML文件内容
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("init config error %s: %v", filePath, err)
	}
	// 解析YAML内容到Config结构体
	err = yaml.Unmarshal(data, &Data)
	if err != nil {
		log.Fatalf("init config Error parsing YAML: %v", err)
	}
	log.Println("init config success")
}
