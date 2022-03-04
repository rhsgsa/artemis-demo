# Artemis Demo

This demo runs on Artemis in containers.


## Pre-Requisites

1. `docker`
1. `docker-compose`
1. `make`
1. `sed`
1. `patch`


## Setup

Create the `artemis` container image by running

	make install


## Cleaning Up

Delete the `artemis` container image by running

	make clean-image

Or delete the `artemis` container image and all the temp directories by running

	make clean-all


## Resources

* [amqv7-workshop github repo](https://github.com/RedHatWorkshops/amqv7-workshop)

* [AMQ 7 Broker Supported Configurations](https://access.redhat.com/articles/2791941)
