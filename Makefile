.PHONY: build test lint bundle setup clean vuln

BINARY_DIR := bin
TOOLS_DIR := tools
COVERAGE_OUT := coverage.out

build:
	go build -o $(BINARY_DIR)/atheon ./cmd/atheon
	go build -o $(BINARY_DIR)/atheon-mcp ./cmd/mcp

test:
	go test ./... -p 1 -timeout 15m -coverprofile=$(COVERAGE_OUT)
	go tool cover -func=$(COVERAGE_OUT) | grep total:

test-race:
	go test ./... -p 1 -race -timeout 15m -coverprofile=$(COVERAGE_OUT)

test-junit:
	go install github.com/jstemmer/go-junit-report/v2@latest
	go test ./... -p 1 -v -timeout 15m -coverprofile=$(COVERAGE_OUT) 2>&1 | go-junit-report -set-exit-code > report.xml

lint:
	go vet ./...
	gofmt -l . | xargs -r false
	$(TOOLS_DIR)/golangci-lint run --timeout=5m 2>/dev/null || golangci-lint run --timeout=5m 2>/dev/null || true

bundle:
	go run ./bundler

vuln:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

setup:
	git config core.hooksPath .githooks
	mkdir -p $(TOOLS_DIR)
	GOBIN=$(PWD)/$(TOOLS_DIR) go install golang.org/x/tools/cmd/goimports@latest 2>/dev/null && echo "  installed: tools/goimports" || echo "  skipped: goimports (proxy blocked)"
	GOBIN=$(PWD)/$(TOOLS_DIR) go install honnef.co/go/tools/cmd/staticcheck@latest 2>/dev/null && echo "  installed: tools/staticcheck" || echo "  skipped: staticcheck (proxy blocked)"

clean:
	rm -rf $(BINARY_DIR)/ $(COVERAGE_OUT) report.xml
