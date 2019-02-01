GOBIN=$(shell pwd)/bin
GOFILES=$(wildcard *.go)
GONAME=$(shell basename "$(PWD)")
PKGS := $(shell go list ./... | grep -v /vendor)

.PHONY: test
test: 
	go test $()

.PHONY: build
build:
  	@echo "Building $(GOFILES) to ./bin"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -o bin/$(GONAME) $(GOFILES)

.PHONY: clean
clean:
	@echo " CLEAN"
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
		--enable='gofmt -s' \ 
		--enable=goimports \ 
		--enable=`lll - Repor` \
		--enable=misspell \
		--enable=unparam \