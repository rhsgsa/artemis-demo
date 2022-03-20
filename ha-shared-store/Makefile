CWD:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

include $(CWD)/../Makefile

ARTEMIS_BASE=/usr/local/artemis

USER=admin
PASSWORD=password

.PHONY: setup clean stop start

setup: clean
	$(DOCKER) volume create artemis-node1
	$(DOCKER) volume create artemis-node2
	$(DOCKER) volume create artemis-shared

	# Create artemis node1 instance and fix CORS setting
	# Create artemis node2 instance and fix CORS setting
	$(DOCKER) run \
	  --rm \
	  --name node1 \
	  --mount type=volume,source=artemis-node1,target=/storage/instances/node1 \
	  --mount type=volume,source=artemis-node2,target=/storage/instances/node2 \
	  --mount type=volume,source=artemis-shared,target=/storage/shared-store \
	  $(IMAGE_TAG):$(VERSION) \
	  /bin/bash \
	  -c \
	  "$(ARTEMIS_BASE)/bin/artemis create \
	     --allow-anonymous \
	     --name setup \
	     --user $(USER) \
	     --password $(PASSWORD) \
	     --cluster-user $(USER) \
	     --cluster-password $(PASSWORD) \
	     --clustered \
	     --failover-on-shutdown \
	     --shared-store \
	     --http-host 0.0.0.0 \
	     --host node1 \
	     --data /storage/shared-store \
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
	     --shared-store \
	     --http-host 0.0.0.0 \
	     --slave \
	     --host node2 \
	     --data /storage/shared-store \
	     /storage/instances/node2 \
	   && \
	   sed \
	     -i \
		 -e 's|<allow-origin>.*|<allow-origin>*://*</allow-origin>|' \
		 /storage/instances/node2/etc/jolokia-access.xml"

clean: stop
	-@$(DOCKER) volume rm artemis-node1 2>/dev/null || exit 0
	-@$(DOCKER) volume rm artemis-node2 2>/dev/null || exit 0
	-@$(DOCKER) volume rm artemis-shared 2>/dev/null || exit 0

stop:
	-@$(DOCKER_COMPOSE) down 2>/dev/null || exit 0

start:
	@$(DOCKER_COMPOSE) up
