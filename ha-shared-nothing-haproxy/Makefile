CWD:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

include $(CWD)/../Makefile

STORAGE=$(CWD)/demo
ARTEMIS_BASE=/usr/local/artemis

USER=admin
PASSWORD=password

.PHONY: setup clean stop start consumer1 consumer2 producer1 producer2

setup: clean
	mkdir node1 node2 node3
	# Create artemis node1 instance and fix CORS setting
	$(DOCKER) run \
	  --rm \
	  --name node1 \
	  -v ${PWD}/node1/:/home/jboss/broker/ \
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
	     --max-hops 1 \
	     --http-host 0.0.0.0 \
	     --host node1 \
	     --queues demo:anycast \
	     /home/jboss/broker/ \
	   && \
	   sed \
	     -i \
		 -e 's|<allow-origin>.*|<allow-origin>*://*</allow-origin>|' \
		 /home/jboss/broker/etc/jolokia-access.xml"

	# Create artemis node2 instance and fix CORS setting
	$(DOCKER) run \
	  --rm \
	  --name node2 \
	  -v ${PWD}/node2/:/home/jboss/broker/ \
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
	     --max-hops 1 \
	     --http-host 0.0.0.0 \
	     --host node2 \
		 --slave \
	     --queues demo:anycast \
	     /home/jboss/broker/ \
	   && \
	   sed \
	     -i \
		 -e 's|<allow-origin>.*|<allow-origin>*://*</allow-origin>|' \
		 /home/jboss/broker/etc/jolokia-access.xml"

	# Create artemis node3 instance and fix CORS setting
	$(DOCKER) run \
	  --rm \
	  --name node3 \
	  -v ${PWD}/node3/:/home/jboss/broker/ \
	  $(IMAGE_TAG):$(VERSION) \
	  /bin/bash \
	  -c \
	  "$(ARTEMIS_BASE)/bin/artemis create \
	     --allow-anonymous \
	     --name node3 \
	     --user $(USER) \
	     --password $(PASSWORD) \
	     --cluster-user $(USER) \
	     --cluster-password $(PASSWORD) \
	     --clustered \
		 --failover-on-shutdown \
	     --replicated \
	     --max-hops 1 \
	     --http-host 0.0.0.0 \
	     --host node3 \
		 --slave \
	     --queues demo:anycast \
	     /home/jboss/broker/ \
	   && \
	   sed \
	     -i \
		 -e 's|<allow-origin>.*|<allow-origin>*://*</allow-origin>|' \
		 /home/jboss/broker/etc/jolokia-access.xml"

	

	 cp node1_broker.xml ./node1/etc/broker.xml
	 cp node2_broker.xml ./node2/etc/broker.xml
	 cp node3_broker.xml ./node3/etc/broker.xml
clean: stop
	rm -rf node1 node2 node3

stop:
	-@$(DOCKER_COMPOSE) down 2>/dev/null || exit 0

start:
	@$(DOCKER_COMPOSE) up

consumer1:
	$(DOCKER) exec -it node1 \
	  /home/jboss/broker/bin/artemis consumer \
	    --url tcp://haproxy:5672 \
	    --verbose \
	    --destination queue://demo

consumer2:
	$(DOCKER) exec -it node2 \
	  /home/jboss/broker/bin/artemis consumer \
	    --url tcp://haproxy:5672 \
	    --verbose \
	    --destination queue://demo

producer1:
	$(DOCKER) exec -it node1 \
	  /home/jboss/broker/bin/artemis producer \
	    --url tcp://haproxy:5672 \
	    --message-count 10 \
	    --destination queue://demo

producer2:
	$(DOCKER) exec -it node2 \
	  /home/jboss/broker/bin/artemis producer \
	  --url tcp://haproxy:5672 \
	  --message-count 10 \
	  --destination queue://demo
