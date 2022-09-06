migrate_up:
	migrate -path ./schema -database 'postgres://admin:pgsqlpassword@localhost:5432/shop?sslmode=disable' up
migrate_down:
	migrate -path ./schema -database 'postgres://admin:pgsqlpassword@localhost:5432/shop?sslmode=disable' down
up:
	 docker-compose build && docker-compose down && docker-compose up -d
swag:
	swag init -g cmd/app.go
migrate_new:
	migrate create -ext sql -dir ./schema -seq $(name)
build:
	docker build -t hwndrer/backend .