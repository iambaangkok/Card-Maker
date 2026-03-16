package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

const configPath = "app.yaml"

type Config struct {
	HTML     HTMLTemplateConfig `yaml:"html"`
	Renderer RendererConfig    `yaml:"renderer"`
}

type HTMLTemplateConfig struct {
	OutputHTMLEnabled bool `yaml:"output-html-enabled"`
}

type RendererConfig struct {
	URL       string `yaml:"url"`
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