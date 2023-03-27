package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type ApiConfig struct {
	ApiKeyToTokens map[string]int `mapstructure:"api_keys"`
}

func getConfig() ApiConfig {
	configPath := "."
	if configPath, ok := os.LookupEnv("CONF_PATH"); ok {
		fmt.Println("CONF_PATH: " + configPath)
	}
	// 设置viper的配置文件名和路径
	viper.SetConfigName("config")   // 文件名(不带后缀)
	viper.AddConfigPath(configPath) // 相对路径

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// 将映射列表映射到结构体中
	var config ApiConfig
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("Error unmarshaling config file: %s \n", err))
	}

	fmt.Println("config: ", config)

	return config
}
