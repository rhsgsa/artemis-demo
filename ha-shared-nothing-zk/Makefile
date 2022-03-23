CWD:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

include $(CWD)/../Makefile

ARTEMIS_BASE=/usr/local/artemis

USER=admin
PASSWORD=password

.PHONY: setup clean stop start

setup: clean
	$(DOCKER) volume create artemis-node1
	$(DOCKER) volume create artemis-node2

	$(DOCKER) run \
	  --rm \
	  -d \
	  --name setup \
	  --mount type=volume,source=artemis-node1,target=/storage/instances/node1 \
	  --mount type=volume,source=artemis-node2,target=/storage/instances/node2 \
	  $(IMAGE_TAG):$(VERSION) \
	  tail -f /dev/null

	# Create artemis node1 instance and fix CORS setting
	# Create artemis node2 instance and fix CORS setting
	$(DOCKER) exec setup \
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
		 /storage/instances/node1/etc/jolokia-access.xml \
	   && \
	   $(ARTEMIS_BASE)/bin/artemis create \
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

	# Patch node1 and node2 broker.xml
	$(DOCKER) cp $(CWD)/primary.patch setup:/storage/instances/node1/etc/primary.patch
	$(DOCKER) cp $(CWD)/backup.patch setup:/storage/instances/node2/etc/backup.patch
	$(DOCKER) exec \
	  setup \
	  /bin/bash \
	  -c \
	  'cd /storage/instances/node1/etc && patch broker.xml primary.patch && rm -f primary.patch && cd /storage/instances/node2/etc && patch broker.xml backup.patch && rm -f backup.patch'

	$(DOCKER) stop setup

clean: stop
	-@$(DOCKER) volume rm artemis-node1 2>/dev/null || exit 0
	-@$(DOCKER) volume rm artemis-node2 2>/dev/null || exit 0

stop:
	-@$(DOCKER_COMPOSE) down 2>/dev/null || exit 0

start:
	@$(DOCKER_COMPOSE) up
