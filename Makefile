SHELL := /bin/bash

TOOLS_MOD_DIR := ./tools
TOOLS_DIR := $(abspath ./.tools)
CRDS_DIR := $(abspath ./crds)

$(TOOLS_DIR)/controller-gen: $(TOOLS_MOD_DIR)/go.mod $(TOOLS_MOD_DIR)/go.sum $(TOOLS_MOD_DIR)/tools.go
	cd $(TOOLS_MOD_DIR) && \
	go build -o $(TOOLS_DIR)/controller-gen sigs.k8s.io/controller-tools/cmd/controller-gen

.DEFAULT_GOAL := all

.PHONY: all
all: ## build pipeline
all: mod gen build spell lint test

.PHONY: precommit
precommit: ## validate the branch before commit
precommit: all vuln

.PHONY: ci
ci: ## CI build pipeline
ci: precommit diff

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: clean
clean: ## remove files created during build pipeline
	rm -rf dist
	rm -f coverage.*
	rm -f '"$(shell go env GOCACHE)/../golangci-lint"'
	go clean -i -cache -testcache -modcache -fuzzcache -x

.PHONY: run
run: ## go run
	go run .

.PHONY: mod
mod: ## go mod tidy
	go mod tidy

.PHONY: gen
gen: ## go generate
	go generate ./...

.PHONY: build
build: ## goreleaser build
	go tool goreleaser build --clean --single-target --snapshot

.PHONY: spell
spell: ## misspell
	go tool misspell -error -locale=US -w **.md

.PHONY: lint
lint: ## golangci-lint
	go tool golangci-lint run --fix

.PHONY: vuln
vuln: ## govulncheck
	go tool govulncheck ./...

ifeq ($(CGO_ENABLED),0)
RACE_OPT =
else
RACE_OPT = -race
endif

.PHONY: test
test: ## go test
	go test $(RACE_OPT) -covermode=atomic -coverprofile=coverage.out -coverpkg=./... ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: diff
diff: ## git diff
	git diff --exit-code
	RES=$$(git status --porcelain) ; if [ -n "$$RES" ]; then echo $$RES && exit 1 ; fi

.PHONY: codegen
codegen:
	./hack/update-codegen.sh

.PHONY: crdgen
crdgen: $(TOOLS_DIR)/controller-gen
	$(TOOLS_DIR)/controller-gen \
		crd:crdVersions=v1 \
  		paths=./pkg/k8s/apis/cyberark/v1alpha1/... \
  		output:crd:artifacts:config=./crds
	mv $(CRDS_DIR)/biggs.cl_cyberarks.yaml $(CRDS_DIR)/CyberArk.yaml
