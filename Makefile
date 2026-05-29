.PHONY: all build test clean run dev lint install uninstall docker

BINARY_NAME=gitinspect
GO=go
INSTALL_DIR?=/usr/local/bin

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

install: build
	install -d $(DESTDIR)$(INSTALL_DIR)
	install -m 755 $(BINARY_NAME) $(DESTDIR)$(INSTALL_DIR)/$(BINARY_NAME)

uninstall:
	rm -f $(DESTDIR)$(INSTALL_DIR)/$(BINARY_NAME)

docker:
	docker build -t gitinspect .
