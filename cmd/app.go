package main

import "shop_backend/internal/app"

const configPath = "configs"

func main() {
	app.Run(configPath)
}
