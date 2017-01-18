# Converting existing mysql `traffic_ops` database to postgres

* Requires a recent ( 1.12 ) version of `docker-engine` and `docker-compose`.

* Modify the mysql.env for your existing Mysql Database

* Modify the postgres.env for your new Postgres Database
  (NOTE: do not set the POSTGRES_HOST to 'localhost' it needs to be the IP address or DNS available hostname so that the container can reach out to Postgres)

* Ensure that your new Postgres service is running (local or remote)

* Run the Mysql to Postgres Migration Docker flow
  * `$ docker-compose down -v && docker-compose build && TO_SERVER=https://traffic_ops.kabletown.com TO_USER=me TO_PASSWORD='my!passwd' docker-compose -f pgmigration.yml up`

* Run the Postgres datatype conversion
  * `$ docker-compose -f convert.yml up` 
