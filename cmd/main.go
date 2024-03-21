package main

import (
	"bank-api/internal/api/server"
	"bank-api/internal/repository"
	"bank-api/internal/service"
	"context"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	pgxzap "github.com/jackc/pgx-zap"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	time.Sleep(3 * time.Second)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	rawLogger, _ := zap.NewProduction()
	logger := rawLogger.Sugar()
	defer logger.Sync()

	if err := godotenv.Load(); err != nil {
		logger.Fatalf("failed to load env variables: %v", err)
	}

	url := os.Getenv("DB_URL")
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		logger.Fatalf("can't parse pgxpool config: %v", err)
	}

	config.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   pgxzap.NewLogger(logger.Desugar()),
		LogLevel: tracelog.LogLevelDebug,
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		logger.Fatalf("can't create new pgxpool: %v", err)
	}
	defer pool.Close()

	processMigration(os.Getenv("MIGRATION_PATH"), url, logger)

	repo := repository.New(pool, logger)

	userService := service.NewUserService(repo)
	accountService := service.NewAccountService(repo)
	transactionService := service.NewTransactionService(repo)

	port := os.Getenv("HTTP_PORT")
	srv := server.New(userService, accountService, transactionService)

	go func() {
		if err := srv.Run(port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("server run failed: %v", err)
		}
	}()

	<-ctx.Done()

	stop()
	logger.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("server forced to shutdown: %v", err)
	}

	logger.Info("server is shut down")
}

func processMigration(migrationURL string, dbSource string, logger *zap.SugaredLogger) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		logger.Fatalf("failed to create migration: %v", err)
	}

	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Fatalf("failed to migrate: %v", err)

	}
	defer migration.Close()

	logger.Info("migration successful")
}
