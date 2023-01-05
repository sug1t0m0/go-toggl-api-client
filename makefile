.PHONY: lint
lint:
	golangci-lint run -v

.PHONY: test
test:
	go test ./...
