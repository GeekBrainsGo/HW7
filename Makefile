APPNAME = "blog"
# h - help
h help:
	@echo "h help 	- this help"
	@echo "build 	- build and the app"
	@echo "run 	- run the app"
	@echo "clean 	- clean app trash"
.PHONY: h

# build - build the app
build:
	go build -o $(APPNAME)
.PHONY: build

# run - build and run the app
run: build
	./$(APPNAME)
.PHONY: run

clean:
	rm ./$(APPNAME)
.PHONY: clean
