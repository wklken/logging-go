.PHONY: dep
dep:
	go mod tidy
	go mod vendor

.PHONY: lint
lint:
	export GOFLAGS=-mod=vendor
	golangci-lint run --fast -e examples/*

.PHONY: test
test:
	go test -mod=vendor -gcflags=all=-l $(shell go list ./... | grep -v examples) -covermode=count -coverprofile .coverage.cov
	go tool cover -func=.coverage.cov

.PHONY: tools
tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint
