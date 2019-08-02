#!/bin/sh

pushd traffic_ops/traffic_ops_golang

if [ ! $(go get -v) ]; then
	echo "Failed to get dependencies; bailing" >&2
	exit 1
fi

if [ ! $(go test -cover -v ./...)]; then
	echo "TO tests failed" >&2
	exit 1
fi

popd
pushd lib/go-tc

if [ ! $(go test -cover -v ../../)]; then
	echo "Library tests failed" >&2
	exit 1
fi

popd
