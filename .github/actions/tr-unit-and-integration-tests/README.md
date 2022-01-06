<!--
  Licensed to the Apache Software Foundation (ASF) under one
  or more contributor license agreements.  See the NOTICE file
  distributed with this work for additional information
  regarding copyright ownership.  The ASF licenses this file
  to you under the Apache License, Version 2.0 (the
  "License"); you may not use this file except in compliance
  with the License.  You may obtain a copy of the License at

    https://apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing,
  software distributed under the License is distributed on an
  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
  KIND, either express or implied.  See the License for the
  specific language governing permissions and limitations
  under the License.
-->

# tr-unit-and-integration-tests Docker action
This action runs the Traffic Router unit tests and integration tests in an Alpine Docker container.

## Inputs

## Outputs

### `exit-code`
0 for success, nonzero for failure

## Example usage
```yaml
jobs:
  tests:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@master
      - name: Cache local Maven repository
        uses: actions/cache@v2
        with:
          path: ${{ github.workspace }}/.m2/repository
          key: ${{ runner.os }}-maven-${{ hashFiles('**/pom.xml') }}
          restore-keys: |
            ${{ runner.os }}-maven-
      - name: Run unit tests and integration tests
        uses: ./.github/actions/tr-unit-tests
```

To run the tests locally:
```shell
export GITHUB_WORKSPACE='/github/workspace';
docker build -f .github/actions/tr-unit-and-integration-tests/Dockerfile -t tr-unit-and-integration-tests .;
docker run --rm -te GITHUB_WORKSPACE -v "$(pwd):${GITHUB_WORKSPACE}" -w "$GITHUB_WORKSPACE" tr-unit-and-integration-tests;
```
