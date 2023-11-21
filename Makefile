# SOURCES=go.mod go.sum $(shell find . -type f -name '*.go')
SOURCES=$(shell find . -type f -name '*.go')
VERSION=$(shell hack/version.sh)

# # # # # # # # # # # # # # # # # # # #
# Go commands                         #
# # # # # # # # # # # # # # # # # # # #
BUILD_FLAGS=-ldflags "-X main.Version=${VERSION}"
GO_BUILD=go build ${BUILD_FLAGS}
GO_INSTALL=go install ${BUILD_FLAGS}
GO_FMT=go fmt
GO_TEST=go test

ifdef VERBOSE
	GO_BUILD += -v
	GO_INSTALL += -v
	GO_FMT += -x
	GO_TEST += -test.v

	RM += --verbose
endif

$(info using tag '${VERSION}')

# # # # # # # # # # # # # # # # # # # #
# Show help                           #
# # # # # # # # # # # # # # # # # # # #

.PHONY: help

help:
	@echo "Usage: make <target>"
	@echo "  build"
	@echo "  install"
	@echo "  clean"

# # # # # # # # # # # # # # # # # # # #
# Build / Install                     #
# # # # # # # # # # # # # # # # # # # #

.PHONY: build install

build: bin/meetup

bin/meetup: ${SOURCES}
	${GO_BUILD} -o $@ .

install: bin/meetup
	${GO_INSTALL} .


# # # # # # # # # # # # # # # # # # # #
# Test                                #
# # # # # # # # # # # # # # # # # # # #

.PHONY: test

test:
	${GO_TEST} ./...

# # # # # # # # # # # # # # # # # # # #
# Clean                               #
# # # # # # # # # # # # # # # # # # # #

.PHONY: clean

clean:
	${RM} bin