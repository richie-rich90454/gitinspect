.PHONY: all build test clean run dev lint

BINARY_NAME=gitinspect
GO=go

all: build

build:
	$(GO) build -o $(BINARY_NAME) ./cmd/gitinspect

test:
	$(GO) test -v ./...

clean:
	rm -f $(BINARY_NAME)
	$(GO) clean ./...

run: build
	./$(BINARY_NAME) .

dev: build
	./$(BINARY_NAME) --format text --strip .

lint:
	golangci-lint run
