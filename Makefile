BIN_DIR := ./bin
BUILDNAME := $(BIN_DIR)/server
PROTO_DIR := ./pkg/server/grpc/api
PROTOFILE := $(PROTO_DIR)/server.proto

build:
	go build -o $(BUILDNAME) ./cmd/app
.PHONY:
test: 
	go test -cover ./...

.PHONY:
clean:
	go clean
	rm -rf $(BIN_DIR)

.PHONY:
proto:
	protoc -I $(PROTO_DIR) $(PROTOFILE) --go_out=plugins=grpc:$(PROTO_DIR)

.PHONY: lint
lint:
	golangci-lint run --disable-all \
		--enable=vet \
		--enable=typecheck \
		--enable=deadcode \
		--enable=gocyclo \
		--enable=golint \
		--enable=varcheck \
		--enable=structcheck \
		--enable=maligned \
		--enable=errcheck \
		--enable=ineffassign \
		--enable=interfacer \
		--enable=unconvert \
		--enable=goconst \
		--enable=gofmt \
		--enable=goimports \
		--enable=misspell \
		--enable=unparam \
		--enable=unused ./...

mod:
	go mod tidy
