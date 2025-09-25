local-run:
	go build -o bin/main cmd/main.go
	./bin/main
.PHONY: local-run

run:
	docker compose up -d
	docker compose logs -f app
.PHONY: run

docker-build:
	docker compose up -d --build
.PHONY: docker-build

docker-stop:
	docker compose down
.PHONY: docker-stop

app-logs:
	docker compose logs -f app
.PHONY: app-logs

test:
	go test -v ./... ./internal/...
.PHONY: test

swagger:
	swag init -g cmd/main.go -o internal/docs
.PHONY: swagger