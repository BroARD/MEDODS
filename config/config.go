package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug"`
	Listen  struct {
		Type   string `yaml:"type"`
		BindIP string `yaml:"bind_ip"`
		Port   string `yaml:"port"`
	} `yaml:"listen"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
		if err := cleanenv.ReadConfig("config/config.yml", instance); err != nil{
			help, _ := cleanenv.GetDescription(instance, nil)
			log.Fatal(help)
		}
	})
	return instance
}
