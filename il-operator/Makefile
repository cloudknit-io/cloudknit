SHELL = /bin/bash

export CGO_ENABLED = 0
export DOCKER_IMG = zlifecycle-il-operator:latest
export VERSION = $(shell git describe --always --tags 2>/dev/null || echo "initial")

BUILD_CMD = go build -a -o build/zlifecycle-il-operator-$${GOOS}-$${GOARCH}

.PHONY: all
all: deps check docker-build

.PHONY: deps
deps: download-controller-gen download-kustomize

.PHONY: check test
.ONESHELL:
check test: generate manifests
	set -e
	set -o pipefail

	@echo "Setting up envtest..."
	ENVTEST_ASSETS_DIR="$${PWD}/testbin"
	[[ ! -d $$ENVTEST_ASSETS_DIR ]] && mkdir -p $${ENVTEST_ASSETS_DIR}
	if [[ ! -f $$ENVTEST_ASSETS_DIR/setup-envtest.sh ]]; then
		pushd $$ENVTEST_ASSETS_DIR
		curl -fsSLO https://raw.githubusercontent.com/kubernetes-sigs/controller-runtime/v0.6.3/hack/setup-envtest.sh
		popd
	fi

	source $${ENVTEST_ASSETS_DIR}/setup-envtest.sh
	fetch_envtest_tools $${ENVTEST_ASSETS_DIR}
	setup_envtest_env $${ENVTEST_ASSETS_DIR}

	@echo "Running static analysis..."
	go vet ./...

	@echo "Running tests..."
	[[ ! -d test-results ]] && mkdir test-results || true
	go test -parallel 2 -covermode=count -coverprofile test-results/cover.out ./... \
		| tee test-results/go-test.out

	# Generate a human-readable, detailed coverage report in HTML
	@echo "Generating HTML coverage report..."
	GO111MODULE=off go get -u github.com/axw/gocov/gocov
	GO111MODULE=off go get -u github.com/matm/gocov-html
	gocov convert test-results/cover.out | gocov-html > test-results/coverage.html

.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_IMG) .
	docker tag $(DOCKER_IMG) $${DOCKER_IMG%:*}:latest

.PHONY: docker-dev-build 
docker-dev-build:
	docker build . -t ${DOCKER_IMG} --file Dockerfile.dev
	docker tag $(DOCKER_IMG) $${DOCKER_IMG%:*}:latest

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
	docker tag $(DOCKER_IMG) $(ECR_REPO)/$(DOCKER_IMG)
	docker push $(ECR_REPO)/$(DOCKER_IMG)

.PHONY: clean
clean:
	rm -rf build testbin test-results

# Generate code
generate: controller-gen
	controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
	controller-gen crd:trivialVersions=true rbac:roleName=manager-role webhook \
		paths="./..." output:crd:artifacts:config=helm/zlifecycle-il-operator/templates output:rbac:artifacts:config=helm/zlifecycle-il-operator/templates

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
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.3.0 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
endif

.PHONY: kustomize
download-kustomize:
ifeq (, $(shell command -v kustomize))
	@{ \
	echo "Downloading kustomize..."
	set -e ;\
	KUSTOMIZE_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$KUSTOMIZE_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/kustomize/kustomize/v3@v3.5.4 ;\
	rm -rf $$KUSTOMIZE_GEN_TMP_DIR ;\
	}
endif
