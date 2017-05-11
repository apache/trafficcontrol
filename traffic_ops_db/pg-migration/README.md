# Converts existing mysql `traffic_ops` database to postgres

* Requires a fairly recent ( 17.05.0-ce ) version of `docker-engine` and `docker-compose`.

* Modify the mysql-to-postgres.env file for the parameters in your Traffic Ops environment

* Ensure that your new Postgres service is running (local or remote)

* A sample Postgres Docker container has been provided for testing
  1. `cd ../docker`
  2. `$ sh todb.sh run` - to download/start your Postgres Test Docker container
  3. `$ sh todb.sh setup` - to create your new 'traffic_ops' role and 'traffic_ops' database

* Run the Mysql to Postgres Migration Docker flow
  1. `$ cd ../pg-migration`
  2. `$ sh migrate.sh`
