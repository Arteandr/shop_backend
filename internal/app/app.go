package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"os/signal"
	"shop_backend/internal/config"
	delivery "shop_backend/internal/delivery/http"
	"shop_backend/internal/repository"
	"shop_backend/internal/server"
	"shop_backend/internal/service"
	"syscall"
	"time"
)

func Run(configPath string) {
	// Config
	cfg, err := config.Init(configPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	// DB
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.PGSQL.Host, cfg.PGSQL.Port, cfg.PGSQL.User, cfg.PGSQL.Password, cfg.PGSQL.DatabaseName, cfg.PGSQL.SSLMode)
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
		return
	}

	repos := repository.NewRepositories(db)
	services := service.NewServices(service.ServicesDeps{
		Repos: repos,
	})

	handlers := delivery.NewHandler(services, cfg)

	// HTTP server
	srv := server.NewServer(cfg, handlers.Init(cfg))

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error occured while running http server: %s\n", err.Error())
		}
	}()
	fmt.Println("server started")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	const timeout = 5 * time.Second
	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		log.Fatalf("failed to stop server: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		log.Fatalf(err.Error())
	}
}
