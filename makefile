BINARY := ai-orchestra

.PHONY: build test run clean install

build:
	go build -o bin/$(BINARY) .

test:
	go test ./...

run: build
	./bin/$(BINARY)

install:
	go install .

clean:
	rm -rf bin/
