
# Converting existing mysql `traffic_ops` database to postgres

* Requires `docker-engine` and `docker-compose`.

* Ensure database is up-to-date with latest `traffic_ops` migrations for `mysql` (last 1.x release of `traffic_ops`)
  * `cd /opt/traffic_ops/app;  ./db/admin.pl --env=production upgrade`

* Get a database dump from `traffic_ops`
  * `Tools->DB Dump`

* In development environment, `cd traffic_ops/app/db/pg-migration`.

* Move the `mysql` database dump file into `./mysql/initdb.d` directory.  The file must have a `.sql` suffix.

* `docker-compose down -v && docker-compose build && docker-compose up`

* Postgres is still running in a docker container -- dump the database to a file:
  `docker exec -it pgmigration_postgres_host_1 pg_dump -Utraffic_ops >pg.sql`
