MAKE_TOOLS_VERSION := 6.7.1
MAKE_TOOLS_BOOT    := https://artifacts.amfamlabs.com/dl/make-tools/bootstrap-$(MAKE_TOOLS_VERSION)
MAKE_TOOLS_DIR     := $(shell pwd)/.make-tools-$(MAKE_TOOLS_VERSION)
MAKE_TOOLS_BIN     := $(MAKE_TOOLS_DIR)/bin
MAKE_TOOLS         := $(shell test -d $(MAKE_TOOLS_DIR) \
                          && echo 'already installed' \
                          || { curl -s '$(MAKE_TOOLS_BOOT)' | bash ; \
                                 echo 'installed make tools $(MAKE_TOOLS_VERSION)'; })

PRJ_NAME           := terraform-provider-artifactory

include $(MAKE_TOOLS_DIR)/include/coding.mk
include $(MAKE_TOOLS_DIR)/include/aws.mk

.PHONY: build build-ci clean install publish

BINARY_NAME = $(PRJ_NAME)
BINARY_LOCATION = /usr/bin/$(BINARY_NAME)

$(BINARY_NAME):
	go build
	chmod +x $(BINARY_NAME)

$(BINARY_LOCATION): $(BINARY_NAME)
	sudo cp $< $(BINARY_LOCATION)

clean:
	sudo rm $(BINARY_NAME) || true
	sudo rm $(BINARY_LOCATION) || true

build: $(BINARY_NAME)

install: $(BINARY_LOCATION)

publish:
	$(MAKE_TOOLS_BIN)/pub-artifact \
		-s dist/${BINARY_NAME} \
		-d dl-local/terraform-providers/${BINARY_NAME}-${PRJ_VERSION}
