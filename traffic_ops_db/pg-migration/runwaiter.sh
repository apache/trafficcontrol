#!/bin/bash


echo "WAITER_HOST: $WAITER_HOST"
echo "WAITER_PORT: $WAITER_PORT"

for c in "$WAITER_HOST $WAITER_PORT" ; do
  while true; do
     echo "Waiting 3 seconds for Host and Port: $c "
     sleep 3
     nc -z $c && break
  done
done
