package config

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Postgresql struct {
		Host     string `yaml:"host" env-default:"localhost"`
		Port     string `yaml:"port" env-default:"5432"`
		User     string `yaml:"user" env-default:"postgres"`
		Password string `yaml:"password"`
		DbName   string `yaml:"db_name"`
	} `yaml:"postgresql"`
	Proxy struct {
		User         string `yaml:"user"`
		Password     string `yaml:"password"`
		ProxyAddress string `yaml:"proxy_address"`
		Port         string `yaml:"port"`
	} `yaml:"proxy"`
	Tokens struct {
		TelegramToken string `yaml:"telegram_token"`
		OpenAIToken   string `yaml:"openai_token"`
	} `yaml:"tokens"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		log.Println("Reading configuration...")

		instance = &Config{}

		wd, _ := os.Getwd()
		pathToYML := filepath.Dir(wd) + "\\config.yml"

		if err := cleanenv.ReadConfig(pathToYML, instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			log.Println(help)
			log.Fatalf("error reading configuration, error: %v\n", err)
		}

		log.Println("Configuration was successfully read")
	})
	return instance
}
