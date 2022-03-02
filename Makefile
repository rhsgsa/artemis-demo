VERSION=2.20.0
IMAGE_TAG=ghcr.io/kwkoo/artemis

BASE:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

.PHONY: usage clean image

usage:
	@echo "AMQ Demo - 'make image' to create the artemis container image, 'make clean' to delete the artemis container image"

image:
	@docker build \
	  -t $(IMAGE_TAG):$(VERSION) \
	  --build-arg VERSION=$(VERSION) \
	  $(BASE)/image

clean:
	@docker rmi -f $(IMAGE_TAG):$(VERSION)
