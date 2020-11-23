.PHONY: test
test:
	go test ./...

.PHONY: generate
generate:
	go generate ./...

.PHONY: build
build:
	go build ./cmd/aiprime
