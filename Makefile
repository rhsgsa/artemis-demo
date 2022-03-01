VERSION=2.20.0
ARTEMIS_BASE=/usr/local/artemis

BASE:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

STORAGE=$(BASE)/demo

.PHONY: usage clean install

usage:
	@echo "AMQ Demo"

install: clean
	sudo rm -rf /usr/local/artemis
	curl -Lo /tmp/artemis.zip "https://www.apache.org/dyn/closer.cgi?filename=activemq/activemq-artemis/$(VERSION)/apache-artemis-$(VERSION)-bin.zip&action=download"
	cd /tmp && unzip /tmp/artemis.zip
	rm -f /tmp/artemis.zip
	sudo mv /tmp/apache-artemis-* $(ARTEMIS_BASE)

clean:
	@rm -rf $(STORAGE)

