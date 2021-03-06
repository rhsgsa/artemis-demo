# Shared Filesystem High-Availability Demo

Before you run the demo, you should ensure that the `artemis` container image has been created. You can do this by running `make image`.

To run the demo,

1. Setup the instances and shared directory by running

		make setup

1. Start the `node1` and `node2` containers

		make start

1. Login to the `node1` console at <http://localhost:8161/console> with `admin` / `password`

1. `node1` is the live server - the `Cluster Info` section should state that this node is a `Master`

1. Login to the `node2` console at <http://localhost:8261/console> with `admin` / `password` in an incognito window

1. `node2` is the backup server - the `Cluster Info` section should state that this node is a `Slave`

1. Send a message on the `demo` queue

1. Stop `node1`

		docker stop node1

1. `node2` should become the live server

1. Browse the `demo` queue on `node2` and ensure that the message is still there

1. Send another message on the `demo` queue from the `node2` console

1. Start `node1`

		docker start node1

1. `node1` should become the live server and `node2` should become the backup server

1. Browse the `demo` queue on the `node1` console and check that the queue contains 2 messages


## Cleaning Up

When the demo is over, clean everything up with

	make clean
