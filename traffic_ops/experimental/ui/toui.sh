#!/usr/bin/env bash

# Temporary helper wrapper script
ACTION=$1

if [[ "$ACTION" == "run" ]];then
  docker-compose -p trafficops up -d
elif [[ "$ACTION" == "stop" ]];then
  docker-compose -p trafficops down -v
else
  echo "Invalid action: only 'run' or 'stop' available."
fi

