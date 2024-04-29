package main

import (
	"bank-api/internal/app"
	"bank-api/pkg/config"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

func main() {
	rawLogger, _ := zap.NewProduction()
	log := rawLogger.Sugar()
	defer log.Sync()

	cfg := config.LoadConfig(log)

	a := app.New(log, cfg)

	a.Run()
}
