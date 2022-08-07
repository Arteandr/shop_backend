package main

import "shop_backend/internal/app"

const configPath = "configs"

// @title FinlandShop API
// @version 0.5
// @description API server

// @host localhost:8000
// @BasePath /api/v1

// @securityDefinitions.apikey UsersAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey AdminAuth
// @in context
// @name Admin authorization
func main() {
	app.Run(configPath)
}
