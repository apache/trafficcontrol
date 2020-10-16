#!/usr/bin/env bash

IMAGE=notset
WORKDIR=$(mktemp -d)

while getopts i: flag
do
    case "${flag}" in
        i) IMAGE=${OPTARG};;
    esac
done

# Traffic ops
if [ "ops" == "$IMAGE" ]; then
  echo "OPS"
fi

echo "$WORKDIR"
echo "$IMAGE"