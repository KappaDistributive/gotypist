run: build 
	./gotypist

dependencies:
	@echo Installing dependencies.
	go get -d ./v1

build: dependencies
	@echo Compiling gotypist.
	go build -o gotypist ./v1

test:
	@echo Testing gotypist.
	go test -v ./v1

clean:
	@echo Deleting config.
	if [ -d ~/.config/gotypist ]; then rm -r ~/.config/gotypist; fi
