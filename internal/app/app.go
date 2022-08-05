package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"os/signal"
	"shop_backend/internal/config"
	delivery "shop_backend/internal/delivery/http"
	"shop_backend/internal/repository"
	"shop_backend/internal/server"
	"shop_backend/internal/service"
	"shop_backend/pkg/auth"
	"shop_backend/pkg/hash"
	"shop_backend/pkg/logger"
	"syscall"
	"time"
)

func Run(configPath string) {
	// Config
	cfg, err := config.Init(configPath)
	if err != nil {
		logger.Error(err)
		return
	}

	// DB
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.PGSQL.Host, cfg.PGSQL.Port, cfg.PGSQL.User, cfg.PGSQL.Password, cfg.PGSQL.DatabaseName, cfg.PGSQL.SSLMode)
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		logger.Error("[DATABASE] " + err.Error())
		return
	}

	// Migrations
	instance, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		logger.Error("[DB INSTANCE] " + err.Error())
		return
	}
	m, err := migrate.NewWithDatabaseInstance("file://./schema", "postgres", instance)
	if err != nil {
		logger.Error("[MIGRATE] " + err.Error())
		return
	}
	m.Up()

	// Hasher
	hasher := hash.NewSHA1Hasher(cfg.Auth.PasswordSalt)

	// Token manager
	tokenManager, err := auth.NewManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		logger.Error("[AUTH] " + err.Error())
		return
	}

	// Services and repositories
	repos := repository.NewRepositories(db)
	services := service.NewServices(service.ServicesDeps{
		Repos:           repos,
		Hasher:          hasher,
		AccessTokenTTL:  cfg.Auth.AccessTokenTTL,
		RefreshTokenTTL: cfg.Auth.RefreshTokenTTL,
		TokenManager:    tokenManager,
	})

	handlers := delivery.NewHandler(services, cfg)

	// HTTP server
	srv := server.NewServer(cfg, handlers.Init(cfg))

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()
	logger.Info("server started")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	const timeout = 5 * time.Second
	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logger.Errorf("failed to stop server: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logger.Error(err.Error())
	}
}
