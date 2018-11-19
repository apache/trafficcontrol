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

set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER traffic_ops WITH ENCRYPTED PASSWORD '$POSTGRES_PASSWORD';
    CREATE USER telegraf;
    CREATE USER grafana;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" \
    -tc "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME'" | grep -q 1 ||  \
    psql -U postgres -c "CREATE DATABASE $DB_NAME OWNER $DB_USER"

