package main

import "shop_backend/internal/app"

const configPath = "configs"

// @title FinlandShop API
// @version 0.1
// @description API server

// @host localhost:8000
// @BasePath /api/v1
func main() {
	app.Run(configPath)
}
