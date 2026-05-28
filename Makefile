.PHONY: all build test clean run lint

BINARY_NAME=gitinspect
GO=go

all: build

build:
	$(GO) build -o $(BINARY_NAME) ./cmd/gitinspect

test:
	$(GO) test -v ./internal/...

clean:
	rm -f $(BINARY_NAME)
	$(GO) clean ./...

run: build
	./$(BINARY_NAME)

lint:
	golangci-lint run
