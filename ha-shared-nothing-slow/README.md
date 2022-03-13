# Shared Nothing High-Availability Demo - Slow Connection

---

## Todo

* Static connectors

	* [docs](https://activemq.apache.org/components/artemis/documentation/1.0.0/clusters.html#discovery-using-static-connectors)

	* [example](https://github.com/apache/activemq-artemis/blob/913a87c948312aebc244a43f7dd6373b47599ec3/examples/features/clustered/clustered-static-discovery/src/main/resources/activemq/server1/broker.xml#L35)

* Docs to `make sloxy-image`

---

Before you run the demo, you should ensure that the `artemis` container image has been created. You can do this by running `make image`.

To run the demo,

1. Setup the instances by running

		make setup

1. Start the `node1` and `node2` containers

		make start

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
