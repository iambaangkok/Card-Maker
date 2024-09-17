package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

const configPath = "app.yaml"

type Config struct {
	CSV	CSVConfig `yaml:"csv"`
	Renderer RendererConfig `yaml:"renderer"`
}
type CSVConfig struct {
	Entities []EntityConfig `yaml:"entities"`
	Dir string `yaml:"dir"`
}

type EntityConfig struct {
	Name string `yaml:"name"`
}

type RendererConfig struct {
	URL string `yaml:"url"`
	OutputDir string `yaml:"output-dir"`
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