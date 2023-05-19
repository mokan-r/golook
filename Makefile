MODULE=github.com/mokan-r/golook

run: run_postgres build

build:
	go mod tidy
	go build -o golook ${MODULE}/cmd/golook

run_postgres:
	docker-compose up --detach
