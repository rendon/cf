default: bin/cf
SRC=$(shell ls *.go)
bin/cf: $(SRC)
	go fmt
	golint
	go vet
	go build -ldflags "-s" -o bin/cf

install: bin/cf
	mkdir -p ~/bin
	cp bin/cf ~/bin/
test: bin/cf
	go test

coverage-test:
	go test -coverprofile=coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

clean:
	rm -vf bin/cf
