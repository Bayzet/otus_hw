BIN := "./bin/calendar"
DB_USER := "root"
DB_PASS := "YLmcjlpD3b96Bi08KPrf"

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) --config ./configs/config.yaml

test:
	go test -race ./internal/...

lint:
	golangci-lint run ./...

migrate-up:
	mysql -u $(DB_USER) -p $(DB_PASS) calendar < migrations/first_06122022.up.sql

migrate-down:
	mysql -u $(DB_USER) -p $(DB_PASS) calendar < migrations/first_06122022.down.sql