.PHONY: all
all: test

.PHONY: clean
clean:
	rm ../golang-github-k0swe-wsjtx-go*

.PHONY: test
test:
	go test ./...
	go vet ./...
