package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type DatabaseConfig struct {
	Server   string `yaml:"server"`
	Database string `yaml:"database"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	SSL      string `yaml:"ssl"`
}

type TelegramConfig struct {
	Token           string `yaml:"token"`
	OwnerTelegramID int64  `yaml:"owner_telegram_id"`
}

type AIConfig struct {
	Key                 string `yaml:"key"`
	BaseURL             string `yaml:"base_url"`
	MaxCompletionTokens int64  `yaml:"max_completion_tokens"`
	Model               string `yaml:"model"`
	Prompt              string `yaml:"prompt"`
}

type Config struct {
	Env      string          `yaml:"env"`
	Database *DatabaseConfig `yaml:"database"`
	Telegram *TelegramConfig `yaml:"telegram"`
	AI       *AIConfig       `yaml:"ai"`
}

func MustLoad() *Config {
	path := MustGetPath()

	return MustLoadByPath(path)
}

func MustGetPath() string {
	path := getPath()
	if path == "" {
		log.Fatal("config path not set")
	}

	return path
}

func getPath() string {
	if path := getPathByEnv(); path != "" {
		return path
	}
	return getPathByFlag()
}

func getPathByEnv() string {
	path := os.Getenv("CONFIG_PATH")
	return path
}

func getPathByFlag() string {
	var path string

	flag.StringVar(&path, "config_path", "", "path to config file")
	flag.Parse()
	return path
}

func MustLoadByPath(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatal("failed to read config:" + err.Error())
	}

	return &cfg
}
