PRJ_NAME := terraform-provider-artifactory

.PHONY: build install clean download install-tools generate-docs

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

download:
	@echo Download go.mod dependencies
	@go mod download

install-tools: download
	@echo Installing tools from tools.go
	@cat tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

generate-docs: install-tools
	@tfplugindocs validate
	@tfplugindocs generate