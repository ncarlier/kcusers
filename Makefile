.SILENT :

# App name
APPNAME=kcusers

# Go configuration
GOOS?=$(shell go env GOHOSTOS)
GOARCH?=$(shell go env GOHOSTARCH)

# Archive name
ARCHIVE=$(APPNAME)-$(GOOS)-$(GOARCH).tgz

# Extract version infos
PKG_VERSION:=github.com/ncarlier/$(APPNAME)/internal/version
VERSION:=`git describe --always --dirty`
GIT_COMMIT:=`git rev-list -1 HEAD --abbrev-commit`
BUILT:=`date`
define LDFLAGS
-X '$(PKG_VERSION).Version=$(VERSION)' \
-X '$(PKG_VERSION).GitCommit=$(GIT_COMMIT)' \
-X '$(PKG_VERSION).Built=$(BUILT)' \
-s -w -buildid=
endef

# Default task
all: build

## Clean built files
clean:
	echo ">>> Removing generated files..."
	-rm CHANGELOG.md
	-rm -rf release
.PHONY: clean

## Build executable
build:
	-mkdir -p release
	echo ">>> Building $(APPNAME) $(VERSION) for $(GOOS)-$(GOARCH) ..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -tags osusergo,netgo -ldflags "$(LDFLAGS)" -o release/$(APPNAME)
.PHONY: build

release/$(APPNAME): build

## Run tests
test: 
	go test ./...
.PHONY: test

# Check code style
check-style:
	echo ">>> Checking code style..."
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest ./...
.PHONY: check-style

# Check code criticity
check-criticity:
	echo ">>> Checking code criticity..."
	go run github.com/go-critic/go-critic/cmd/gocritic@latest check -enableAll ./...
.PHONY: check-criticity

# Check code security
check-security:
	echo ">>> Checking code security..."
	go run github.com/securego/gosec/v2/cmd/gosec@latest -quiet ./...
.PHONY: check-security

## Code quality checks
checks: check-style check-criticity
.PHONY: checks

## Install executable
install: release/$(APPNAME)
	echo "Installing $(APPNAME) to ${HOME}/.local/bin/$(APPNAME) ..."
	cp release/$(APPNAME) ${HOME}/.local/bin/$(APPNAME)
.PHONY: install

# Generate changelog
CHANGELOG.md:
	standard-changelog --first-release

## Create archive
archive: release/$(APPNAME) CHANGELOG.md
	echo "Creating release/$(ARCHIVE) archive..."
	tar czf release/$(ARCHIVE) README.md LICENSE CHANGELOG.md -C release/ $(APPNAME)
	rm release/$(APPNAME)
.PHONY: archive

## Create distribution binaries
distribution:
	GOARCH=amd64 make build archive
	GOARCH=arm64 make build archive
	GOARCH=arm make build archive
	GOOS=darwin make build archive
.PHONY: distribution