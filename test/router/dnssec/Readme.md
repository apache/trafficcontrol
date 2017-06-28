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

DNSSEC Tests
============

Running the test

`ginkgo -- -ns=router-01.thecdn.example.com:53  -ds=ds-01.thecdn.example.com.`

Sample Output
```
Running Suite: Dnssec Suite
===========================
Random Seed: 1476984556
Will run 4 of 4 specs

2016/10/20 11:29:17 Nameserver router-01.thecdn.example.com:53
2016/10/20 11:29:17 DeliveryService ds-01.thecdn.example.com.
••••
Ran 4 of 4 Specs in 0.110 seconds
SUCCESS! -- 4 Passed | 0 Failed | 0 Pending | 0 Skipped PASS

Ginkgo ran 1 suite in 825.345359ms
Test Suite Passed
```