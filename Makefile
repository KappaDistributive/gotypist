run: build 
	./gotypist

build:
	@echo Compiling gotypist.
	go build -o gotypist ./v1

test:
	@echo Testing gotypist.
	go test -v ./v1

clean:
	@echo Deleting config.
	rm -r ~/.config/gotypist
