#!/bin/bash

# Script for running the Dockerfile for Traffic Ops PostgREST
# The Dockerfile sets up a Docker image which can be used for any new Traffic Ops PostgREST container;
# This script, which should be run when the container is run (it's the ENTRYPOINT), will configure the container.
#
# The following environment variables must be set, ordinarily by `docker run -e` arguments:
# USER
# PASS
# URI - without protocol, i.e. just the Fully Qualified Domain Name and the port, e.g. example.net:5432
# DATABASE

start() {
	/postgrest postgres://${USER}:${PASS}@${URI}/${DATABASE} --port 9001 --schema public --anonymous postgres
}

init() {
	echo "INITIALIZED=1" >> /etc/environment
}

source /etc/environment
if [ -z "$INITIALIZED" ]; then init; fi
start
