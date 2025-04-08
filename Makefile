
lint: .golangci.yml
	golangci-lint run


OS:=$(shell go env GOOS)

test:
ifeq ($(OS),windows)
	@echo "Skipping test on windows, issue with -- and testscript"
else
	@echo "running tests on $(OS)"
	go test -race ./...
endif

.golangci.yml: Makefile
	curl -fsS -o .golangci.yml https://raw.githubusercontent.com/fortio/workflows/main/golangci.yml

.PHONY: lint
