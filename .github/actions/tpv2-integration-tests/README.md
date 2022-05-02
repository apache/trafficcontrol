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

# tp-integration-tests javascript action
this action runs the traffic portal integration tests

## inputs

## outputs

### `exit-code`
returns non-zero exit code on failure

## example usage
```yaml
jobs:
  E2E_tests:
    if: github.event.pull_request.draft == false
    runs-on: ubuntu-latest
    services:
      hub:
        image: selenium/hub
        ports:
          - "4444:4444"
        options: --health-cmd=/opt/bin/check-grid.sh --health-interval=5s --health-timeout=15s --health-retries=5
      chrome:
        image: selenium/node-chrome
        env:
          HUB_HOST: hub
          HUB_PORT: 4444
        volumes:
          - /dev/shm:/dev/shm
      postgres:
        image: postgres:11
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
      - name: Initialize Traffic Ops Database
        id: todb
        uses: ./.github/actions/todb-init
      - name: Run TP
        uses: ./.github/actions/tpv2-integration-tests
```
