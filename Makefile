# Go configuration
GOBIN     := $(shell which go)
GOENV     ?=
GOOPT     ?=
GOLDF     ?=

# Logging configuration
VERBOSE   ?= false

# If run as 'make VRBOSE=true', it will pass the '-v' option to GOBIN and will restore docker build output
ifeq ($(VERBOSE),true)
GOOPT     += -v
endif

# Target configuration
TARGET    := bin/apk
GOPKGDIR   = $(@:bin/%=./cmd/%)

all: $(TARGET)

# Build binary with GOBIN using target name & GOPKGDIR
$(TARGET): GOOPT += -ldflags '$(GOLDF)'
$(TARGET):
	$(info >>> Building $@ from $(GOPKGDIR) using $(GOBIN))
	$(if $(GOENV),$(info >>> with $(GOENV) and GOOPT=$(GOOPT)),)
	@$(GOENV) $(GOBIN) build -o $@ $(GOOPT) $(GOPKGDIR)

# Build binary staticly
static: GOLDF += -extldflags "-static"
static: GOENV += CGO_ENABLED=0 GOOS=linux
static: $(TARGET)

# Run tests using GOBIN
test: GOPKGLIST = $(shell $(GOBIN) list ./... | grep -v vendor)
test: GOOPT += -ldflags '$(GOLDF)'
test:
	$(info >>> Testing ./... using $(GOBIN))
	@$(GOENV) $(GOBIN) test $(GOOPT) -cover $(GOPKGLIST)

# Run linters using gometalinter
lint:
	$(info >>> Linting ./... using gometalinter)
	@gometalinter --vendor  --disable-all \
				--enable=vet \
				--enable=vetshadow \
				--enable=golint \
				--enable=ineffassign \
				--enable=goconst \
				--enable=interfacer \
				--enable=goconst \
				--enable=unparam \
				--enable=gofmt \
				--enable=goimports \
				--enable=misspell \
				-- ./...

# Run megacheck on codebase
check:
	$(info >>> Checking ./... using megacheck)
	@megacheck ./...

# Clean
clean:
	$(info >>> Cleaning up binaries and distribuables)
	@rm -rv $(TARGET)

.PHONY: all $(TARGET) static test lint check clean
