package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"

)

type HTTPServer struct {
	Addr string `yaml:"address"`
}

type Config struct {
	Env  string `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
	StoragePath string `yaml:"storage_path" env-required:"true" env:"DB_PATH" env-default:"data/database.sqlite"`
	HTTPServer HTTPServer `yaml:"http_server"`
}

func MustLoad() *Config{
	// Load config from file
	var configPath string
	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "path to config file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("config path is required")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file not found: %s", configPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("failed to load config: %v", err.Error())
	}

	return &cfg
}