BINARY_NAME=chrono
MAIN_PACKAGE=./cmd/chrono/

.PHONY: build clean install test coverage 

build:
	go build -o $(BINARY_NAME) $(MAIN_PACKAGE)

clean:
	rm -f $(BINARY_NAME)

install:
	go install $(MAIN_PACKAGE)

test:
	go test ./...

coverage:
	go test -cover ./...

.DEFAULT_GOAL := build
