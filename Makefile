VERSION := latest
GO_FLAGS := '-gcflags=all=-N -l ' -ldflags='-X main.version=$(VERSION)'
GOOS := linux
GOARCH := amd64
BINS := $(addprefix build/,$(notdir $(wildcard cmd/*)))
all: $(BINS)

$(BINS): $(shell find . -name '*.go')
	@mkdir -p $(dir $@)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build $(GO_FLAGS) -o $@ ./cmd/$(notdir $@)

.PHONY: release
release:
	$(MAKE) clean
	$(MAKE) GO_FLAGS="'-ldflags=-s -w -buildid= -X main.version=$(VERSION)' -trimpath -buildvcs=false" VERSION=$(VERSION)

DOCKER := docker
.PHONY: oci
oci: $(BINS)
	$(DOCKER) image build -t ms1-server:$(VERSION) -t ms1-server:latest .

.PHONY: integration
integration: oci
	go test -v -count 1 ./integration

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor
	go mod verify

.PHONY: lint
lint: build/tools/golangci-lint
	$< run --timeout=10m ./...

.PHONY: gen
gen:
	go generate ./...

.PHONY: test
test:
	go test -v ./...

build/tools/golangci-lint:
	mkdir -p $(dir $@)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(dir $@) v1.54.2

audit: build/tools/trivy
	$< fs --scanners vuln .

build/tools/trivy:
	mkdir -p $(dir $@)
	curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b $(dir $@) v0.45.1

.PHONY: ci
ci:
	dagger run go run ./ci

.PHONY: check-git-clean
check-git-clean:
	test $(shell git status -s | wc -l) -eq 0

.PHONY: clean
clean:
	rm -rf build bin

pbgen: build/tools/protoc build/tools/protoc-gen-go build/tools/protoc-gen-go-grpc go.mod $(wildcard pb/*/*.proto)
	PATH=$(PWD)/build/tools build/tools/protoc/bin/protoc --go_out=. --go_opt=module=$(shell go list -m) \
		--go-grpc_out=. --go-grpc_opt=module=$(shell go list -m) \
		$(wildcard pb/*/*.proto)
	@touch $@

build/tools/protoc:
	curl --progress-bar -Lo protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v25.3/protoc-25.3-$(shell uname | tr '[:upper:]' '[:lower:]')-$(shell uname -m).zip
	rm -rf $@
	mkdir -p build/tools
	unzip -qd build/tools/protoc protoc.zip
	rm protoc.zip
	chmod +x build/tools/protoc/bin/protoc

build/tools/protoc-gen-go:
	GOBIN=$(PWD)/build/tools go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28

build/tools/protoc-gen-go-grpc:
	GOBIN=$(PWD)/build/tools go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

.PHONY: pblint
pblint: build/tools/protolint
	$< pb

build/tools/protolint:
	GOBIN=$(PWD)/build/tools go install github.com/yoheimuta/protolint/cmd/protolint@v0.47.5
