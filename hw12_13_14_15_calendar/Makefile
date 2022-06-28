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

migrate:
	go run ./migrations --config=./configs/config.yaml
