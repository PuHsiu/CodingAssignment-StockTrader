.PHONY: build

all: clean build

build:
	go build -o dist/trader

clean:
	rm -rf ./dist

test:
	go test ./...
