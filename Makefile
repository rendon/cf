default: build

build:
	go fmt
	go vet
	go build -ldflags "-s" -o bin/cf

test: build
	go test

coverage-test:
	go test -coverprofile=coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out
