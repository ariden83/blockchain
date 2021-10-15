REPOSITORY ?= $(shell git config --get remote.origin.url| cut -d':' -f2 |rev |cut -c5-|rev)
REGISTRY   ?= localhost:5000
GIT_TAG_NAME ?= $(shell git describe --abbrev=1 --tags 2> /dev/null || git describe --always)

IMAGE      ?= $(REGISTRY)/ariden83/blockchain:$(GIT_TAG_NAME)

DATE = $(shell date +'%Y%m%d%H%M%S')

GOBUILDER_IMAGE ?= "golang:1.16"
BRANCH_NAME  ?= $(shell git rev-parse --abbrev-ref HEAD)
PROJECT_ROOT     ?= /go/src/ariden83/blockchain

BUILD_LABELS += --label application_branch=$(BRANCH_NAME)

BUILD_OPTIONS  = -t $(IMAGE)
BUILD_OPTIONS += $(BUILD_LABELS)
BUILD_OPTIONS += --build-arg cache=$(DATE)
BUILD_OPTIONS += --build-arg PROJECT_ROOT=$(PROJECT_ROOT)
BUILD_OPTIONS += --build-arg GOBUILDER_IMAGE=$(GOBUILDER_IMAGE)

MAKEFILE_DIRECTORY = $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))


default: build-app

build: build-app

test: initialize

initialize:
	@echo "> initialize..."

clean:
	@echo "> start clean..."

build-app:
	@echo "> start building..."
	docker build $(BUILD_OPTIONS) .

push:
	@echo "> start push..."
	docker push $(IMAGE)

run: build
	@echo "> launch local docker image"
	docker run -p 8080/tcp -p 8082:8082/tcp -p 8081:8081/tcp --rm $(IMAGE)

local-proof:
	@echo "> Launch local proof of work ..."
	go fmt ./...
	export GO111MODULE=on;
	CGO_ENABLED=0 GOOS=linux go build -mod vendor -ldflags "-X main.Version=$$GIT_TAG_NAME" -o bin/main ./cmd/proof-work/.
	-cli_level=INFO ./bin/main

local-p2p:
	@echo "> Launch local p2p ..."
	go fmt ./...
	go run ./cmd/p2p/main.go -l 10000 -secio

local-networking:
	@echo "> Launch local networking ..."
	go fmt ./...
	export GO111MODULE=on;
	CGO_ENABLED=0 GOOS=linux go build -mod vendor -ldflags "-X main.Version=$$GIT_TAG_NAME" -o bin/main ./cmd/networking/.
	-cli_level=INFO ./bin/main

local-proof-stake:
	@echo "> Launch local proof of stake ..."
	go fmt ./...
	export GO111MODULE=on;
	CGO_ENABLED=0 GOOS=linux go build -mod vendor -ldflags "-X main.Version=$$GIT_TAG_NAME" -o bin/main ./cmd/proof-stake/.
	-cli_level=INFO ./bin/main

local-seed:
	@echo "> Launch local seed ..."
	go fmt ./...
	export GO111MODULE=on;
	CGO_ENABLED=0 GOOS=linux go build -mod vendor -ldflags "-X main.Version=$$GIT_TAG_NAME" -o bin/main ./cmd/seed/.
	-cli_level=INFO ./bin/main

local:
	@echo "> Launch local ..."
	go fmt ./...
	export GO111MODULE=on;
	CGO_ENABLED=0 GOOS=linux go build -mod vendor -ldflags "-X main.Version=$$GIT_TAG_NAME" -o bin/main ./cmd/app/.
	-cli_level=INFO ./bin/main

local-pod2:
	@echo "> Launch local ..."
	go fmt ./...
	export GO111MODULE=on;
	CGO_ENABLED=0 GOOS=linux go build -mod vendor -ldflags "-X main.Version=$$GIT_TAG_NAME" -o bin/main ./cmd/app/.
	-cli_level=INFO  ./bin/main


local-vendor:
	@echo "> Regenerate vendor ..."
	# dep init
	dep ensure -update

print-%:
	@echo '$($*)'

.PHONY: build lint push local-proto run test local test-local
