# Image configuration
REGISTRY ?= ghcr.io/kelos-dev
VERSION ?= latest
IMAGE_DIRS ?= cmd/kelos-controller cmd/kelos-spawner cmd/kelos-token-refresher claude-code codex gemini opencode cursor
LOCAL_ARCH ?= $(shell go env GOARCH)

# Version injection for the kelos CLI – only set ldflags when an explicit
# version is given so that dev builds fall through to runtime/debug info.
VERSION_PKG = github.com/kelos-dev/kelos/internal/version
ifneq ($(VERSION),latest)
LDFLAGS ?= -X $(VERSION_PKG).Version=$(VERSION)
endif

# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION = 1.31.0

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: test
test: ## Run unit tests.
	go test $$(go list ./... | grep -v /test/) --skip=E2E

.PHONY: test-integration
test-integration: envtest ## Run integration tests (envtest).
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path)" go test ./test/integration/... -v

.PHONY: test-e2e
test-e2e: ginkgo ## Run e2e tests (requires cluster and agent credentials).
	$(GINKGO) -v --timeout 20m ./test/e2e/...

.PHONY: update
update: controller-gen yamlfmt shfmt ## Run all generators and formatters.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."
	hack/update-install-manifest.sh $(CONTROLLER_GEN)
	hack/update-codegen.sh
	go fmt ./...
	go mod tidy
	$(YAMLFMT) .
	find . -name '*.sh' -not -path './bin/*' -exec $(SHFMT) -w -i 2 -ci {} +

.PHONY: verify
verify: controller-gen yamlfmt shfmt ## Verify everything is up-to-date and correct.
	@hack/verify.sh $(CONTROLLER_GEN) $(YAMLFMT) $(SHFMT)
	go vet ./...

##@ Build

.PHONY: build
build: ## Build binaries (use WHAT=cmd/kelos to build specific binary).
	@for dir in $$(go list ./$(or $(WHAT),cmd/...)); do \
		bin_name=$$(basename $$dir); \
		CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/$$bin_name $$dir; \
	done

.PHONY: run
run: ## Run a controller from your host.
	go run ./cmd/kelos-controller

.PHONY: image
image: ## Build docker images for the local platform (use WHAT to build specific image).
	@for dir in $(filter cmd/%,$(or $(WHAT),$(IMAGE_DIRS))); do \
		name=$$(basename $$dir); \
		GOOS=linux GOARCH=$(LOCAL_ARCH) $(MAKE) build WHAT=$$dir; \
		cp bin/$$name bin/$${name}-linux-$(LOCAL_ARCH); \
	done
	@GOOS=linux GOARCH=$(LOCAL_ARCH) $(MAKE) build WHAT=cmd/kelos-capture
	@cp bin/kelos-capture bin/kelos-capture-linux-$(LOCAL_ARCH)
	@for dir in $(or $(WHAT),$(IMAGE_DIRS)); do \
		docker build -t $(REGISTRY)/$$(basename $$dir):$(VERSION) -f $$dir/Dockerfile . --build-arg TARGETARCH=$(LOCAL_ARCH); \
	done

.PHONY: push
push: ## Push docker images (use WHAT to push specific image).
	@for dir in $(or $(WHAT),$(IMAGE_DIRS)); do \
		docker push $(REGISTRY)/$$(basename $$dir):$(VERSION); \
	done

IMAGE_PLATFORMS ?= linux/amd64,linux/arm64
IMAGE_ARCHES = $(shell echo "$(IMAGE_PLATFORMS)" | tr ',' '\n' | cut -d'/' -f2 | tr '\n' ' ')

.PHONY: image-multiarch
image-multiarch: ## Build and push multi-platform docker images via buildx (use WHAT to build specific image).
	@for dir in $(filter cmd/%,$(or $(WHAT),$(IMAGE_DIRS))); do \
		name=$$(basename $$dir); \
		for arch in $(IMAGE_ARCHES); do \
			GOOS=linux GOARCH=$$arch $(MAKE) build WHAT=$$dir; \
			mv bin/$$name bin/$${name}-linux-$$arch; \
		done; \
	done
	@for arch in $(IMAGE_ARCHES); do \
		GOOS=linux GOARCH=$$arch $(MAKE) build WHAT=cmd/kelos-capture; \
		mv bin/kelos-capture bin/kelos-capture-linux-$$arch; \
	done
	@for dir in $(or $(WHAT),$(IMAGE_DIRS)); do \
		docker buildx build --platform $(IMAGE_PLATFORMS) \
			-t $(REGISTRY)/$$(basename $$dir):$(VERSION) \
			-f $$dir/Dockerfile . --push; \
	done

RELEASE_PLATFORMS ?= linux/amd64 linux/arm64 darwin/amd64 darwin/arm64

.PHONY: release-binaries
release-binaries: ## Cross-compile CLI binaries for release and generate checksums.
	@for platform in $(RELEASE_PLATFORMS); do \
		os=$${platform%/*}; \
		arch=$${platform#*/}; \
		GOOS=$$os GOARCH=$$arch $(MAKE) build WHAT=cmd/kelos; \
		mv bin/kelos "bin/kelos-$${os}-$${arch}"; \
	done
	@cd bin && sha256sum kelos-* > checksums.txt

.PHONY: clean
clean: ## Clean build artifacts.
	rm -rf bin/
	rm -f cover.out

##@ Tool Dependencies

## Tool Binaries
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
ENVTEST ?= $(LOCALBIN)/setup-envtest
GINKGO ?= $(LOCALBIN)/ginkgo
YAMLFMT ?= $(LOCALBIN)/yamlfmt
SHFMT ?= $(LOCALBIN)/shfmt

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN)
$(CONTROLLER_GEN): $(LOCALBIN)
	test -s $(LOCALBIN)/controller-gen || GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen

.PHONY: envtest
envtest: $(ENVTEST)
$(ENVTEST): $(LOCALBIN)
	test -s $(LOCALBIN)/setup-envtest || GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest

.PHONY: ginkgo
ginkgo: $(GINKGO)
$(GINKGO): $(LOCALBIN)
	test -s $(LOCALBIN)/ginkgo || GOBIN=$(LOCALBIN) go install github.com/onsi/ginkgo/v2/ginkgo

.PHONY: yamlfmt
yamlfmt: $(YAMLFMT)
$(YAMLFMT): $(LOCALBIN)
	test -s $(LOCALBIN)/yamlfmt || GOBIN=$(LOCALBIN) go install github.com/google/yamlfmt/cmd/yamlfmt

.PHONY: shfmt
shfmt: $(SHFMT)
$(SHFMT): $(LOCALBIN)
	test -s $(LOCALBIN)/shfmt || GOBIN=$(LOCALBIN) go install mvdan.cc/sh/v3/cmd/shfmt
