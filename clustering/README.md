# Shared Nothing High-Availability Demo

Before you run the demo, you should ensure that the `artemis` container image has been created. You can do this by running `make image`.

To run the demo,

1. Setup the instances by running

		make setup

1. Start the `node1` and `node2` containers

		make start

1. Open a new tab and start a consumer on `node1`

		make consumer1

1. Split the window horizontally (Shift-Cmd-D) and start a consumer on `node2`

		make consumer2

1. Split the window horizontally (Shift-Cmd-D) and start a producer on `node1`

		make producer1

1. The producer should send 10 messages - 5 messages should be sent to `node1` and 5 messages should be sent to `node2`

1. The behavior should be the same if you start a producer on `node2`

		make producer2

1. If you stop the consumer on `node2` and start a producer, all messages should go to `node1`


## Cleaning Up

When the demo is over, clean everything up with

	make clean
