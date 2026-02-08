BINARY=agentsec

.PHONY: build test cover docs-smoke fmt fmt-check lint hooks clean

build:
	@mkdir -p bin
	go build -o bin/$(BINARY) ./cmd/agentsec

test:
	go test ./...

cover:
	@export GOCACHE=$$(pwd)/.gocache; \
	go test -coverprofile=.cover_internal.out ./internal/...; \
	total=$$(go tool cover -func=.cover_internal.out | awk '/^total:/{gsub("%","",$$3); print $$3}'); \
	awk -v t="$$total" 'BEGIN { if (t < 80.0) { printf "coverage %.1f%% is below 80.0%%\n", t; exit 1 } else { printf "coverage %.1f%% (threshold 80.0%%)\n", t } }'

docs-smoke: build
	bash scripts/docs-smoke.sh

fmt:
	gofmt -w .

fmt-check:
	@test -z "$$(gofmt -l .)" || (echo "Run 'make fmt' to format files." && gofmt -l . && exit 1)

lint:
	@LINT_BIN="./bin/golangci-lint"; \
	if [ ! -x "$$LINT_BIN" ]; then LINT_BIN="golangci-lint"; fi; \
	if ! command -v "$$LINT_BIN" >/dev/null 2>&1; then \
		echo "golangci-lint not found. Install it or place binary at ./bin/golangci-lint"; \
		exit 1; \
	fi; \
	HOME=$$(pwd) XDG_CACHE_HOME=$$(pwd)/.cache GOCACHE=$$(pwd)/.gocache GOLANGCI_LINT_CACHE=$$(pwd)/.golangci-cache "$$LINT_BIN" run

hooks:
	git config core.hooksPath .githooks
	chmod +x .githooks/pre-commit
	@echo "Configured git hooks to use .githooks"

clean:
	rm -rf bin _demo dist
