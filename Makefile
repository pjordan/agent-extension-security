BINARY=agentsec

.PHONY: build test fmt lint clean

build:
	@mkdir -p bin
	go build -o bin/$(BINARY) ./cmd/agentsec

test:
	go test ./...

fmt:
	gofmt -w .

lint:
	@echo "No linter configured yet. Consider golangci-lint."

clean:
	rm -rf bin _demo dist
