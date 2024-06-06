package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
	Kafka struct {
		Host      string `yaml:"host"`
		Port      string `yaml:"port"`
		Topic     string `yaml:"topic"`
		Partition string `yaml:"partition"`
	} `yaml:"kafka"`
}

func Parse(env string) *Config {
	var cfg Config

	path := fmt.Sprintf("configs/%s.yml", env)

	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	return &cfg
}
