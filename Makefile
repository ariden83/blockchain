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
	CGO_ENABLED=0 GOOS=linux go build -mod vendor -ldflags "-X main.Version=$$GIT_TAG_NAME" -o bin/main ./tutorial/proof-work/.
	-cli_level=INFO ./bin/main

local-light:
	@echo "> Launch local p2p ..."
	go fmt ./...
	# make p2p_target=/ip4/127.0.0.1/tcp/8098/p2p/QmWV1qKRBSy8vggYgMSWDGukmwcus8wbuSoru31oNaEWdd local-light
	go run ./cmd/light/main.go -p2p_target $(p2p_target)

local-networking:
	@echo "> Launch local networking ..."
	go fmt ./...
	export GO111MODULE=on;
	CGO_ENABLED=0 GOOS=linux go build -mod vendor -ldflags "-X main.Version=$$GIT_TAG_NAME" -o bin/main ./tutorial/networking/.
	-cli_level=INFO ./bin/main

local-proof-stake:
	@echo "> Launch local proof of stake ..."
	go fmt ./...
	export GO111MODULE=on;
	CGO_ENABLED=0 GOOS=linux go build -mod vendor -ldflags "-X main.Version=$$GIT_TAG_NAME" -o bin/main ./tutorial/proof-stake/.
	-cli_level=INFO ./bin/main

local-seed:
	@echo "> Launch local seed ..."
	go fmt ./...
	export GO111MODULE=on;
	CGO_ENABLED=0 GOOS=linux go build -mod vendor -ldflags "-X main.Version=$$GIT_TAG_NAME" -o bin/main ./tutorial/seed/.
	-cli_level=INFO ./bin/main

local-relay:
	@echo "> Launch local relay ..."
	go fmt ./...
	export GO111MODULE=on;
	CGO_ENABLED=0 GOOS=linux go build -mod vendor -ldflags "-X main.Version=$$GIT_TAG_NAME" -o bin/main ./tutorial/relay/.
	-cli_level=INFO ./bin/main

local-rsa-tutorial:
	@echo "> Launch local rsa tutorial ..."
	go fmt ./...
	export GO111MODULE=on;
	CGO_ENABLED=0 GOOS=linux go build -mod vendor -ldflags "-X main.Version=$$GIT_TAG_NAME" -o bin/tutorial ./tutorial/rsa-encryption/.
	-cli_level=INFO ./bin/tutorial

local-wscat-tutorial:
	@echo "> Launch wscat tutorial ..."
	go fmt ./...
	export GO111MODULE=on;
	CGO_ENABLED=0 GOOS=linux go build -mod vendor -ldflags "-X main.Version=$$GIT_TAG_NAME" -o bin/wscat-tutorial ./tutorial/wscat/.
	-cli_level=INFO ./bin/wscat-tutorial

local-sign-script:
	@echo "> Launch local sign script ..."
	go fmt ./...
	export GO111MODULE=on;
	go test ./internal/blockchain/signSchnorr/...

local:
	@echo "> Launch local ..."
	go fmt ./...
	# gosec -tests -exclude-dir=example -exclude-dir=tutorial ./...
	export GO111MODULE=on;
	CGO_ENABLED=0 GOOS=linux go build -mod vendor -ldflags "-X main.Version=$$GIT_TAG_NAME" -o bin/main ./cmd/app/.
	-cli_level=INFO ./bin/main

local-web:
	@echo "> Launch local web  ..."
	go fmt ./...
	export GO111MODULE=on;
	CGO_ENABLED=0 GOOS=linux go build -mod vendor -ldflags "-X main.Version=$$GIT_TAG_NAME" -o bin/web ./cmd/web/.
	-cli_level=INFO ./bin/web

local-vendor:
	@echo "> Regenerate vendor ..."
	# dep init
	dep ensure -update

proto:
	@echo "> protos ..."
	docker run -v $(MAKEFILE_DIRECTORY):/app -w /app --rm jaegertracing/protobuf -I./protos api.proto --go_out=plugins=grpc:./pkg/api
	# sudo chmod u+x ./pkg/api/api.pb.go
	# docker run --rm -u $(id -u) -v${PWD}:${PWD} -w${PWD} jaegertracing/protobuf:latest --proto_path=${PWD}/api \
    # --golang_out=${PWD}/pkg/api/ --language golang /usr/include/github.com/gogo/protobuf/gogoproto/gogo.proto
	# docker run --rm -v${PWD}:${PWD} qarlm/protoc:latest --working-dir ${PWD}/protos/*.proto --grpc --language go --grpc_out ${PWD}/pkg/api


print-%:
	@echo '$($*)'

.PHONY: build lint push local-proto run test local test-local
