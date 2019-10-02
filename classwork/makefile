# Project filename
NAME = "serv"
BENCH_COUNT = 5
BENCH_TIME = 1s

# h - help
h help:
	@echo "h help 	- this help"
	@echo "build 	- build and the app"
	@echo "run 	- run the app"
	@echo "clean 	- clean app trash"
	@echo "swag 	- generate swag docs"
	@echo "dev 	- generate docs and run"
	@echo "test 	- run all tests"
	@echo "bench 	- run all benchmarks"
.PHONY: h

# build - build the app
build:
	go build -o $(NAME)
.PHONY: build

# run - build and run the app
run: build
	./$(NAME)
.PHONY: run

clean:
	rm ./$(NAME)
.PHONY: clean

# swag - generate swagger docs
swag:
	swag init
.PHONY: swag

# dev - generate docs and run
dev: swag run
.PHONY: dev

# test - run all tests
test:
	go test ./...
.PHONY: test

# bench - run all benchmarks
bench:
	go test -bench=. ./... -count=$(BENCH_COUNT) -benchtime=$(BENCH_TIME) -benchmem
.PHONY: bench