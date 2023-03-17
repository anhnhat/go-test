build:
	go build -o out.out ./cmd/server/main.go

run:
	go run ./cmd/server/main.go

migrate/init:
	go run ./cmd/migrate/migrate.go migrate

migrate/drop:
	go run ./cmd/migrate/migrate.go drop

migrate/seed:
	go run ./cmd/migrate/migrate.go seed

docker/up:
	docker compose up -d

fmt:
	goimports-reviser -rm-unused -format ./...
