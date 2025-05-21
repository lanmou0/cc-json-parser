.PHONY: build run test fmt clean

build:
	go build -o bin/ccjp src/*.go

run: build
	./bin/ccjp $(ARGS)

test:
	go test ./...

fmt:
	go fmt ./...

clean:
	rm -rf bin/
