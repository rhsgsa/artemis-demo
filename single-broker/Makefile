CWD:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

include $(CWD)/../Makefile

STORAGE=$(CWD)/demo
ARTEMIS_BASE=/usr/local/artemis

USER=admin
PASSWORD=password

.PHONY: setup clean stop start debug

setup: clean
	$(DOCKER) volume create artemis-node1

	# Create artemis node instance and fix CORS setting
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
	     --http-host 0.0.0.0 \
	     /storage/instances/node1 \
	   && \
	   sed \
	     -i \
		 -e 's|<allow-origin>.*|<allow-origin>*://*</allow-origin>|' \
		 /storage/instances/node1/etc/jolokia-access.xml"

clean: stop
	-@$(DOCKER) volume rm artemis-node1 2>/dev/null || exit 0

stop:
	-@$(DOCKER) stop node1 2>/dev/null || exit 0
	-@$(DOCKER) rm node1 2>/dev/null || exit 0

start:
	$(DOCKER) run \
	  --rm \
	  --name node1 \
	  -it \
	  -p 8161:8161 \
	  -p 61616:61616 \
	  -p 5445:5445 \
	  -p 5672:5672 \
	  -p 1883:1883 \
	  -p 61613:61613 \
	  --mount type=volume,source=artemis-node1,target=/storage/instances/node1 \
	  $(IMAGE_TAG):$(VERSION) \
	  /storage/instances/node1/bin/artemis run

debug:
	$(DOCKER) exec -it -w /storage/instances/node1 node1 /bin/bash
