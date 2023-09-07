<!--
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
-->

# Squashing database migrations

For convenience, [`squash_migrations.sh`](https://github.com/apache/trafficcontrol/blob/master/traffic_ops/app/db/squash_migrations.sh) script squashes the migrations, but whoever PRs the result is responsible for verifying that the migrations are squashed, regardless of the result of having run the script.

--------

Each major release of Apache Traffic Control combines database migrations from previous ATC releases into [`create_tables.sql`](https://github.com/apache/trafficcontrol/blob/master/traffic_ops/app/db/create_tables.sql).

For example, suppose the latest version of Apache Traffic Control is 147.5.8 and contains these migrations:
* `1_my-migration.up.sql`
* `1_my-migration.down.sql`
* `3_another-migration.up.sql`
* `3_another-migration.down.sql`

And suppose the ATC [`master`](https://github.com/apache/trafficcontrol/commits/master) branch contains these migrations:
* `1_my-migration.up.sql`
* `1_my-migration.down.sql`
* `3_another-migration.up.sql`
* `3_another-migration.down.sql`
* `4_migration-name.up.sql`
* `4_migration-name.down.sql`
* `9_add-column-to-table.up.sql`
* `9_add-column-to-table.down.sql`

1. In order to prepare database migrations for the next major release, in this case, ATC 148.0.0, migrations `1` and `3` should be collapsed into `create_tables.sql` and migrations `4` and `9` should remain in [`traffic_ops/app/db/migrations/`](https://github.com/apache/trafficcontrol/tree/master/traffic_ops/app/db/migrations/).

2. * After migrations from ATC 147.5.8 have been collapsed, the first migration version will be `4`. Find the definition for `FirstMigrationTimestamp` in [`traffic_ops/app/db/admin.go`](https://github.com/apache/trafficcontrol/blob/master/traffic_ops/app/db/admin.go) and change it to `4`.

Past PRs that have collapsed the DB migrations:
- https://github.com/apache/trafficcontrol/pull/6065
- https://github.com/apache/trafficcontrol/pull/3524
