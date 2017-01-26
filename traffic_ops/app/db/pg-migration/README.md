# Converts existing mysql `traffic_ops` database to postgres

* Requires a recent ( 1.12 ) version of `docker-engine` and `docker-compose`.

* Modify the mysql-to-postgres.env file for the parameters in your Migration 

* Ensure that your new Postgres service is running (local or remote)

* Run the Mysql to Postgres Migration Docker flow
  * `sh migrate.sh
