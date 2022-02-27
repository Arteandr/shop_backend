migrate_up:
	migrate -path ./schema -database 'postgres://admin:pgsqlpassword@localhost:5433/shop?sslmode=disable' up
migrate_down:
	migrate -path ./schema -database 'postgres://admin:pgsqlpassword@localhost:5433/shop?sslmode=disable' down
