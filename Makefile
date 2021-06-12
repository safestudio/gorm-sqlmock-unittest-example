test:
	go test -v -cover -covermode=atomic ./...

unittest:
	go test -short  ./...

clean:
	[ -d dist ] && rm dist/*

build:
	go build -o dist/app

lint-prepare:
	@echo "Installing golangci-lint"
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s latest

lint:
	./bin/golangci-lint run ./...

.PHONY: test unittest clean build lint-prepare lint