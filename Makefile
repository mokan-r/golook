MODULE=github.com/mokan-r/golook

run: migrate_up build run_golook

build:
	go mod tidy
	go build -o golook ${MODULE}/cmd/golook

run_golook:
	./golook

run_postgres:
	docker-compose up --detach

migrate_up:
	migrate -path ./migrations -database 'postgres://golook:password@localhost:5432/golook?sslmode=disable' up

migrate_down:
	migrate -path ./migrations -database 'postgres://golook:password@localhost:5432/golook?sslmode=disable' down

create_migrations:
	migrate create -ext sql -dir ./migrations -seq create_users_table
