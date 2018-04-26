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
Traffic_stats
=========

Traffic stats is an optional Apache Traffic Control component which collects CDN metrics from Traffic Monitor and stores them in the timeseries database InfluxDB.  This role performs the initial Traffic Stats installation and InfluxDB database creation.

Requirements
------------

* A valid RPM in an available yum repository.
* An one or more InfluxDB

Role Variables
--------------

Refer to the defaults/main.yml for most information.

ts_version: This is an optional string that can be provided to specify a particular version of Traffic Stats to install.  It should be something like `3.0.0-10063.5db80eca.el7`.  The absense of this variable entails automatically using the latest version available to yum at the time of initial installation.

additional_yum_repos: An optional list of additional yum repositories to enable specifically when installing this component.  This could be used to enable non-production ready rpms in a separate repository and not supplying the specific RPM version to automatically use the latest available.

Dependencies
------------

None

Example Playbook
----------------
```yaml
  - name: Deploy Traffic Stats
    import_role:
      name: traffic_stats
    vars:
      install_traffic_stats: true
      ts_influxUrls: "{{ groups['influxrelay'] | map('regex_replace', '^(.*)', 'http://\\1:9086') | list }}"
      ts_influxdb_hosts: "{{ groups['influxdb'] }}"
```

License
-------

Apache 2.0

Author Information
------------------

Apache Traffic Control
