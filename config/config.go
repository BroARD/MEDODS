package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	IP string `envconfig:"APP_IP" default:"0.0.0.0"`
	Port string `envconfig:"APP_PORT" default:"8080"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
    once.Do(func() {
        err := godotenv.Load(".env")
        if err != nil {
            log.Println("Warning: .env file not found or failed to load, using system env variables")
        }

        instance = &Config{}
        if err := cleanenv.ReadEnv(instance); err != nil {
            help, _ := cleanenv.GetDescription(instance, nil)
            log.Fatalf("Error reading env: %v\n%s", err, help)
        }
    })
    return instance
}