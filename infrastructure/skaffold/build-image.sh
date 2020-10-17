#!/usr/bin/env bash

SERVICE=notset
WORKDIR=$(mktemp -d)
GOIMAGE=tcgo:latest
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

# Always clean up
trap "{ rm -rf $WORKDIR; }" EXIT

# TODO: Detect if the go image has changed and
#  rebuild all images if so

# Check if our base image exists
# docker image inspect $GOIMAGE > /dev/null 2>&1

# Just always build the go image for now
# if [ $? == 1 ]; then
docker build -t $GOIMAGE "$DIR/base"
# fi

# Command flags
while getopts s: flag
do
    case "${flag}" in
        s) SERVICE=${OPTARG};;
    esac
done

# Traffic ops
if [ "ops" == "$SERVICE" ]; then
  cp -r lib $WORKDIR
  cp -r traffic_ops ${WORKDIR}
  cp -r vendor $WORKDIR
  cp infrastructure/skaffold/ops/* $WORKDIR
  docker build $WORKDIR
fi