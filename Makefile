migrate_up:
	migrate -path ./schema -database 'postgres://admin:pgsqlpassword@localhost:5432/shop?sslmode=disable' up
migrate_down:
	migrate -path ./schema -database 'postgres://admin:pgsqlpassword@localhost:5432/shop?sslmode=disable' down
up:
	docker-compose down && docker-compose build && docker-compose up -d