CWD:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

include $(CWD)/../Makefile

STORAGE=$(CWD)/demo
ARTEMIS_BASE=/usr/local/artemis

USER=admin
PASSWORD=password

.PHONY: setup clean stop start debug

setup: clean
	$(DOCKER) volume create artemis-node1

	$(DOCKER) run \
	  --rm \
	  -d \
	  --name setup \
	  --mount type=volume,source=artemis-node1,target=/storage/instances/node1 \
	  --mount type=volume,source=browser-amqp,target=/storage/shared-store/docroot \
	  $(IMAGE_TAG):$(VERSION) \
	  tail -f /dev/null

	# Create artemis node instance and fix CORS setting
	$(DOCKER) exec setup \
	  /bin/bash \
	  -c \
	  "$(ARTEMIS_BASE)/bin/artemis create \
	     --allow-anonymous \
	     --name node1 \
	     --user $(USER) \
	     --password $(PASSWORD) \
	     --http-host 0.0.0.0 \
	     --queues demo:anycast \
	     /storage/instances/node1 \
	   && \
	   sed \
	     -i \
		 -e 's|<allow-origin>.*|<allow-origin>*://*</allow-origin>|' \
		 /storage/instances/node1/etc/jolokia-access.xml"

	# Setup browser-amqp docroot
	$(DOCKER) cp $(CWD)/browser-amqp-image/docroot/ setup:/storage/shared-store/

	$(DOCKER) stop setup

clean: stop
	-@$(DOCKER) volume rm artemis-node1 2>/dev/null || exit 0
	-@$(DOCKER) volume rm browser-amqp 2>/dev/null || exit 0

stop:
	-@$(DOCKER_COMPOSE) down 2>/dev/null || exit 0

start:
	@(until nc -w 1 localhost 8080; do sleep 1; done && open "http://localhost:8080") &
	@$(DOCKER_COMPOSE) up

debug:
	$(DOCKER) exec -it -w /storage/instances/node1 node1 /bin/bash
