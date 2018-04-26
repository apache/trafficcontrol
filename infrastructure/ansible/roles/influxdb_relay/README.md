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
Influxdb-relay
=========

When using the Traffic Stats component of Apache Traffic Control, InfluxDB is required.  This playbook deploys an additional project [influx-relay](https://github.com/influxdata/influxdb-relay) between Traffic Stats and InfluxDB to provide basic HA capabilities to the open-source InfluxDB project.

Requirements
------------

A valid RPM in an available yum repository.

Role Variables
--------------

Refer to the defaults/main.yml for most information and the [influx-relay](https://github.com/influxdata/influxdb-relay) project.

influxdb_relay_version: This is an optional string that can be provided to specify a particular version of ORT to install.  It should be something like `adaa2ea-1`.  The absense of this variable entails automatically using the latest version available to yum at the time of initial installation.

additional_yum_repos: An optional list of additional yum repositories to enable specifically when installing this component.  This could be used to enable non-production ready rpms in a separate repository and not supplying the specific RPM version to automatically use the latest available.

Dependencies
------------

None

Example Playbook
----------------
See the [sample playbook](../../influxdb_relay.yml).

License
-------

Apache 2.0

Author Information
------------------

Apache Traffic Control
