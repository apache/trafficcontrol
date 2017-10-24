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

# Traffic Ops Config File / Snapshot Compare

This test allows you to compare all generated config files (from a sample of caches across all profiles) and CDN snapshots (CRConfig.json) from 2 instances of Traffic Ops. For example, you could compare config files / snapshots between 2 versions of Traffic Ops.

*Prerequisites*

1. Make sure the data in your 2 Traffic Ops databases are synced to avoid getting false positives.
2. Modify test.config with proper settings. Set perform_snapshot=1 if you want to force a snapshot in both Traffic Ops instances.

*Running the Test*

1. `./cfg_test.pl getref test.config` - files from Traffic Ops #1 go into `/tmp/files/ref`
2. `./cfg_test.pl getnew test.config` - files from Traffic Ops #2 go into `/tmp/files/new`
3. `./cfg_test.pl compare test.config` - all `not ok` lines should be explained.

It will compare all config files for all profiles, _including_ the CRConfig.json.

