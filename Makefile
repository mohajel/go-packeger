F := .

run:
	go run main.go $(F)

build:
	go build -o ./bin/packager

install:
	go install

clean:
	rm -rf ./bin
	rm -rf ./packager-result

.PHONY: run build install clean