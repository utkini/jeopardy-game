package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HttpServer  `yaml:"http_server"`
}

type HttpServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type GameConfig struct {
	Categories []struct {
		Name      string `yaml:"name"`
		Questions []struct {
			Question  string `yaml:"question"`
			Answer    string `yaml:"answer"`
			Points    int    `yaml:"points"`
			MediaType string `yaml:"media_type"`
			MediaURL  string `yaml:"media_url"`
		} `yaml:"questions"`
	} `yaml:"categories"`
}

type GameConfigManager struct {
	Config GameConfig
}

func NewGameConfigManager(configPath string) (*GameConfigManager, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config GameConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &GameConfigManager{
		Config: config,
	}, nil
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("config_path not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return &cfg
}
