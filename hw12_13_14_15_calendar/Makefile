BIN := "./bin/calendar"
DB_USER := "root"
DB_PASS := "YLmcjlpD3b96Bi08KPrf"

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) --config ./configs/config.yaml

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

migrate:
	go run ./migrations --config=./configs/config.yaml

generate:
	go generate ./...
	protoc ./proto/calendar.proto --go_out=./internal/gen/pb --go-grpc_out=./internal/gen/pb

evans:
	evans --proto ./proto/calendar.proto repl
