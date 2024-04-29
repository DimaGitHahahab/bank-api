package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

type Config struct {
	HttpPort      string `envconfig:"HTTP_PORT" default:"8080"`
	DbUrl         string `envconfig:"DB_URL" required:"true"`
	MigrationPath string `envconfig:"MIGRATION_PATH" required:"true"`
	JwtSecret     string `envconfig:"JWT_SECRET" required:"true"`
}

func LoadConfig(log *zap.SugaredLogger) *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("Failed to load env variables: ", err)
	}

	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalln("Failed to process env variables: ", err)
	}

	log.Infof("Loaded env variables successfully")
	return &cfg
}
