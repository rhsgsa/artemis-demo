VERSION=2.20.0
IMAGE_TAG=ghcr.io/kwkoo/artemis
DOCKER_COMPOSE=docker-compose

BASE:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

.PHONY: usage clean-image image clean-all

usage:
	@echo "AMQ Demo - 'make image' to create the artemis container image, 'make clean' to delete the artemis container image, 'make clean-all' to remove the container image and all data directories"

image:
	@docker build \
	  -t $(IMAGE_TAG):$(VERSION) \
	  --build-arg VERSION=$(VERSION) \
	  $(BASE)/image

clean-image:
	@docker rmi -f $(IMAGE_TAG):$(VERSION)

clean-all: clean-image
	@for dir in clustering ha-shared-nothing ha-shared-store; do \
	  cd $(BASE)/$$dir; \
	  make clean; \
	done
