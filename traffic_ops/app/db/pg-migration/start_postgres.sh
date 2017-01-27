#!/bin/bash 
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

. mysql-to-postgres.env

#Traffic Ops Settings
# The following configs should be configured to point to the 
# Traffic Ops instances that is connected to the MySQL that 
# you want to convert
separator="---------------------------------------"

function display_env() {

  echo "TO_SERVER: $TO_SERVER"
  echo "TO_USER: $TO_USER"
  echo "TO_PASSWORD: $TO_PASSWORD"
  echo "MYSQL_HOST: $MYSQL_HOST"
  echo "MYSQL_PORT: $MYSQL_PORT"
  echo "MYSQL_DATABASE: $MYSQL_DATABASE"
  echo "MYSQL_USER: $MYSQL_USER"
  echo "MYSQL_PASSWORD: $MYSQL_PASSWORD"

  echo "POSTGRES_HOST: $POSTGRES_HOST"
  echo "POSTGRES_PORT: $POSTGRES_PORT"
  echo "POSTGRES_DATABASE: $POSTGRES_DATABASE"
  echo "POSTGRES_USER: $POSTGRES_USER"
  echo "POSTGRES_PASSOWRD: $POSTGRES_PASSWORD"
  echo "PGDATA: $PGDATA"
  echo "PGDATA_VOLUME: $PGDATA_VOLUME"
  echo "PGLOGS_VOLUME: $PGLOGS_VOLUME"

}

function start_staging_postgres_server() {

  docker-compose -v -p trafficops -f postgres.yml down

  echo "PGDATA_VOLUME: $PGDATA_VOLUME"
  echo "PGLOGS_VOLUME: $PGLOGS_VOLUME"
  PGLOGS_VOLUME=$PGLOGS_VOLUME PGDATA_VOLUME=$PGDATA_VOLUME docker-compose -p trafficops -f postgres.yml up  -d

}

start_staging_postgres_server
