.SILENT:

LOCAL_BIN=$(CURDIR)/bin
XO_OUTPUT_PATH="./internal/repository/entities"

DB_HOST="localhost"
DB_USER="diana"
DB_PORT="6532"
DB_NAME="wb"
DB_PASSWORD=""
DB_SSLMODE="disable"
DB_MIGRATIONS="./tools/migrations"

MAIN_PATH="./cmd/main.go"

.PHONY: bin-deps
bin-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.20.0

.PHONY: gen-db
gen-db:
	rm -r $(XO_OUTPUT_PATH)
	mkdir -p $(XO_OUTPUT_PATH)
	$(LOCAL_BIN)/xo schema "postgres://$(DB_USER)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable" \
		-o $(XO_OUTPUT_PATH) --schema public

.PHONY: add-migrations
add-migrations:
	$(LOCAL_BIN)/goose -dir  $(DB_MIGRATIONS) -no-color create $(NAME) sql

.PHONY: migrations-up-to
migrations-up-to:
	$(LOCAL_BIN)/goose -dir $(DB_MIGRATIONS) -no-color postgres "host=$(DB_HOST) user=$(DB_USER) dbname=$(DB_NAME) password=$(DB_PASSWORD) sslmode=$(DB_SSLMODE) port=$(DB_PORT)" up-to $(NAME)

.PHONY: migrations-down-to
migrations-down-to:
	$(LOCAL_BIN)/goose -dir $(DB_MIGRATIONS) -no-color postgres "host=$(DB_HOST) user=$(DB_USER) dbname=$(DB_NAME) password=$(POSTGRES_PASSWORD) sslmode=$(DB_SSLMODE) port=$(DB_PORT)" down-to $(NAME)

.PHONY: run
run:
	go run $(MAIN_PATH)
