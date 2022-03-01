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
	  --max-hops 1 \
	  --http-host 0.0.0.0 \
	  --host $${IP} \
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
	  --max-hops 1 \
	  --http-host 0.0.0.0 \
	  --host $${IP} \
	  --queues demo:anycast \
	  --port-offset 100 \
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

consumer1:
	IP="`ip addr show ens33 | grep -o -e 'inet [^/]*' | awk '{ print $$2 }'`" \
	&& \
	$(STORAGE)/instances/node1/bin/artemis consumer \
	  --url tcp://$${IP}:61616 \
	  --verbose \
	  --destination queue://demo

consumer2:
	IP="`ip addr show ens33 | grep -o -e 'inet [^/]*' | awk '{ print $$2 }'`" \
	&& \
	$(STORAGE)/instances/node1/bin/artemis consumer \
	  --url tcp://$${IP}:61716 \
	  --verbose \
	  --destination queue://demo

producer1:
	IP="`ip addr show ens33 | grep -o -e 'inet [^/]*' | awk '{ print $$2 }'`" \
	&& \
	$(STORAGE)/instances/node1/bin/artemis producer \
	  --url tcp://$${IP}:61616 \
	  --message-count 10 \
	  --destination queue://demo

producer2:
	IP="`ip addr show ens33 | grep -o -e 'inet [^/]*' | awk '{ print $$2 }'`" \
	&& \
	$(STORAGE)/instances/node1/bin/artemis producer \
	  --url tcp://$${IP}:61716 \
	  --message-count 10 \
	  --destination queue://demo
