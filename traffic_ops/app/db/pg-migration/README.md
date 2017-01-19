# Converting existing mysql `traffic_ops` database to postgres

## Overview
  - This conversion will convert an existing Mysql database by pulling the data from Traffic Ops (using Tools.dbdump) and loading that 
    data into a temporary Mysql Database to make it easier to work with in the conversion process.  Additionally, it will use the 'pgloader' tool
    to perform that conversion in yet another Docker container which will then push that converted data into your permanent Postgres instance(s).

Software requirements
* Requires a recent ( 1.12 ) version of `docker-engine` and `docker-compose`.

* Modify the postgres.env for your new Postgres Database
  (NOTE: do not set the POSTGRES_HOST to 'localhost' it needs to be the IP address or DNS available hostname so that the container can reach out to Postgres)

* Ensure that your new Postgres service is running (local or remote)

* Run the Mysql to Postgres Migration Docker flow
  * `$ docker-compose -f pgmigration.yml down -v && docker-compose -f pgmigration.yml build && TO_SERVER=https://traffic_ops.kabletown.com TO_USER=me TO_PASSWORD='my!passwd' docker-compose -f pgmigration.yml up`

* Run the Postgres datatype conversion
  * `$ docker-compose -f convert.yml up` 
