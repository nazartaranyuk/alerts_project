SHELL := /bin/bash

PG_CONT ?= pg
PG_IMG  ?= postgres:16
PG_USER ?= myuser
PG_PASS ?= secret
PG_DB   ?= alarms
PG_PORT ?= 5432
PG_VOL  ?= pgdata

GO_MAIN ?= ./cmd/main.go
PG_DSN  ?= postgres://$(PG_USER):$(PG_PASS)@localhost:$(PG_PORT)/$(PG_DB)?sslmode=disable
JWT_SECRET ?= devsecretchange

.PHONY: up db-init run dev psql tables logs stop down reset-db help

up:
	@if ! docker ps -a --format '{{.Names}}' | grep -w '$(PG_CONT)' >/dev/null; then \
	  docker volume inspect $(PG_VOL) >/dev/null 2>&1 || docker volume create $(PG_VOL) >/dev/null; \
	  docker run -d --name $(PG_CONT) \
	    -e POSTGRES_USER=$(PG_USER) \
	    -e POSTGRES_PASSWORD=$(PG_PASS) \
	    -p $(PG_PORT):5432 \
	    --health-cmd='pg_isready -U $(PG_USER)' \
	    --health-interval=5s --health-timeout=3s --health-retries=10 \
	    --restart unless-stopped \
	    -v $(PG_VOL):/var/lib/postgresql/data \
	    $(PG_IMG); \
	else \
	  docker start $(PG_CONT) >/dev/null; \
	fi
	@echo "postgres up on port $(PG_PORT)"

db-init: up
	@docker exec -e PGPASSWORD="$(PG_PASS)" -i $(PG_CONT) \
	  psql -U "$(PG_USER)" -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='$(PG_DB)'" | grep -q 1 || \
	  docker exec -e PGPASSWORD="$(PG_PASS)" -i $(PG_CONT) \
	    psql -U "$(PG_USER)" -d postgres -c "CREATE DATABASE $(PG_DB)"
	@docker exec -i $(PG_CONT) bash -lc "cat >/tmp/schema.sql" < schema.sql
	@docker exec -e PGPASSWORD="$(PG_PASS)" -i $(PG_CONT) \
	  psql -U "$(PG_USER)" -d "$(PG_DB)" -f /tmp/schema.sql
	@echo "schema applied to $(PG_DB)"

run: db-init
	@PG_DSN="$(PG_DSN)" JWT_SECRET="$(JWT_SECRET)" go run $(GO_MAIN)

dev: db-init
	@PG_DSN="$(PG_DSN)" JWT_SECRET="$(JWT_SECRET)" reflex -r '\.go$$' -- sh -c 'go run $(GO_MAIN)'

psql:
	@docker exec -it $(PG_CONT) psql -U $(PG_USER) -d $(PG_DB)

tables:
	@docker exec -it $(PG_CONT) psql -U $(PG_USER) -d $(PG_DB) -c "\dt"

logs:
	@docker logs -f $(PG_CONT)

stop:
	@docker stop $(PG_CONT) >/dev/null || true
	@echo "stopped $(PG_CONT)"

down:
	@docker rm -f $(PG_CONT) >/dev/null || true
	@echo "removed $(PG_CONT)"

reset-db: down
	@docker volume rm -f $(PG_VOL) >/dev/null || true
	@$(MAKE) up db-init

help:
	@echo "make up        # start postgres"
	@echo "make db-init   # create db + apply schema.sql"
	@echo "make run       # run go app with PG_DSN"
	@echo "make psql      # open psql shell"
	@echo "make tables    # list tables"
	@echo "make logs      # docker logs -f pg"
	@echo "make stop      # stop container"
	@echo "make down      # remove container"
	@echo "make reset-db  # drop volume and recreate"
