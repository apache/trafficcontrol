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

Related files for use with this project and git.

## Hooks

A collection of useful pre-commit hooks can be found in the `pre-commit-hooks` directory.

#### Installing pre-commit hooks

In the `$GOPATH/src/github.com/apache/trafficcontrol/` directory, create a symbolic link from the `pre-commit` executable in this directory to the `.git/hooks/` directory:

```shell
ln -s ../../misc/git/pre-commit .git/hooks/
```

Now, all executables in the `pre-commit-hooks` directory will be run on commit.

#### Adding pre-commit check

Once the pre-commit file is in place, all executables in the `pre-commit-hooks` directory will be run. Simply add an executable there. Exiting with non-zero status from this script causes the git commit to abort (the commit contents will be unaffected).

#### Skipping

To commit without running the hooks, use the `no-verify` flag.

```bash
git commit --no-verify
```
