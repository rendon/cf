default: build
SRC=$(shell ls *.go)
bin/cf: $(SRC)
	go fmt
	go vet
	go build -ldflags "-s" -o bin/cf

install: bin/cf
	cp bin/cf ~/bin/
test:
	go test

coverage-test:
	go test -coverprofile=coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out