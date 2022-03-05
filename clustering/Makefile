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

	docker run \
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
	cd $(STORAGE)/instances/node1/etc/ \
	&& \
	cat jolokia-access.xml \
	| \
	sed -e 's|<allow-origin>.*|<allow-origin>*://*</allow-origin>|' \
	> jolokia-access.xml.new \
	&& \
	mv jolokia-access.xml.new jolokia-access.xml

	docker run \
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
	cd $(STORAGE)/instances/node2/etc/ \
	&& \
	cat jolokia-access.xml \
	| \
	sed -e 's|<allow-origin>.*|<allow-origin>*://*</allow-origin>|' \
	> jolokia-access.xml.new \
	&& \
	mv jolokia-access.xml.new jolokia-access.xml

clean: stop
	@rm -rf $(STORAGE)

stop:
	-@docker compose down

start:
	@docker compose up

consumer1:
	docker exec -it node1 \
	  /storage/instances/node1/bin/artemis consumer \
	    --url tcp://node1:61616 \
	    --verbose \
	    --destination queue://demo

consumer2:
	docker exec -it node2 \
	  /storage/instances/node2/bin/artemis consumer \
	    --url tcp://node2:61616 \
	    --verbose \
	    --destination queue://demo

producer1:
	docker exec -it node1 \
	  /storage/instances/node1/bin/artemis producer \
	    --url tcp://node1:61616 \
	    --message-count 10 \
	    --destination queue://demo

producer2:
	docker exec -it node2 \
	  /storage/instances/node2/bin/artemis producer \
	  --url tcp://node2:61616 \
	  --message-count 10 \
	  --destination queue://demo
