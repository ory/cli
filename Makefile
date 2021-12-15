SHELL=/bin/bash -o pipefail

#  EXECUTABLES = docker-compose docker node npm go
#  K := $(foreach exec,$(EXECUTABLES),\
#          $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH")))

export GO111MODULE := on
export PATH := .bin:${PATH}
export PWD := $(shell pwd)

GO_DEPENDENCIES = github.com/ory/go-acc \
				  github.com/ory/x/tools/listx \
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
$(call make-lint-dependency)

.bin/clidoc: Makefile go.mod go.sum cmd
		go build -tags nodev -o .bin/clidoc ./cmd/clidoc/.

docs/cli: .bin/clidoc
		curl -o docs/sidebar.json https://raw.githubusercontent.com/ory/docs/master/docs/sidebar.json
		clidoc .

.bin/cli: go.mod go.sum Makefile
		go build -o .bin/cli -tags sqlite github.com/ory/cli

.bin/golangci-lint: Makefile
		bash <(curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh) -d -b .bin v1.28.3

.PHONY: lint
lint: .bin/golangci-lint
		golangci-lint run -v ./...

.PHONY: install
install:
		GO111MODULE=on go install -tags sqlite .

.PHONY: test
test:
		go test -p 1 -tags sqlite -count=1 -failfast ./...

# Formats the code
.PHONY: format
format: .bin/goimports
		goimports -w -local github.com/ory .

# Runs tests in short mode, without database adapters
.PHONY: docker
docker:
		docker build -f .docker/Dockerfile-build -t oryd/ory:latest-sqlite .
