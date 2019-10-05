NAME = "swaggerlblog"

h help:
	@echo "h help - Help"
	@echo "build - Build project"
	@echo "run - Build project and run"
.PHONY: h help

build:
	go build -o $(NAME)
.PHONY: build

run: build
	./swaggerlblog
.PHONY: run

test:
	go test ./...
.PHONY: test

swag:
	swag init
.PHONY: swag

dev: swag test run
.PHONY: dev