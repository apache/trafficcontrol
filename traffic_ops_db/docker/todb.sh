#!/usr/bin/env bash

#
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.
#

# Docker wrapper script to start/manage your Postgres container

ACTION=$1

PROJECT_NAME=trafficops
CONTAINER_NAME=trafficops_db_1

if [[ "$ACTION" == "setup" ]];then
  createuser traffic_ops -U postgres -h localhost
  createdb traffic_ops --owner traffic_ops -U postgres -h localhost 
  psql -c "\l" -U postgres -h localhost
  echo "Setup complete"

elif [[ "$ACTION" == "run" ]];then
  docker compose -f docker-compose.dev.yml -p $PROJECT_NAME up -d
  echo "Started Docker container: $CONTAINER_NAME"

elif [[ "$ACTION" == "stop" ]];then
  docker compose -f docker-compose.dev.yml -p $PROJECT_NAME down -v
  echo "Stopping Docker container: $CONTAINER_NAME"

elif [[ "$ACTION" == "clean" ]];then
  docker stop $CONTAINER_NAME
  docker rm $CONTAINER_NAME
  docker rmi $PROJECT_NAME
  rm -rf pgdata

elif [[ "$ACTION" == "seed" ]];then
  docker compose -f docker-compose.dev.yml -p $PROJECT_NAME down -v
else
  echo "Valid actions: 'setup', 'run', 'stop', 'clean'"
  exit 0
fi

docker ps
