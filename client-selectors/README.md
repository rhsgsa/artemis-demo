# Message Selector Demo

This demo shows the effect of message selectors on producers and consumers.

Before you run the demo, you should ensure that the `artemis` container image has been created. You can do this by running `make image`.

To run the demo,

1. Setup the `node1` instance and the nginx docroot by running

		make setup

1. Start the containers

		make start

1. A browser should not open to the AMQP browser client - click the Start button and experiment with different settings for the producer and consumers

After you're done, you can terminate the containers with Ctrl-C and by running

	make stop


## Cleaning Up

When the demo is over, clean everything up with

	make clean
