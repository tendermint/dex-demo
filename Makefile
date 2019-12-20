#!/usr/bin/make

VERSION := $(shell echo $(shell git describe) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
LEDGER_ENABLED ?= true

export GO111MODULE = on

# process build tags

build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

ifeq ($(WITH_CLEVELDB),yes)
  build_tags += gcc
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=dex \
		  -X github.com/cosmos/cosmos-sdk/version.ServerName=dexd \
		  -X github.com/cosmos/cosmos-sdk/version.ClientName=dexcli \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)"

ifeq ($(WITH_CLEVELDB),yes)
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

all: dexd dexcli

all-cross:
	GOOS=linux GOARCH=amd64 go build -mod=readonly $(BUILD_FLAGS) -o build/dexd-linux-x64 cmd/dexd/main.go
	GOOS=linux GOARCH=amd64 go build -mod=readonly $(BUILD_FLAGS) -o build/dexcli-linux-x64 cmd/dexcli/main.go
	GOOS=darwin GOARCH=amd64 go build -mod=readonly $(BUILD_FLAGS) -o build/dexd-darwin-x64 cmd/dexd/main.go
	GOOS=darwin GOARCH=amd64 go build -mod=readonly $(BUILD_FLAGS) -o build/dexcli-darwin-x64 cmd/dexcli/main.go

abigen:
	./update-contracts.sh

update-ui:
	./update-ui.sh
	packr -z

update-testnet-hardfork:
	ssh testnet "bash -s" < ./update-testnet-hardfork.sh

update-testnet-nofork:
	ssh testnet "bash -s" < ./update-testnet-nofork.sh

build:
	mkdir -p build/

build/dexd: build
	go build -mod=readonly $(BUILD_FLAGS) -o $@ ./cmd/dexd/

dexcli: build/dexcli

dexd: build/dexd

build/dexcli: build
	go build -mod=readonly $(BUILD_FLAGS) -o $@ ./cmd/dexcli/

clean:
	rm -rf ./build

install:
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/dexd/
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/dexcli/

test: all
	UEX_TEST_BIN_DIR=$(shell pwd)/build go test -v ./...

test-unit:
	go test -v ./...

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./x/embedded/ui/a_ui-packr.go" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./x/embedded/ui/a_ui-packr.go" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./x/embedded/ui/a_ui-packr.go" | xargs goimports -w -local github.com/cosmos/cosmos-sdk

.PHONY: test all clean install dexd dexcli
