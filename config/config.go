// @File: config
// @Author: Nanjia Ding
// @Date: 2024/06/21
package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Redis struct {
		Addr     string `json:"addr"`
		Password string `json:"password"`
		DB       int    `json:"db"`
	} `json:"redis"`
}

func LoadConfig(file string) (*Config, error) {
	config := &Config{}
	configFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
