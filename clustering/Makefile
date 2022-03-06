CWD:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

include $(CWD)/../Makefile

STORAGE=$(CWD)/demo
ARTEMIS_BASE=/usr/local/artemis

USER=admin
PASSWORD=password

.PHONY: setup clean stop start consumer1 consumer2 producer1 producer2

setup: clean
	mkdir -p \
	  $(STORAGE)/instances/node1 \
	  $(STORAGE)/instances/node2
	
	chmod -R ugo+w $(STORAGE)

	$(DOCKER) run \
	  --rm \
	  --name node1 \
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
	    --max-hops 1 \
	    --http-host 0.0.0.0 \
	    --host node1 \
	    --queues demo:anycast \
	    /storage/instances/node1

	# Fix CORS setting
	$(DOCKER) run \
	  --rm \
	  --name node1 \
	  -v $(STORAGE)/instances/node1:/storage/instances/node1 \
	  -w /storage/instances/node1/etc \
	  $(IMAGE_TAG):$(VERSION) \
	  sed \
	    -i \
		-e 's|<allow-origin>.*|<allow-origin>*://*</allow-origin>|' \
		jolokia-access.xml

	$(DOCKER) run \
	  --rm \
	  --name node2 \
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
	    --max-hops 1 \
	    --http-host 0.0.0.0 \
	    --host node2 \
	    --queues demo:anycast \
	    /storage/instances/node2

	# Fix CORS setting
	$(DOCKER) run \
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
	-@$(DOCKER) run \
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
