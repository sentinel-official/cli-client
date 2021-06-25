PACKAGES := $(shell go list ./...)
VERSION := $(shell git rev-parse --short HEAD)
COMMIT := $(shell git log -1 --format='%H')

BUILD_TAGS := $(strip netgo,ledger)
LD_FLAGS := -s -w \
	-X github.com/cosmos/cosmos-sdk/version.Name=sentinel \
	-X github.com/cosmos/cosmos-sdk/version.AppName=sentinelcli \
	-X github.com/cosmos/cosmos-sdk/version.Version=${VERSION} \
	-X github.com/cosmos/cosmos-sdk/version.Commit=${COMMIT} \
	-X github.com/cosmos/cosmos-sdk/version.BuildTags=${BUILD_TAGS}

.PHONY: benchmark
benchmark:
	@go test -mod=readonly -v -bench= ${PACKAGES}

.PHONY: clean
clean:
	rm -rf ./build

.PHONY: build
build:
	GOOS=darwin GOARCH=amd64 go build -mod=readonly -tags="${BUILD_TAGS}" -ldflags="${LD_FLAGS}" \
		-o ./build/sentinelcli-${VERSION}-darwin-amd64 main.go
	GOOS=linux GOARCH=amd64 go build -mod=readonly -tags="${BUILD_TAGS}" -ldflags="${LD_FLAGS}" \
		-o ./build/sentinelcli-${VERSION}-linux-amd64 main.go

.PHONY: install
install:
	go build -mod=readonly -tags="${BUILD_TAGS}" -ldflags="${LD_FLAGS}" \
		-o ${GOPATH}/bin/sentinelcli main.go

.PHONY: go-lint
go-lint:
	@golangci-lint run --fix

.PHONY: test
test:
	@go test -mod=readonly -v -cover ${PACKAGES}

.PHONY: tools
tools:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.27.0
