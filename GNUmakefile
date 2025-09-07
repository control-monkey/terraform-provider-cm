TEST?=./...
PKGNAME?=./internal/provider
VERSION?=$(shell grep -o 'Version = \".*\"' version/version.go | grep -o \[0-9.]\\+)
RELEASE?=v$(VERSION)
ENV?=local

V := 0
Q := $(if $(filter 1,$(V)),,@)
GO := GO111MODULE=on go

default: build

.PHONY: build
build: fmtcheck
	go install

.PHONY: test
test: fmtcheck
	go test $(TEST) -timeout=30s -parallel=4

.PHONY: testacc
testacc: fmtcheck
	set -a; if [ -f .env.$(ENV) ]; then source .env.$(ENV); fi; set +a; TF_ACC=1 go test $(TEST) -v -count 1 -parallel 20 $(TESTARGS) -timeout 120m

.PHONY: testcompile # https://stackoverflow.com/questions/72721580/how-to-compile-all-tests-across-a-repo-without-executing-them
testcompile:
	go test -run=SHOULD_NEVER_MATCH $(TEST) $(TESTARGS)

.PHONY: vet
vet:
	go vet ./...

.PHONY: fmt
fmt:
	@gofmt -s -w $(CURDIR)

.PHONY: fmtcheck
fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

.PHONY: depscheck
depscheck:
	@echo "==> Checking source code with go mod tidy..."
	@go mod tidy
	@git diff --exit-code -- go.mod go.sum || \
		(echo; echo "Unexpected difference in go.mod/go.sum files. Run 'go mod tidy' command or revert any go.mod/go.sum changes and commit."; exit 1)

.PHONY: docs
docs: tools
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

.PHONY: docscheck
docscheck: docs
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs validate

.PHONY: tools
tools:
	@go generate -tags tools tools/tools.go

.PHONY: release
release: CI_JOB_NUMBER=$(shell read -p "» Last successful CI job number: " n && [[ -z "$$n" ]] && n="unknown"; echo $$n)
release:
	@git commit -a -m "chore(release): $(RELEASE)" -m "CI_JOB_NUMBER: $(CI_JOB_NUMBER)"
	@git tag -f -m    "chore(release): $(RELEASE)" $(RELEASE)
	@git push --follow-tags


# The section below is used only for local development.
# In order to build a ControlMonkey terraform provider, enter the terminal and enter 'make cm_provider'.
# Then, enter sandbox directory (using `cd sandbox`) and after `terraform init` command, you are ready to go and use
# ControlMonkey local terraform provider.
OS ?= darwin
ARCH ?= arm64
OS_ARCH := $(OS)_$(ARCH)

CM_TF_PROVIDER := terraform-provider-cm_$(RELEASE)


# Builds the go binary
.PHONY: binary
binary:
	go fmt ./...
	echo "Building Go binary"
	go build -o $(CM_TF_PROVIDER)

# Sets up your local workstation to "accept" this local provider binary
.PHONY: init
init: binary
	echo $(VERSION)
	echo "Initializing..."
	echo "Setting up for local provider..."
	rm -f ~/.terraform.d/plugins/example.com/control-monkey/cm/$(VERSION)/$(OS_ARCH)/$(CM_TF_PROVIDER)
	mkdir -p ~/.terraform.d/plugins/example.com/control-monkey/cm/$(VERSION)/$(OS_ARCH)
	ln -s $(CURDIR)/$(CM_TF_PROVIDER) ~/.terraform.d/plugins/example.com/control-monkey/cm/$(VERSION)/$(OS_ARCH)/$(CM_TF_PROVIDER)

# Builds the go binary, and cleans up Terraform lock file just in case
.PHONY: build_local
build: binary
	if [ -f "sandbox/.terraform.lock.hcl" ]; then \
	  rm sandbox/.terraform.lock.hcl; \
	fi
	if [ -f "sandbox/imports/.terraform.lock.hcl" ]; then \
	  rm sandbox/imports/.terraform.lock.hcl; \
	fi
	if [ -f "sandbox/data-sources/.terraform.lock.hcl" ]; then \
	  rm sandbox/data-sources/.terraform.lock.hcl; \
	fi
	terraform -chdir=sandbox init && terraform -chdir=sandbox/imports init && terraform -chdir=sandbox/data-sources init

# Creates ControlMonkey provider for local usage
.PHONY: cm_provider
cm_provider:
	make binary && make init && make build

# Mimicking build - run before push to vcs
.PHONY: pre_build
pre_build: fmt docscheck vet testcompile testacc

# Optimize imports
.PHONY: imports
imports:
	$(Q) goimports -w $$($(GO) list -f {{.Dir}} ./... | grep -v /vendor/)

.PHONY: testentity
testentity: ## Run a specific test (requires TESTNAME)
	@if [ -z "$(TESTNAME)" ]; then \
		echo "ERROR: TESTNAME must be set, e.g. make testentity TESTNAME=NotificationEndpointResource"; \
		exit 1; \
	fi
	@bash -lc 'set -a; if [ -f .env.$(ENV) ]; then source .env.$(ENV); fi; set +a; TF_ACC=1 go test $(PKGNAME) -v -run $(TESTNAME) $(TESTARGS) | cat'
