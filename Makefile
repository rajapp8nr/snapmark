APP=snapmark

.PHONY: build build-linux build-windows build-mac run

build:
	go build -o bin/$(APP) ./...

build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/$(APP)-linux ./...

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/$(APP).exe ./...

build-mac:
	GOOS=darwin GOARCH=amd64 go build -o bin/$(APP)-mac ./...

run:
	go run .
