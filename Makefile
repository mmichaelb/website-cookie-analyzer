PROJECT_NAME=website-cookie-analyzer

GIT_VERSION=$(shell git describe --always)
GIT_BRANCH=$(shell git branch --show-current)

LD_FLAGS = -X main.GitVersion=${GIT_VERSION} -X main.GitBranch=${GIT_BRANCH}

OUTPUT_SUFFIX=$(go env GOEXE)

OUTPUT_PREFIX=./bin/${PROJECT_NAME}-${GIT_VERSION}

# test go program
test:
	@go test ./...

# builds and formats the project with the built-in Golang tool
.PHONY: build
build:
	@go build -ldflags '${LD_FLAGS}' -o "${OUTPUT_PREFIX}-${GOOS}-${GOARCH}${OUTPUT_FILE_ENDING}" ./cmd/websitecookieanalyzer/main.go

# installs and formats the project with the built-in Golang tool
install:
	@go install -ldflags '${LD_FLAGS}' ./cmd/websitecookieanalyzer/main.go
