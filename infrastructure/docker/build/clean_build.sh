#!/usr/bin/env sh

set -x
cp -a /trafficcontrol /tmp/. && \
	cd /tmp/trafficcontrol && \
	rm -rf dist && \
	ln -fs /trafficcontrol/dist dist &&
	((((./build/build.sh $1 2>&1; echo $? >&3) | tee ./dist/build-$1.log >&4) 3>&1) | (read x; exit $x)) 4>&1
