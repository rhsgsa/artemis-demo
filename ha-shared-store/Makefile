include ../Makefile

.PHONY: setup node1 node2

setup:
	@rm -rf $(STORAGE)

	mkdir $(STORAGE)

	IP="`ip addr show ens33 | grep -o -e 'inet [^/]*' | awk '{ print $$2 }'`" \
	&& \
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
	  --host $${IP} \
	  --data $(STORAGE)/shared-store \
	  --queues demo:anycast \
	  $(STORAGE)/instances/node1

	# Fix CORS setting
	sed \
	  -i \
	  -e 's|<allow-origin>.*|<allow-origin>*://*</allow-origin>|' \
	  $(STORAGE)/instances/node1/etc/jolokia-access.xml

	IP="`ip addr show ens33 | grep -o -e 'inet [^/]*' | awk '{ print $$2 }'`" \
	&& \
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
	  --host $${IP} \
	  --data $(STORAGE)/shared-store \
	  $(STORAGE)/instances/node2

	# Fix CORS setting
	sed \
	  -i \
	  -e 's|<allow-origin>.*|<allow-origin>*://*</allow-origin>|' \
	  $(STORAGE)/instances/node2/etc/jolokia-access.xml

node1:
	$(STORAGE)/instances/node1/bin/artemis run

node2:
	$(STORAGE)/instances/node2/bin/artemis run
