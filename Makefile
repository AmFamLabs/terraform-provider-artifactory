.PHONY: build clean install

FILENAME = terraform-provider-artifactory
LOCATION = /usr/bin/$(FILENAME)

$(FILENAME):
	go build
	chmod +x $(FILENAME)

$(LOCATION): $(FILENAME)
	sudo cp $< $(LOCATION)

clean:
	sudo rm $(FILENAME) || true
	sudo rm $(LOCATION) || true

build: $(FILENAME)

install: $(LOCATION)
