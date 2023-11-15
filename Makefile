SHELL=/bin/bash -o pipefail

#  EXECUTABLES = docker-compose docker node npm go
#  K := $(foreach exec,$(EXECUTABLES),\
#          $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH")))

export GO111MODULE := on
export PATH := .bin:${PATH}
export PWD := $(shell pwd)

GOLANGCI_LINT_VERSION = 1.55.2

GO_DEPENDENCIES = github.com/ory/go-acc \
				  github.com/golang/mock/mockgen \
				  github.com/go-swagger/go-swagger/cmd/swagger \
				  golang.org/x/tools/cmd/goimports \
				  github.com/mikefarah/yq \
				  github.com/mattn/goveralls

define make-go-dependency
  # go install is responsible for not re-building when the code hasn't changed
  .bin/$(notdir $1): go.mod go.sum Makefile
		GOBIN=$(PWD)/.bin/ go install $1
endef
$(foreach dep, $(GO_DEPENDENCIES), $(eval $(call make-go-dependency, $(dep))))

.bin/clidoc: Makefile go.mod go.sum cmd
	go build -tags nodev -o .bin/clidoc ./cmd/clidoc/.

docs/cli: .bin/clidoc
	curl -o docs/sidebar.json https://raw.githubusercontent.com/ory/docs/master/docs/sidebar.json
	clidoc .

.bin/cli: go.mod go.sum Makefile
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
	GO111MODULE=on go install -tags sqlite .

.PHONY: test
test: lint
	go test -p 1 -tags sqlite -count=1 -failfast ./...

.PHONY: refresh
refresh:
	UPDATE_SNAPSHOTS=true go test -tags sqlite,json1,refresh ./...

# Formats the code
.PHONY: format
format: .bin/cli .bin/goimports node_modules
	.bin/cli dev headers copyright --type=open-source
	goimports -w -local github.com/ory .
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
