run:
    go run .

build:
    go build -o ./bin/packager

install:
    go install

clean:
    rm -rf ./bin

.PHONY: run build install clean