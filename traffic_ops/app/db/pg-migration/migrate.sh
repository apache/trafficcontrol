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

: ${TO_SERVER?"Please set the TO_SERVER environment variable: ie: https://kabletown.net"}
: ${TO_USER?"Please set the TO_USER environment variable: ie: <your Traffic Ops userid>"}
: ${TO_PASSWORD?"Please set the TO_PASSWORD environment variable: ie: <your Traffic Ops password>"}

MYSQL_PORT=3306
POSTGRES_PORT=5432

separator="---------------------------------------"

function shutdown_trafficops_database() {
     sudo systemctl stop trafficops-db
}

function start_staging_mysql_server() {
      docker-compose -p trafficops -f mysql_host.yml up --build -d
      while [[ ! `netstat -lnt | grep :$MYSQL_PORT` ]]; do
	    # wait for signal that other container is waiting
	    echo "Waiting for Mysql to Start..."
	    sleep 3
      done
      echo $separator
      echo "Mysql Host is started..."
      echo $separator
}

function start_staging_postgres_server() {
	sudo systemctl start trafficops-db
	while [[ ! `netstat -lnt | grep :$POSTGRES_PORT` ]]; do
	    # wait for signal that other container is waiting
	    echo "Waiting for Postgres to Start..."
	    sleep 3
	done
	echo $separator
	echo "Postgres started.."
	echo $separator
}


function run_postgres_datatypes_conversion() {
	echo $separator
	echo "Starting Mysql to Postgres Migration..."
	echo $separator
	docker-compose -p trafficops -f convert.yml up --build
}


function migrate_data_from_mysql_to_postgres() {
	echo $separator
	echo "Starting Mysql to Postgres Migration..."
	echo $separator
	docker-compose -p trafficops -f mysql-to-postgres.yml up --build
}

function clean() {
        echo $separator
        echo "Cleaning up..."
        echo $separator
        docker kill trafficops_mysql_host_1
        docker-compose -p trafficops -f mysql-to-postgres.yml down --remove-orphans
        docker-compose -p trafficops -f convert.yml down --remove-orphans
        docker rm trafficops_mysql-to-postgres_1 
        docker rm trafficops_convert_1
        docker rm trafficops_mysql_host_1
        docker rmi trafficops_mysql-to-postgres
        docker rmi trafficops_convert 
        docker rmi trafficops_mysql_host
        docker rmi mysql:5.6 
        docker rmi dimitri/pgloader:latest
}

start_staging_mysql_server
start_staging_postgres_server
migrate_data_from_mysql_to_postgres
run_postgres_datatypes_conversion
clean
