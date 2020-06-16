GIT_COMMIT := $(shell git rev-parse HEAD)
GIT_DIRTY := $(if $(shell git status --porcelain),+CHANGES)
GIT_TAG := $(shell git name-rev --tags --name-only $(git rev-parse HEAD))
GO_LDFLAGS := "-X github.com/matoous/linkfix/internal/version.GitCommit=$(GIT_COMMIT)$(GIT_DIRTY) -X github.com/matoous/linkfix/internal/version.GitTag=$(GIT_TAG)"

define echo_green
    $(call echo_custom_color,$(color_green),$(1))
endef

define echo_cyan
	$(call echo_custom_color,$(color_cyan),$(1))
endef

define echo_warning
	$(call echo_custom_color,$(color_yellow),$(1))
endef

# Make this makefile self-documented with target `help`
.PHONY: help
.DEFAULT_GOAL := help

help:
	@grep -Eh '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

configure: ## Prepare environment for compiling, e.g. download dependencies
	go mod download

build: configure ## Build binaries
	$(call prepare_build_vars)
	go build -o bin/linkfix --ldflags $(GO_LDFLAGS) main.go

test: configure ## Run tests
	go test -race ./... -cover

lint: go-mod-tidy lint-golangci ## Run all linters

go-mod-tidy: ## Check if go.mod and go.sum does not contains any unnecessary dependencies and remove them.
ifndef TMPDIR
	$(eval TMPDIR=$(shell mktemp -d))
endif
	cp -fv go.mod $(TMPDIR)
	cp -fv go.sum $(TMPDIR)
	go mod tidy -v
	diff -u $(TMPDIR)/go.mod go.mod
	diff -u $(TMPDIR)/go.sum go.sum
	rm -f $(TMPDIR)go.mod $(TMPDIR)go.sum

lint-golangci: ## Runs golangci-lint. It outputs to the code-climate json file in if CI is defined.
	golangci-lint run --max-same-issues 0 --max-issues-per-linter 0 $(if $(CI),--out-format code-climate > gl-code-quality-report.json 2>golangci-stderr-output)
