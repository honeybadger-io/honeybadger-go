all: test

prepare:
	# needed for `make fmt`
	go get golang.org/x/tools/cmd/goimports
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	# needed for `make cover`
	go get golang.org/x/tools/cmd/cover
	@echo Now you should be ready to run "make"

test:
	@go test -parallel 4 -race ./...

# goimports produces slightly different formatted code from go fmt
fmt:
	find . -name "*.go" -exec goimports -w {} \;

lint:
	golangci-lint run

cover:
	go test -cover -coverprofile cover.out
	go tool cover -html=cover.out

.PHONY: all prepare test fmt lint cover
