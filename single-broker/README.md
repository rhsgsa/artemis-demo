# Single Broker Demo

Before you run the demo, you should ensure that the `artemis` container image has been created. You can do this by running `make image`.

To run the demo,

1. Setup the `node1` instance by running

		make setup

1. Start the `node1` container

		make start

1. You should now be able to access the admin console at <http://localhost:8161/console>, logging in with `admin` / `password`

You can get a shell into the container by running

	make debug

After you're done, you can terminate the `node1` container with Ctrl-C or by running

	make stop


## Cleaning Up

When the demo is over, clean everything up with

	make clean
