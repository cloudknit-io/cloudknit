.PHONY: build
build: ## Build all binaries
	go build -gcflags "all=-N -l" -o build/debug/ .

.PHONY: build-release
build-release: ## Build all binaries without debug information, clean up after
	go build -ldflags '-s -w' -trimpath -o build/release/ .

.PHONY: run
run:
	go run main.go
