CWD:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

include $(CWD)/../Makefile

ARTEMIS_BASE=/usr/local/artemis

USER=admin
PASSWORD=password

.PHONY: setup clean stop start consumer1 consumer2 producer1 producer2

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
	     --max-hops 1 \
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
	     --max-hops 1 \
	     --http-host 0.0.0.0 \
	     --host node2 \
	     --queues demo:anycast \
	     /storage/instances/node2 \
	   && \
	   sed \
	     -i \
		 -e 's|<allow-origin>.*|<allow-origin>*://*</allow-origin>|' \
		 /storage/instances/node2/etc/jolokia-access.xml"

clean: stop
	-@$(DOCKER) volume rm artemis-node1 2>/dev/null || exit 0
	-@$(DOCKER) volume rm artemis-node2 2>/dev/null || exit 0

stop:
	-@$(DOCKER_COMPOSE) down 2>/dev/null || exit 0

start:
	@$(DOCKER_COMPOSE) up

consumer1:
	$(DOCKER) exec -it node1 \
	  /storage/instances/node1/bin/artemis consumer \
	    --url tcp://node1:61616 \
	    --verbose \
	    --destination queue://demo

consumer2:
	$(DOCKER) exec -it node2 \
	  /storage/instances/node2/bin/artemis consumer \
	    --url tcp://node2:61616 \
	    --verbose \
	    --destination queue://demo

producer1:
	$(DOCKER) exec -it node1 \
	  /storage/instances/node1/bin/artemis producer \
	    --url tcp://node1:61616 \
	    --message-count 10 \
	    --destination queue://demo

producer2:
	$(DOCKER) exec -it node2 \
	  /storage/instances/node2/bin/artemis producer \
	  --url tcp://node2:61616 \
	  --message-count 10 \
	  --destination queue://demo
