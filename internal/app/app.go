package app

import (
	"log"
	"shop_backend/internal/config"
)

func Run(configPath string) {
	// Config
	_, err := config.Init(configPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	// HTTP server

}
