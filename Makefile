NAME := randstr
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
DIST := dist

PLATFORMS := \
	linux/amd64 \
	linux/arm64 \
	darwin/amd64 \
	darwin/arm64 \
	windows/amd64 \
	windows/arm64

.PHONY: build clean test cross-build

build:
	go build -o $(NAME)

test:
	go test -v ./...

cross-build: clean
	@mkdir -p $(DIST)
	@$(foreach platform,$(PLATFORMS), \
		$(eval GOOS := $(word 1,$(subst /, ,$(platform)))) \
		$(eval GOARCH := $(word 2,$(subst /, ,$(platform)))) \
		$(eval EXT := $(if $(filter windows,$(GOOS)),.exe,)) \
		$(eval OUT := $(DIST)/$(NAME)-$(GOOS)-$(GOARCH)$(EXT)) \
		echo "Building $(OUT)" && \
		GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-s -w" -o $(OUT) . && \
	) true

clean:
	rm -rf $(DIST)
