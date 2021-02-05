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
# save-alpine-tar action

This action lets save an already-downloaded Alpine image to a tar archive.

## Arguments

* `load-or-save` (required) `'load'` or `'save'`, depending on whether you want to load a Docker image from a tar archive or save an existing Docker image to a new tar archive.
* `digest` (required): A Docker image digest. Find this by running
  ```shell
  docker image ls --digests
  ```
  Example: `sha256:08d6ca16c60fe7490c03d10dc339d9fd8ea67c6466dea8d558526b1330a85930`

## Outputs
### `exit-code`

* Exit code is 0 unless either
  - An invalid value for `load-or-save` was given, or
  - An invalid Docker image digest for the `digest` input was given

## Example usage
Load an Alpine Docker image:
  ```yaml
  - run: .github/actions/save-alpine-tar/entrypoint.sh load sha256:08d6ca16c60fe7490c03d10dc339d9fd8ea67c6466dea8d558526b1330a85930
  ```

Save an Alpine Docker image:
  ```yaml
  - run: .github/actions/save-alpine-tar/entrypoint.sh save sha256:08d6ca16c60fe7490c03d10dc339d9fd8ea67c6466dea8d558526b1330a85930
  ```
