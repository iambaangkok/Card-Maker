package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

const configPath = "app.yaml"

type Config struct {
	CSV	CSVConfig `yaml:"csv"`
}
type CSVConfig struct {
	Entities []EntityConfig `yaml:"entities"`
	Dir string `yaml:"dir"`
}

type EntityConfig struct {
	Name string `yaml:"name"`
}


func (cfg *Config) Init() {
	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		log.Fatal("failed to load config")
	}
}

func LoadConfig() Config {
	cfg := Config{}
	cfg.Init()
	return cfg
}