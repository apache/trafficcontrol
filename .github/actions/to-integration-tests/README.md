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

# to-integration-tests JavaScript action
This action runs the Traffic Ops Go client integration tests with the Traffic Ops API.
- Requires an SMTP service (see `smtp_address` input)

## Inputs

### `version`
**Required** Major API version to test e.g. 1, 2, 3 etc.

### `smtp_address`
**Required** The address of an SMTP server for use by Traffic Ops.

### `smtp_port`
**Required** The address of an SMTP server for use by Traffic Ops. Required but defaults to `25`.

### `smtp_user`
**Optional** The user to authenticate with for the SMTP server.

### `smtp_password`
**Optional** The password to authenticate with for the SMTP server.

## Outputs

### `exit-code`
1 if the Go program(s) could be built successfully.

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
          POSTGRES_DB: traffic_ops
        ports:
        - 5432:5432
        options: --health-cmd pg_isready --health-interval 10s --health-timeout 5s --health-retries 5

      smtp:
        image: maildev/maildev:2.0.0-beta3
        ports:
          - 25:25
        options: >-
          --entrypoint=bin/maildev
          --user=root
          --health-cmd="sh -c \"[[ \$(wget -qO- http://smtp/healthz) == true ]]\""
          --
          maildev/maildev:2.0.0-beta3
          --smtp=25
          --hide-extensions=STARTTLS
          --web=80

    steps:
      - name: Checkout
        uses: actions/checkout@master
      - name: initialize database
        uses: ./.github/actions/todb-init
      - name: Run API v5 tests
        uses: ./.github/actions/to-integration-tests
        with:
          version: 5
          smtp_address: localhost
```
