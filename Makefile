
lint: .golangci.yml
	golangci-lint run

test:
ifeq ($(OS),windows)
	@echo "Skipping test on windows, issue with -- and testscript"
else
	go test -race ./...
endif

.golangci.yml: Makefile
	curl -fsS -o .golangci.yml https://raw.githubusercontent.com/fortio/workflows/main/golangci.yml

.PHONY: lint
