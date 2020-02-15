build:
	@echo Compiling gotypist.
	go build -o gotypist ./src

run: build 
	./gotypist
