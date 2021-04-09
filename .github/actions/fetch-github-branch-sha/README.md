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

# fetch-github-branch-sha docker action
This action queries for the latest git commit 
sha on a github repo branch

## Outputs

### `sha` and `exit-code`
the commit sha from a github repo branch
0 if the tests passed, 1 otherwise.

## Example usage
```yaml
jobs:
  tests:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout
        uses: actions/checkout@master
      - name: Fetch GitHub commit
        uses: ./.github/actions/fetch-github-branch-sha
          with:
            - owner: apache
            - repo: trafficserver
            - branch: 8.1.x
```
