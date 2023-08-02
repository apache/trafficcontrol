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

# repo-info docker action
This action builds an RPM name, after fetching the latest git commit sha and latest tag on a GitHub repo branch.

## Inputs

### `owner`
The owner username or organization owning the GitHub repository containing the branch being queried

### `repo`
The name of the GitHub repository containing the branch being queried

### `branch`
The name of the branch being queried

## Outputs

### `sha`
the commit sha from a GitHub repo branch

### `latest-tag`
the latest tag on the branch

### `exit-code`
0 if the tests passed, 1 otherwise.

## Example usage
```yaml
jobs:
  tests:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout
        uses: actions/checkout@master
      - name: Fetch RPM name
        uses: ./.github/actions/repo-info
          with:
            - owner: apache
            - repo: trafficserver
            - branch: 9.2.x
```
