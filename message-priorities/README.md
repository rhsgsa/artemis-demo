# Message Priorities Demo

This demo shows the effect of message priorities on producers and consumers.

Before you run the demo, you should ensure that the `artemis` container image has been created. You can do this by running `make image`.

To run the demo,

1. Setup the `node1` instance and the nginx docroot by running

		make setup

1. Start the containers

		make start

1. A browser should now open to the AMQP browser client - click the Start button

	* Set `Priority` to `0` (the lowest priority) and click on `Send Message` twice - this will send 2 messages with priority `0`
	* Set `Priority` to `5` and click on `Send Message` once
	* Set `Priority` back to `0` and click on `Send Message` twice
	* Click on `Start Consumer`
	* You should see `message 2` appear before all the other messages because it has a higher priority

After you're done, you can terminate the containers with Ctrl-C and by running

	make stop


## Cleaning Up

When the demo is over, clean everything up with

	make clean
