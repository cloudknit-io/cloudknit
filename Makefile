SHELL = /bin/bash

CGO_ENABLED = 0
DOCKER_TAG ?= latest
DOCKER_IMG = zlifecycle-il-operator

BUILD_CMD = go build -a -o build/zlifecycle-il-operator-$${GOOS}-$${GOARCH}

.PHONY: all
all: deps check docker-build

.PHONY: deps
deps: download-controller-gen download-mockgen

.PHONY: check test
.ONESHELL:
check test: generate manifests
	set -e
	set -o pipefail

	@echo "Running static analysis..."
	go vet ./...

	@echo "Running tests..."
	[[ ! -d test-results ]] && mkdir test-results || true
	go test -parallel 4 -covermode=count -coverprofile test-results/cover.out ./... \
		| tee test-results/go-test.out

	# Generate a human-readable, detailed coverage report in HTML
	@echo "Generating HTML coverage report..."
	GO111MODULE=off go get -u github.com/axw/gocov/gocov
	GO111MODULE=off go get -u github.com/matm/gocov-html
	gocov convert test-results/cover.out | gocov-html > test-results/coverage.html

.PHONY: docker-build
docker-build:
	go generate ./...
	docker build -t $(DOCKER_IMG):$(DOCKER_TAG) .

.PHONY: docker-dev-build
docker-dev-build:
	go generate ./...
	docker build . -t $(DOCKER_IMG):$(DOCKER_TAG) --file Dockerfile.dev

# Push the docker image to ECR -- reminder: never push 'latest'
.PHONY: docker-push
.ONESHELL:
docker-push:
ifndef ECR_REPO
	echo "ECR_REPO environment variable must be set before running 'make docker-push'"
	exit 1
endif
	set -e
	aws --version
	aws ecr get-login-password --region us-east-1 \
		| docker login --username AWS --password-stdin $(ECR_REPO)
	docker tag $(DOCKER_IMG):$(DOCKER_TAG) $(ECR_REPO)/$(DOCKER_IMG):$(DOCKER_TAG)
	docker push $(ECR_REPO)/$(DOCKER_IMG):$(DOCKER_TAG)

.PHONY: clean
clean:
	rm -rf build testbin test-results

# Generate code
generate: controller-gen
	controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
	controller-gen crd rbac:roleName=zlifecycle-il-operator-manager-role webhook \
		paths="./..." output:crd:artifacts:config=charts/zlifecycle-il-operator/crds output:rbac:artifacts:config=charts/zlifecycle-il-operator/templates

manifests-local: controller-gen
	# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
	controller-gen crd rbac:roleName=zlifecycle-il-operator-manager-role webhook \
		paths="./..." output:crd:artifacts:config=charts/zlifecycle-il-operator/crds output:rbac:artifacts:config=charts/zlifecycle-il-operator/templates

fmt: ## Run go fmt against code.
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...

# Run against the configured Kubernetes cluster in ~/.kube/config
# without first building the binary
run: generate manifests
	go run ./main.go

# Ensure controller-gen is available
.PHONY: controller-gen
download-controller-gen:
ifeq (, $(shell command -v controller-gen))
	@{ \
	echo "Downloading controller-gen..."
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.7.0 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
endif

.PHONY: mockgen
download-mockgen:
ifeq (, $(shell command -v mockgen))
	@{ \
	echo "Downloading mockgen..."
	set -e ;\
	MOCKGEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$MOCKGEN_TMP_DIR ;\
	go mod init tmp ;\
	go get github.com/golang/mock/mockgen@v1.5.0 ;\
	rm -rf $$MOCKGEN_TMP_DIR ;\
	}
endif
