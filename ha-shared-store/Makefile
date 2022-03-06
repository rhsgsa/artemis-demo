CWD:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

include $(CWD)/../Makefile

STORAGE=$(CWD)/demo
ARTEMIS_BASE=/usr/local/artemis

USER=admin
PASSWORD=password

.PHONY: setup clean stop start

setup: clean
	mkdir -p \
	  $(STORAGE)/shared-store \
	  $(STORAGE)/instances/node1 \
	  $(STORAGE)/instances/node2
	
	chmod -R ugo+w $(STORAGE)

	docker run \
	  --rm \
	  --name node1 \
	  -v $(STORAGE)/shared-store:/storage/shared-store \
	  -v $(STORAGE)/instances/node1:/storage/instances/node1 \
	  $(IMAGE_TAG):$(VERSION) \
	  $(ARTEMIS_BASE)/bin/artemis create \
	    --allow-anonymous \
	    --name node1 \
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
	    /storage/instances/node1

	# Fix CORS setting
	docker run \
	  --rm \
	  --name node1 \
	  -v $(STORAGE)/instances/node1:/storage/instances/node1 \
	  -w /storage/instances/node1/etc \
	  $(IMAGE_TAG):$(VERSION) \
	  sed \
	    -i \
		-e 's|<allow-origin>.*|<allow-origin>*://*</allow-origin>|' \
		jolokia-access.xml

	docker run \
	  --rm \
	  --name node2 \
	  -v $(STORAGE)/shared-store:/storage/shared-store \
	  -v $(STORAGE)/instances/node2:/storage/instances/node2 \
	  $(IMAGE_TAG):$(VERSION) \
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
	    /storage/instances/node2

	# Fix CORS setting
	docker run \
	  --rm \
	  --name node2 \
	  -v $(STORAGE)/instances/node2:/storage/instances/node2 \
	  -w /storage/instances/node2/etc \
	  $(IMAGE_TAG):$(VERSION) \
	  sed \
	    -i \
		-e 's|<allow-origin>.*|<allow-origin>*://*</allow-origin>|' \
		jolokia-access.xml

clean: stop
	-@docker run \
	  --rm \
	  --name node1 \
	  -v $(STORAGE):/storage \
	  $(IMAGE_TAG):$(VERSION) \
	  /bin/bash -c 'rm -rf /storage/*'
	@rm -rf $(STORAGE)

stop:
	-@$(DOCKER_COMPOSE) down

start:
	@$(DOCKER_COMPOSE) up
