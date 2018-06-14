# Default config
VERSION   ?= $(shell echo `git describe --tag 2>/dev/null || git rev-parse --short HEAD` | sed -E 's|^v||g')
VERBOSE   ?= false

# Go configuration
GOBIN     := $(shell which go)
GOENV     ?=
GOOPT     ?=
GOLDF      = -X main.Version=$(VERSION)

# Gulp configuration
GULPBIN   := node_modules/.bin/gulp
GULPDIST  := dist/
GULPOPT   += --no-color

# Docker configuration
DOCKBIN   := $(shell which docker)
DOCKIMG   := blippar/balrog
DOCKOPTS  += --build-arg VERSION="$(VERSION)"

# If run as 'make VERBOSE=true', it will pass the '-v' option to GOBIN and will restore docker build output
ifeq ($(VERBOSE),true)
GOOPT     += -v
else
DOCKOPTS  += -q
GULPOPT   += -LL
.SILENT:
endif

# Target configuration
TARGET    := bin/balrog
GOPKGDIR   = $(@:bin/%=./cmd/%)

# Local meta targets
all: $(TARGET) $(GULPDIST)
balrog: $(TARGET)
assets: $(GULPDIST)

# Build binary with GOBIN using target name & GOPKGDIR
$(TARGET): GOOPT += -ldflags '$(GOLDF)'
$(TARGET):
	$(info >>> Building $@ from $(GOPKGDIR) using $(GOBIN))
	$(if $(GOENV),$(info >>> with $(GOENV) and GOOPT=$(GOOPT)),)
	$(GOENV) $(GOBIN) build -o $@ $(GOOPT) $(GOPKGDIR)

# Build binary staticly
static: GOLDF += -extldflags "-static"
static: GOENV += CGO_ENABLED=0 GOOS=linux
static: $(TARGET)

# Build assets using GULPBIN
$(GULPDIST):
	$(info >>> Building all assets using $(GULPBIN))
	$(GULPBIN) $(GULPOPT) build

# Run tests using GOBIN
test: GOOPT += -ldflags '$(GOLDF)'
test:
	$(info >>> Testing ./... using $(GOBIN))
	$(GOENV) $(GOBIN) test $(GOOPT) -cover ./...

# Docker
docker:
	$(info >>> Building docker image $(DOCKIMG) using $(DOCKBIN))
	$(DOCKBIN) build $(DOCKOPTS) -t $(DOCKIMG):$(VERSION) -t $(DOCKIMG):latest .

# Run linters using gometalinter
lint:
	$(info >>> Linting source code using gometalinter)
	gometalinter \
		--deadline 10m \
		--vendor \
		--sort="path" \
		--aggregate \
		--enable-gc \
		--disable-all \
		--enable goimports \
		--enable misspell \
		--enable vet \
		--enable deadcode \
		--enable varcheck \
		--enable ineffassign \
		--enable structcheck \
		--enable unconvert \
		--enable gofmt \
		./...

# Clean
clean:
	$(info >>> Cleaning up binaries and assets)
	rm -rv $(TARGET) $(GULPDIST)

.PHONY: all balrog $(TARGET) $(GULPDIST) static test lint check clean
