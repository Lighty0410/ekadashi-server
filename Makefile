BIN_DIR := ./bin
BUILDNAME := $(BIN_DIR)/server

build:
	go build -o $(BUILDNAME) ./cmd/app

test: 
	go test -cover ./...

.PHONY:
clean:
	go clean
	rm -rf $(BIN_DIR)

.PHONY: lint
lint:
	gometalinter --vendor --disable-all \
		--enable=vet \
		--enable=gotype \
		--enable=deadcode \
		--enable=gocyclo \
		--enable=golint \
		--enable=varcheck \
		--enable=structcheck \
		--enable=maligned \
		--enable=errcheck \
		--enable=staticcheck \
		--enable=ineffassign \
		--enable=interfacer \
		--enable=unconvert \
		--enable=goconst \
		--enable=gofmt \
		--enable=goimports \
		--enable=misspell \
		--enable=unparam ./...