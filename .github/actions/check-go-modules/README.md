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
# check-go-modules action

This action lets you perform checks to verify that
- `go.mod` is unmodified after installing modules
- `go.sum` is unmodified after installing modules
- After installing modules, the `vendor` directory contains no untracked files, modified files, or deleted files.

## Outputs
### `exit-code`

Exit code is 0 if the check succeeded.

## Example usage
```yaml
- run: .github/actions/check-go-modules/entrypoint.sh vendor_dependencies
- run: .github/actions/check-go-modules/entrypoint.sh check_vendored_deps
- run: .github/actions/check-go-modules/entrypoint.sh check_go_file go.mod
```