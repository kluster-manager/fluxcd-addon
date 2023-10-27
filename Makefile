# Image URL to use all building/pushing image targets
IMG ?= docker.io/imtiazcho/fluxcd-addon:latest

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development
.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

##@ Build
.PHONY: build
build: fmt vet ## Build manager binary.
	GOFLAGS="" CGO_ENABLED=0 go build -o bin/fluxcd-addon cmd/main.go

.PHONY: run
run: fmt vet ## Run a controller from your host.
	go run cmd/main.go

.PHONY: docker-build
docker-build: ## Build docker image with the manager.
	docker build --no-cache -t ${IMG} .

.PHONY: docker-push
docker-push: docker-build ## Build and Push docker image with the manager.
	docker push ${IMG}

.PHONY: deploy-addon
deploy-addon: ## Deploy addon manifests to the hub cluster
	 kustomize build deploy/raw/ | kubectl apply -f -


.PHONY: undeploy-addon
undeploy-addon: ## Delete deployed manifests from the hub cluster
	kubectl delete -k deploy/raw --ignore-not-found

.PHONY: deploy-crd
deploy-crd: ## Apply flux config crd
	cd api/ && make manifests
	kustomize build api/config/crd/ | kubectl apply -f -

.PHONY: deploy-addon-all
deploy-addon-all:
	make deploy-crd
	make undeploy-addon
	make docker-push
	make deploy-addon