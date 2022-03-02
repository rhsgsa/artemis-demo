CWD:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

include $(CWD)/../Makefile

STORAGE=$(CWD)/demo
ARTEMIS_BASE=/usr/local/artemis

.PHONY: setup cleanstorage stop start

setup: cleanstorage
	mkdir -p $(STORAGE)

	docker run \
	  --rm \
	  --name node1 \
	  -v $(STORAGE):/storage \
	  $(IMAGE_TAG):$(VERSION) \
	  $(ARTEMIS_BASE)/bin/artemis create \
	    --allow-anonymous \
	    --name node1 \
	    --user admin \
	    --password password \
	    --cluster-user admin \
	    --cluster-password password \
	    --clustered \
	    --failover-on-shutdown \
	    --shared-store \
	    --http-host 0.0.0.0 \
	    --host node1 \
	    --data /storage/shared-store \
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
	    --user admin \
	    --password password \
	    --cluster-user admin \
	    --cluster-password password \
	    --clustered \
	    --failover-on-shutdown \
	    --shared-store \
	    --http-host 0.0.0.0 \
	    --slave \
	    --port-offset 100 \
	    --host node2 \
	    --data /storage/shared-store \
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

cleanstorage: stop
	@rm -rf $(STORAGE)

stop:
	-@docker-compose down

start:
	@docker-compose up
