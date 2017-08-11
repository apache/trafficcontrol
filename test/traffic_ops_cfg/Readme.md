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

This test allows you to compare all generated config files and CDN snapshots (CRConfig.json) from 2 instances of Traffic Ops. For example, you could compare config files / snapshots of a MySQL vs Postgres Traffic Ops. You could even compare across releases (1.7.0 vs 1.8.0).

*Prerequisites*

1. Make sure the data in your 2 Traffic Ops databases are synced to avoid getting false positives.
2. Modify test.config with proper settings. Set perform_snapshot=1 if you want to force a snapshot in both Traffic Ops instances.

*Running the Test*

1. `./cfg_test.pl getref test.config` your ref files go into `/tmp/files/ref`
2. `./cfg_test.pl getnew test.config` your new files go into `/tmp/files/new`
3. `./cfg_test.pl compare test.config` - all `not ok` lines should be explained.

It will compare all files for all profiles, _including_ the CRConfig.json. 

