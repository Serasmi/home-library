package config

import (
	"github.com/Serasmi/home-library/pkg/logging"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	App struct {
		Host     string `env:"APP_HOST" envDefault:"0.0.0.0"`
		Port     string `env:"APP_PORT,required"`
		LogLevel string `env:"APP_LOG_LEVEL" envDefault:"info"`
	}
	DB struct {
		Host     string `env:"MONGODB_HOST,required"`
		Port     string `env:"MONGODB_PORT" envDefault:"27017"`
		Name     string `env:"MONGODB_DATABASE,required"`
		Username string `env:"MONGODB_USERNAME,required"`
		Password string `env:"MONGODB_PASSWORD,required"`
	}
	JWT struct {
		Secret string `env:"JWT_SECRET,required"`
	}
}

var config Config

func GetConfig() *Config {
	return &config
}

func InitConfig(logger *logging.Logger) *Config {
	logger.Info("init config")

	err := godotenv.Load()
	if err != nil {
		logger.Fatalf("Error loading .env file: %s", err.Error())
	}

	err = env.Parse(&config)
	if err != nil {
		logger.Fatal("init config error ", err)
	}

	return &config
}
