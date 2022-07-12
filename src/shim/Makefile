# Go command to use for build
GO ?= go
INSTALL ?= install

# Root directory of the project (absolute path).
ROOT_DIR = $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

# Base path used to install.
DESTDIR ?= /usr/local

ifneq ($(GO111MODULE),off)
    MOD_VENDOR := "-mod=vendor"
endif

golang_version_raw=$(shell $(GO) version 2>/dev/null)
ifeq (,$(GOPATH))
    golang_version_raw=
endif

# Get ARCH.
ifeq ($(ARCH),)
    ifneq (,$(golang_version_raw))
        override ARCH = $(shell $(GO) env GOARCH)
    else
        override ARCH = $(shell uname -m)
    endif
endif

ifeq ($(ARCH),x86_64)
    override ARCH = amd64
endif

ifeq ($(ARCH),aarch64)
    override ARCH = arm64
endif

# Project binaries.
SHIM = containerd-shim-rune-v2
COMMANDS += $(SHIM)

GO_BUILD_FLAGS :=
SHIM_CGO_ENABLED ?= 0

BINARIES=$(addprefix bin/,$(COMMANDS))

.PHONY: clean all binaries help install uninstall package
.DEFAULT: default

all: binaries

# Build a binary from a cmd.
bin/containerd-shim-rune-v2:
	@echo "bin/$(SHIM)"
	@CGO_ENABLED=$(SHIM_CGO_ENABLED) $(GO) build $(MOD_VENDOR) $(GO_BUILD_FLAGS) -o bin/$(SHIM) ./cmd/$(SHIM)

binaries: clean $(BINARIES) ## build binaries

clean: ## clean up binaries
	@rm -f $(BINARIES)

install: ## install binaries
	@echo "$@ $(BINARIES)"
	@$(INSTALL) -d $(DESTDIR)/bin
	@$(INSTALL) -D -m0755 $(BINARIES) $(DESTDIR)/bin

uninstall:
	@echo "$@"
	@rm -f $(addprefix $(DESTDIR)/bin/,$(notdir $(BINARIES)))

help: ## this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort
