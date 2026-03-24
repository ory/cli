SHELL=/bin/bash -o pipefail

export PATH := .bin:${PATH}
export PWD := $(shell pwd)

GOLANGCI_LINT_VERSION = 2.11.4

.bin/clidoc: Makefile go.mod cmd
	go build -tags nodev -o .bin/clidoc ./cmd/clidoc/.

docs/cli: .bin/clidoc
	curl -o docs/sidebar.json https://raw.githubusercontent.com/ory/docs/master/docs/sidebar.json
	clidoc .

.bin/cli: go.mod Makefile
	go build -o .bin/cli -tags sqlite github.com/ory/cli

.bin/golangci-lint-$(GOLANGCI_LINT_VERSION):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b .bin v$(GOLANGCI_LINT_VERSION)
	mv .bin/golangci-lint .bin/golangci-lint-$(GOLANGCI_LINT_VERSION)

.bin/licenses: Makefile
	curl https://raw.githubusercontent.com/ory/ci/master/licenses/install | sh

.PHONY: lint
lint: .bin/golangci-lint-$(GOLANGCI_LINT_VERSION)
	.bin/golangci-lint-$(GOLANGCI_LINT_VERSION) run --timeout=10m ./...

.PHONY: install
install:
	go install -tags sqlite .

.PHONY: refresh
refresh:
	UPDATE_SNAPSHOTS=true go test -tags sqlite,json1,refresh ./...

# Formats the code
.PHONY: format
format: .bin/cli node_modules go.mod
	.bin/cli dev headers copyright --type=open-source
	go tool goimports -w -local github.com/ory .
	npm exec -- prettier --write "{**/,}*{.js,.md,.ts}"

licenses: .bin/licenses node_modules  # checks open-source licenses
	.bin/licenses

# Runs tests in short mode, without database adapters
.PHONY: docker
docker:
	docker build -f .docker/Dockerfile-build -t oryd/ory:latest-sqlite .


# Runs tests in short mode, without database adapters
.PHONY: post-release
post-release:
	echo "nothing to do"

node_modules: package-lock.json
	npm ci
	touch node_modules
