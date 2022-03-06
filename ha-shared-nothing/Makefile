CWD:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

include $(CWD)/../Makefile

STORAGE=$(CWD)/demo
ARTEMIS_BASE=/usr/local/artemis

USER=admin
PASSWORD=password

.PHONY: setup clean stop start

setup: clean
	$(DOCKER) volume create artemis-node1
	$(DOCKER) volume create artemis-node2

	# Create artemis node1 instance and fix CORS setting
	$(DOCKER) run \
	  --rm \
	  --name node1 \
	  --mount type=volume,source=artemis-node1,target=/storage/instances/node1 \
	  $(IMAGE_TAG):$(VERSION) \
	  /bin/bash \
	  -c \
	  "$(ARTEMIS_BASE)/bin/artemis create \
	     --allow-anonymous \
	     --name node1 \
	     --user $(USER) \
	     --password $(PASSWORD) \
	     --cluster-user $(USER) \
	     --cluster-password $(PASSWORD) \
	     --clustered \
	     --failover-on-shutdown \
	     --replicated \
	     --http-host 0.0.0.0 \
	     --host node1 \
	     --queues demo:anycast \
	     /storage/instances/node1 \
	   && \
	   sed \
	     -i \
		 -e 's|<allow-origin>.*|<allow-origin>*://*</allow-origin>|' \
		 /storage/instances/node1/etc/jolokia-access.xml"

	# Create artemis node2 instance and fix CORS setting
	$(DOCKER) run \
	  --rm \
	  --name node2 \
	  --mount type=volume,source=artemis-node2,target=/storage/instances/node2 \
	  $(IMAGE_TAG):$(VERSION) \
	  /bin/bash \
	  -c \
	  "$(ARTEMIS_BASE)/bin/artemis create \
	     --allow-anonymous \
	     --name node2 \
	     --user $(USER) \
	     --password $(PASSWORD) \
	     --cluster-user $(USER) \
	     --cluster-password $(PASSWORD) \
	     --clustered \
	     --failover-on-shutdown \
	     --replicated \
	     --http-host 0.0.0.0 \
	     --host node2 \
	     --slave \
	     /storage/instances/node2 \
	   && \
	   sed \
	     -i \
		 -e 's|<allow-origin>.*|<allow-origin>*://*</allow-origin>|' \
		 /storage/instances/node2/etc/jolokia-access.xml"

	# Patch node1 broker.xml
	$(DOCKER) run \
	  --rm \
	  -d \
	  --name node1 \
	  --mount type=volume,source=artemis-node1,target=/storage/instances/node1 \
	  $(IMAGE_TAG):$(VERSION) \
	  tail -f /dev/null
	$(DOCKER) cp $(CWD)/check-for-live-server.patch node1:/storage/instances/node1/etc/check-for-live-server.patch
	$(DOCKER) stop node1
	$(DOCKER) run \
	  --rm \
	  --name node1 \
	  --mount type=volume,source=artemis-node1,target=/storage/instances/node1 \
	  -w /storage/instances/node1/etc \
	  $(IMAGE_TAG):$(VERSION) \
	  /bin/bash \
	  -c \
	  'patch broker.xml check-for-live-server.patch && rm -f check-for-live-server.patch'

clean: stop
	-@$(DOCKER) volume rm artemis-node1 2>/dev/null || exit 0
	-@$(DOCKER) volume rm artemis-node2 2>/dev/null || exit 0

stop:
	-@$(DOCKER_COMPOSE) down 2>/dev/null || exit 0

start:
	@$(DOCKER_COMPOSE) up
