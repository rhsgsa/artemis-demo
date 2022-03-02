CWD:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

include $(CWD)/../Makefile

STORAGE=$(CWD)/demo
ARTEMIS_BASE=/usr/local/artemis

USER=admin
PASSWORD=password

.PHONY: setup clean stop start

setup: clean
	mkdir -p $(STORAGE)

	docker run \
	  --rm \
	  --name node1 \
	  -v $(STORAGE):/storage \
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
	    --replicated \
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
	  -v $(STORAGE):/storage \
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
	    --replicated \
	    --http-host 0.0.0.0 \
	    --host node2 \
	    --slave \
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

	patch $(STORAGE)/instances/node1/etc/broker.xml $(CWD)/check-for-live-server.patch

clean: stop
	@rm -rf $(STORAGE)

stop:
	-@docker-compose down

start:
	@docker-compose up
