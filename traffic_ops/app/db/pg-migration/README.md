
# Converting existing mysql `traffic_ops` database to postgres

* Requires a recent ( 1.12 ) version of `docker-engine` and `docker-compose`.

* Ensure database is up-to-date with latest `traffic_ops` migrations for `mysql` (last 1.x release of `traffic_ops`)
  * `cd /opt/traffic_ops/app;  ./db/admin.pl --env=production upgrade`

* In development environment, `cd traffic_ops/app/db/pg-migration`.

* Provide URL, username password for existing mysql install of `traffic_ops`:

  * `docker-compose -f docker-compose-pgmigration.yml down -v && \
	 docker-compose -f docker-compose-pgmigration.yml build && \
	 TO_SERVER=https://traffic_ops.kabletown.com TO_USER=me TO_PASSWORD='my!passwd' docker-compose -f docker-compose-pgmigration.yml up`

* Postgres is still running in a docker container -- dump the database to a file:
  `docker exec -it pgmigration_postgres_host_1 pg_dump -Utraffic_ops >pg.sql`

* Or examine the database directly:
  * `docker run -it --rm --network pgmigration_default --link pgmigration_postgres_host_1:postgres postgres psql -h postgres -U traffic_ops -d traffic_ops`
  * `\dt`
  * `select * from cdns;`
