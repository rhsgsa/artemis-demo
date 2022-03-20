# Shared Nothing High-Availability Demo - Slow Connection

This scenario shows what happens when the live and backup nodes have a slow connection between them.

	                 ┌────────┐
	                 │        │
	     ┌───────────► slow2  │
	     │           │ proxy  │
	     │           │        │
	     │           └────┬───┘
	     │                │
	     │                │
	┌────┴───┐       ┌────▼───┐
	│        │       │        │
	│ node1  │       │ node2  │
	│ broker │       │ broker │
	│        │       │        │
	└────▲───┘       └────┬───┘
	     │                │
	     │                │
	     │                │
	┌────┴───┐            │
	│        │            │
	│ slow1  ◄────────────┘
	│ proxy  │
	│        │
	└────────┘

`node1` is the live server and `node2` is the backup server.

Before you run the demo, you should ensure that the `artemis`, and `slow-proxy` container images have been created. You can do this by running

	make image

	make slow-image

To run the demo,

1. Setup the instances by running

		make setup

1. Start the `node1` and `node2` containers

		make start

1. The nodes will not be considered highly-available until `node1` has finished syncing its journal and bindings over to `node2`. When the syncing has been completed, you will see the following in the `node2`'s logs

		[org.apache.activemq.artemis.core.server] AMQ221024: Backup server ActiveMQServerImpl::name=node2 is synchronized with live server, nodeID=c3a32dda-a5d9-11ec-9c50-0242ac1e0004.
		[org.apache.activemq.artemis.core.server] AMQ221031: backup announced

1. Open a browser to the browser-based AMQP producer

		make producer

1. Send a message by typing something in the text box and pressing enter - you should see an acceptance and settlement message almost immediately

1. Simulate a slow connection from `node2` to `node1`

		make slow-node1

1. Wait a few seconds for the command above to take effect - look for the following message in the `slow1` container logs

		creating new read buffer of size 1 bytes

1. Send another message with the browser-based AMQP producer - the acceptance and settlement messages should now appear after about 20 seconds

1. Switch `node2`'s connection to `node1` back to normal

		make fast-node1


## Cleaning Up

When the demo is over, clean everything up with

	make clean
