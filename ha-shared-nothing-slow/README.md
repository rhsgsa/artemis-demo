# Shared Nothing High-Availability Demo - Slow Connection

This scenario shows what happens when the live and backup nodes have a slow connection between them.

Before you run the demo, you should ensure that the `artemis` and `slow-proxy` container images have been created. You can do this by running `make image` and `make slow-image`.

To run the demo,

1. Setup the instances by running

		make setup

1. Start the `node1` and `node2` containers

		make start

1. The nodes will not be considered highly-available until `node1` has finished syncing its journal and bindings over to `node2`. When the syncing has been completed, you will see the following in the `node2`'s logs

		[org.apache.activemq.artemis.core.server] AMQ221024: Backup server ActiveMQServerImpl::name=node2 is synchronized with live server, nodeID=c3a32dda-a5d9-11ec-9c50-0242ac1e0004.
		[org.apache.activemq.artemis.core.server] AMQ221031: backup announced

1. Login to the `node1` console at <http://localhost:8161/console> with `admin` / `password`

1. `node1` is the live server

1. Login to the `node2` console at <http://localhost:8261/console> with `admin` / `password` in an incognito window

1. `node2` is the backup server

1. Send a message on the `demo` queue

1. Stop `node1`

		$(DOCKER) stop node1

1. `node2` should become the live server

1. Browse the `demo` queue on `node2` and ensure that the message is still there

1. Send another message on the `demo` queue from the `node2` console

1. Start `node1`

		docker start node1

1. `node1` should become the backup server - it doesn't have access to `acceptors` and `addresses`

1. Stop `node2`

		docker stop node2
		
1. `node1` should become the live server

1. Browse the `demo` queue on the `node1` console and check that the queue contains 2 messages


## Cleaning Up

When the demo is over, clean everything up with

	make clean
