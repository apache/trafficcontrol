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

# tvdb-init Docker action
This action initializes a Traffic Vault database in a locally running PostgresQL
instance. It creates the `traffic_vault` user and database.

## Outputs

### `exit-code`
1 if the database initialized successfully.

## Example usage
```yaml
jobs:
  tests:
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:11.9
        env:
          POSTGRES_USER: traffic_ops
          POSTGRES_PASSWORD: twelve
          POSTGRES_DB: traaffic_ops
        ports:
        - 5432:5432
        options: --health-cmd pg_isready --health-interval 10s --health-timeout 5s --health-retries 5

    steps:
      - name: Checkout
        uses: actions/checkout@master
      - name: ./.github/actions/tvdb-init
```
