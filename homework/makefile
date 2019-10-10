NAME = "serv"

h help:
	@echo "h help	-  Help"
	@echo "build	-  Compile"
	@echo "run	-  Compile and Run"
	@echo "clean	-  Remove compiled binary file"
	@echo "swag	-  Generate Swagger"
	@echo "test	-  Run tests"
.PHONY: h help

build:
	go build -o $(NAME)
.PHONY: build

run: build
	./$(NAME)
.PHONY: run

clean:
	rm -rf ./$(NAME)
.PHONY: clean

test:
	go test ./...
.PHONY: test

swag:
	swag init
.PHONY: swag

dev: swag test run
.PHONY: dev