#!/usr/bin/env bash
#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.


set -ex

d=/docker-entrypoint-initdb.d
for dump in "$d"/*dump; do
    [[ -f "$dump" ]] || break
    t=$(mktemp -p /tmp XXX.sql)
    # convert to sql -- can't load a dump until db initialized,  but sql works
    echo "Restoring from $dump"
    pg_restore -f "$t" "$dump"
    if [[ "${POSTGRES_VERSION%%.*}" -gt 10 ]]; then
      sed -i '/^CREATE SCHEMA public;$/d' "$t"
    fi
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$DB_NAME" <"$t"
done
