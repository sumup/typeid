# Make this makefile self-documented with target `help`
.PHONY: help
.DEFAULT_GOAL := help
help: ## Show help
	@grep -Eh '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: fmt
fmt: ## Format go files
	goimports -w .

.PHONY: lint
lint: ## Lint go files
	golangci-lint run -v

.PHONY: test
test: ## Run tests
	go test -v -failfast -race -timeout 1m ./...

.PHONY: generate
generate: ## Generate files
	go generate ./...

.PHONY: download
download: ## Download dependencies
	@echo Download go.mod dependencies
	@go mod download

.PHONE: vulncheck
vulncheck: ## Check for Vulnerabilities (make sure you have the tools install: `make install-tools`)
	govulncheck ./...
