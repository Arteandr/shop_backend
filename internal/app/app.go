package app

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"shop_backend/internal/config"
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
	db.Exec("SELECT *")
	// HTTP server

}
