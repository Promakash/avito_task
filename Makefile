build:
	docker compose build

run:
	docker compose up -d

stop:
	docker compose down

build_and_run: build run

generate_docs:
	swag fmt
	swag init -g cmd/main.go -o docs

lint_code:
	golangci-lint run

run_tests:
	docker compose run --rm tests

run_unit_tests:
	go test ./internal/...
