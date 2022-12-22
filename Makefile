PACKAGE = violet
BUILD_DIR = $(PWD)/dist
# The specific binaries for each supported platform
LINUX_BIN = $(BUILD_DIR)/$(PACKAGE)_linux_amd64_v1/$(PACKAGE)
WINDOWS_BIN = $(BUILD_DIR)/$(PACKAGE)_windows_amd64/$(PACKAGE).exe
MAC_BIN = $(BUILD_DIR)/$(PACKAGE)_darwin_amd64_v1/$(PACKAGE)

ORANGE = \e[0;33m
GREEN = \e[0;32m
END = \e[0m

.PHONY: test clean

define help_text
Targets:
  - make build          Build Violet for all targets (Linux, Windows, Mac, 64-bit)
  - make run            Build and run Violet for current host
  - make test           Run all Go tests
  - make clean          Delete built artifacts
  - make [help]         Print this help.
endef
export help_text

help:
	@echo "$$help_text"

build:
	@goreleaser build --snapshot --rm-dist
	@echo -e "${GREEN}Build is complete!${END}"

clean:
	@echo -e "${ORANGE}Cleaning...${END}"
	@rm -rf $(BUILD_DIR)

test:
	@echo -e "${ORANGE}Testing...${END}"
	@go test $(PWD)/test/...

ifeq ($(OS)),Windows_NT)
BIN = $(WINDOWS_BIN)
else
UNAME := $(shell uname -s)
ifeq ($(UNAME),Linux)
BIN = $(LINUX_BIN)
endif
ifeq ($(UNAME),Darwin)
BIN = $(MAC_BIN)
endif
endif

$(BIN): build

run: $(BIN)
	@exec $?
