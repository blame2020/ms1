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
	$(MAKE) GO_FLAGS="'-ldflags=-s -w -buildid= -X main.version=$(VERSION)' -trimpath" VERSION=$(VERSION)
	$(MAKE) oci

DOCKER := docker
.PHONY: docker
oci: $(BINS)
	$(DOCKER) image build -t ms1-server:$(VERSION) -t ms1-server:latest .

.PHONY: integration
integration: oci
	go test -v -count 1 ./integration

.PHONY: lint
lint: bin/golangci-lint
	./bin/golangci-lint run ./...

.PHONY: gen
gen:
	go generate ./...

.PHONY: test
test:
	go test -v ./...

bin/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(PWD)/bin v1.54.2

audit: bin/trivy
	./bin/trivy fs --scanners vuln .

bin/trivy:
	curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b $(PWD)/bin v0.45.1

.PHONY: ci
ci:
	dagger run go run ./ci

.PHONY: check-git-clean
check-git-clean:
	test $(shell git status -s | wc -l) -eq 0

.PHONY: clean
clean:
	rm -rf build bin

pbgen: protoc bin/protoc-gen-go bin/protoc-gen-go-grpc go.mod $(wildcard proto/*/*.proto)
	PATH=$(PWD)/bin ./protoc/bin/protoc --go_out=. --go_opt=module=$(shell go list -m) \
		--go-grpc_out=. --go-grpc_opt=module=$(shell go list -m) \
		$(wildcard proto/*/*.proto)
	@touch $@

protoc:
	curl -Lo protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v24.4/protoc-24.4-linux-$(shell uname -m).zip
	rm -rf $@
	unzip -d protoc protoc.zip
	rm protoc.zip
	chmod +x protoc/bin/protoc

bin/protoc-gen-go:
	GOBIN=$(PWD)/bin go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28

bin/protoc-gen-go-grpc:
	GOBIN=$(PWD)/bin go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

.PHONY: pblint
pblint: bin/protolint
	./$< proto

bin/protolint:
	GOBIN=$(PWD)/bin go install github.com/yoheimuta/protolint/cmd/protolint@v0.42.2
