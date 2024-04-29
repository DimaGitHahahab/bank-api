package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"bank-api/internal/handlers"
	"bank-api/internal/repository"
	"bank-api/internal/router"
	"bank-api/internal/server"
	"bank-api/internal/service"
	"bank-api/pkg/config"
	"bank-api/pkg/signal"

	"github.com/golang-migrate/migrate/v4"
	pgxzap "github.com/jackc/pgx-zap"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"go.uber.org/zap"
)

type App struct {
	config  *config.Config
	sigQuit chan os.Signal
	ctx     context.Context
	server  *server.Server
	log     *zap.SugaredLogger
}

func New(log *zap.SugaredLogger, cfg *config.Config) *App {
	sigQuit := signal.GetShutdownChannel()

	ctx := context.Background()

	userRepo, accountRepo := setupRepo(ctx, log, cfg)

	processMigration(cfg.MigrationPath, cfg.DbUrl, log)

	userService := service.NewUserService(userRepo)
	accountService := service.NewAccountService(accountRepo)
	transactionService := service.NewTransactionService(accountRepo)

	h := handlers.NewHandler(cfg.JwtSecret, userService, accountService, transactionService)

	srv := server.New(router.NewRouter(h))

	return &App{
		config:  cfg,
		sigQuit: sigQuit,
		ctx:     ctx,
		server:  srv,
		log:     log,
	}
}

func (a *App) Run() {
	go func() {
		a.log.Infoln("Starting server on port ", a.config.HttpPort)
		if err := a.server.Run(a.config.HttpPort); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Fatalln("Failed to start server: ", err)
		}
	}()

	<-a.sigQuit
	a.log.Infoln("Gracefully shutting down server")

	ctx, cancel := context.WithTimeout(a.ctx, 2*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		log.Fatalln("Failed to shutdown the server gracefully: ", err)
	}

	a.log.Infoln("Server shutdown is successful")
}

func setupRepo(ctx context.Context, log *zap.SugaredLogger, cfg *config.Config) (repository.UserRepository, repository.AccountRepository) {
	pool, err := setupPgxPool(ctx, log, cfg)
	if err != nil {
		log.Fatalln(err)
	}
	return repository.New(pool, log)
}

func setupPgxPool(ctx context.Context, log *zap.SugaredLogger, cfg *config.Config) (*pgxpool.Pool, error) {
	pgxConfig, err := pgxpool.ParseConfig(cfg.DbUrl)
	if err != nil {
		log.Errorln("Failed to parse pgxpool pgxConfig: ", err)
		return nil, err
	}

	pgxConfig.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   pgxzap.NewLogger(log.Desugar()),
		LogLevel: tracelog.LogLevelDebug,
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		log.Errorln("Failed to create new pool: ", err)
		return nil, err
	}

	log.Infoln("Pgx pool initialization successful")
	return pool, nil
}

func processMigration(migrationURL string, dbSource string, log *zap.SugaredLogger) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatalln("Failed to create migration: ", err)
	}

	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalln("Failed to migrate: ", err)
	}
	defer migration.Close()

	log.Infoln("Migration successful")
}
