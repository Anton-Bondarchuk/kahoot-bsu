package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)


type Config struct {
	Env            string     `yaml:"env" env-default:"prod"`
	StorageConfig  StorageConfig     `yaml:"storage" env-required:"true"`
	BotConfig      BotConfig  `yaml:"bot" env-required:"true"`
}

type StorageConfig struct {
	DatabaseUrl string  `yaml:"database_url" env-required:"true"`
}

type BotConfig struct {
	Token string `yaml:"token" env-required:"true"`
	Timeout int  `yaml:"timeout" env-default:"60"`
	Debug    bool `yaml:"debug" env-default:"false"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	// TODO: check the current path
	fmt.Println("config path: ", res)
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	if res == "" {
		res = "./config/prod.yml"
	}

	return res
}


// TODO: integrate to logic

type ServerConfig struct {
	Host            string
	Port            int
	ReadTimeout     int
	WriteTimeout    int
	ShutdownTimeout int
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// RabbitMQConfig holds RabbitMQ configuration
type RabbitMQConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	VHost    string
	Exchange string
}
