#!/usr/bin/env bash
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
trap 'echo "Error on line ${LINENO} of ${0}" >/dev/stderr; exit 1' ERR
set -o errexit -o nounset -o pipefail
cd "$(dirname "$0")"

apache_remote="$(git remote -v | grep 'apache/trafficcontrol.*(fetch)$' | cut -f1 | head -n1)"
git fetch --tags --force "$apache_remote"
last_release="$(git tag --list --sort=v:refname | grep -E '^RELEASE-[0-9]+[.][0-9]+[.][0-9]+$' | tail -n1)"
migrations_to_squash="$(git ls-tree -r "$last_release" -- migrations | grep -o 'migrations/[0-9].*')"

cp create_tables.sql to_squash.sql
echo "$migrations_to_squash" | grep '\.up\.sql$' | xargs cat >>to_squash.sql
last_squashed_migration="$(<<<"$migrations_to_squash" tail -n1)"
last_squashed_migration_timestamp="$(<<<"$last_squashed_migration" sed -E 's|migrations/([0-9]+).*|\1|')"
first_migration="$(ls migrations/*.sql | grep -A1 "/${last_squashed_migration_timestamp}_" | tail -n1)"
first_migration_timestamp="$(<<<"$first_migration" sed -E 's|migrations/([0-9]+).*|\1|')"
sed -i.bak '/^--/,$d' create_tables.sql # keeps the Apache License 2.0 header
sed -Ei.bak "s|(LastSquashedMigrationTimestamp\s+uint\s+= ).*|\1${last_squashed_migration_timestamp} // ${last_squashed_migration}|" admin.go
sed -Ei.bak "s|(FirstMigrationTimestamp\s+uint\s+= ).*|\1${first_migration_timestamp} // ${first_migration}|" admin.go

dump_db_with_migrations() {
  trap 'echo "Error on line ${LINENO} of dump_db_with_migrations" >/dev/stderr; exit 1' ERR
  set -o errexit -o nounset
  {
    docker-entrypoint.sh postgres &
    sleep 10
    psql -f to_squash.sql
  } >/dev/stderr
  pg_dump
}
docker run --rm -iw/db \
  -v "$(pwd):/db" \
  -e PGUSER=traffic_ops \
  -e PGPASSWORD=twelve \
  -e POSTGRES_USER=traffic_ops \
  -e POSTGRES_PASSWORD=twelve \
  postgres:13-alpine bash -c "$(type dump_db_with_migrations | tail -n+2); dump_db_with_migrations" >>create_tables.sql
rm to_squash.sql

git add create_tables.sql
git commit -m "Redump create_tables.sql with migrations through timestamp ${last_squashed_migration_timestamp}"

echo "$migrations_to_squash" | xargs git rm
git commit -m "Remove migrations that existed at ${last_release}"

git add -p admin.go
git commit -m 'Update LastSquashedMigrationTimestamp and FirstMigrationTimestamp'

echo 'Migrations squashed successfully!'
