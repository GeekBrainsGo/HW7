APPNAME = "blog_srv"

.PHONY: build
build:
	go build -o $(APPNAME)

.PHONY: run
run: build
	./$(APPNAME)

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.PHONY: docs
docs:
	swag init

.PHONY: dev
dev: test docs run

.DEFAULT_GOAL := build