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

# tm-integration-tests JavaScript action
This action runs the Traffic Monitor integration tests with a fake TO API and fake ATS.

## Inputs

## Outputs

### `exit-code`
1 if the Go program(s) could be built successfully.

## Example usage
```yaml
jobs:
  TM_Integration_tests:
    if: github.event.pull_request.draft == false
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@master
        with:
          fetch-depth: 1
      - name: Run integration tests
        uses: ./.github/actions/tm-integration-tests
        env:
          GOPATH: /github/workspace

```
