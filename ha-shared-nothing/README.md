# Shared Nothing High-Availability Demo

Before you run the demo, you should ensure that `artemis` is installed by running `make install` in the parent directory.

To run the demo,

1. Setup the instances by running

		make setup

1. Start up `node1`

		make node1

1. Open another terminal and start `node2`

		make node2

1. Login to the `node1` console at <http://localhost:8161/console> with `admin` / `password`

1. Login to the `node2` console at <http://localhost:8261/console> with `admin` / `password` in an incognito window

1. Send a message on the `demo` queue

1. Stop `node1`

1. `node2` should become active

1. Browse the `demo` queue on `node2` and ensure that the message is still there
