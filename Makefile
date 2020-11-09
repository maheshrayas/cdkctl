# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
BINARY_NAME=cdkctl
SYSTEMS=darwin linux windows

build: 
	$(GOBUILD) -o $(BINARY_NAME) -v

local:
	go build -o cdkctl ./cmd
	mv cdkctl /usr/local/bin

windows:
	env GOOS=windows GOARCH=386  go build -o cdkctl ./cmd

run:
	go run ./cmd/main.go

releases:
	$(foreach SYSTEM, $(SYSTEMS), \
	CGO_ENABLED=0 GOOS=$(SYSTEM) GOARCH=amd64 $(GOBUILD) -o release/$(SYSTEM)/$(BINARY_NAME) ./cmd; \
	cd release/$(SYSTEM)/; \
	tar -zcvf $(SYSTEM)-cdkctl.tar.gz cdkctl; \
	cd ../../; \
	)