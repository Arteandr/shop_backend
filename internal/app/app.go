package app

import (
	"fmt"
	"log"
	"shop_backend/internal/config"
)

func Run(configPath string) {
	cfg, err := config.Init(configPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(cfg.HTTP.Port)
}
