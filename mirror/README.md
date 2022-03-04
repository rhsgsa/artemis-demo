# Disaster Recovery Demo

This demo illustrates how you can setup a broker to [mirror](https://activemq.apache.org/components/artemis/documentation/latest/amqp-broker-connections.html#mirroring) data to a second broker.

Before you run the demo, you should ensure that the `artemis` container image has been created. You can do this by running `make install` in the parent directory.

To run the demo,

1. Setup the instances by running

		make setup

1. The command above sets `node1` up with a queue named `demo`; `node1` is also configured to mirror to `node2`

1. Start the `node1` and `node2` containers

		make start

1. Send some messages to the `demo` queue on `node1`

		make producer1

1. Consume the messages on `node1`

		make consumer1

1. The command in the previous step should output 10 messages - press Ctrl-C to stop the consumer

1. Check to see if we can consume the messages on `node2`

		make consumer2

1. The command in the previous step should not output any messages - this is because they have already been consumed on `node1` and the message acknowledgements have been mirrored to `node2`

1. Press Ctrl-C to stop the consumer

1. Send some more messages to the `demo` queue on `node1`

		make producer1

1. Without consuming the messages on `node1`, simulate a disaster by stopping `node1`

		docker stop node1

1. Check if we can consume the messages on `node2`

		make consumer2

1. The command in the previous step should output 10 messages - this is because they were not consumed at `node1`

Note: You will not be able to browse the `demo` queue on `node2` on the Artemis console. A NullPointerException will be thrown when you try browsing the `demo` queue.


## Cleaning Up

When the demo is over, clean everything up with

	make clean
