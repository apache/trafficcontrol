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
ATS
=========

Apache Traffic Control traditionally leverages the Apache Traffic Server caching proxy.  ATC uses its own configuration management tool to assist with the management of ATS named t3c.  This role deploys t3c to facilitate deployment and configuration of ATS according to the dataset modeled inside ATC.

Requirements
------------

A valid RPM in an available yum repository.

Role Variables
--------------

Refer to the defaults/main.yml for most information.

ort_version: This is an optional string that can be provided to specify a particular version of t3c to install.  It should be something like `3.0.0-10063.5db80eca.el7`.  The absense of this variable entails automatically using the latest version available to yum at the time of initial installation.

additional_yum_repos: An optional list of additional yum repositories to enable specifically when installing this component.  This could be used to enable non-production ready rpms in a separate repository and not supplying the specific RPM version to automatically use the latest available.

Dependencies
------------

None

Example Playbook
----------------
```yaml
  - name: Deploy ORT for ATS
    include_role:
      name: ats
    vars:
      install_ats: true
      ort_traffic_ops_url: "{{ to_url }}"
      ort_traffic_ops_username: ort_user
      ort_traffic_ops_password: "{{ ort_passwd}}"
      ort_crontab:
        syncds:
          schedule: '0,20,40 * * * *'
          user: root
          job: "t3c apply --run-mode=syncds --log-location-warning=stderr --log-location-error=stderr --traffic-ops-url='{{ ort_traffic_ops_url }}' --traffic-ops-user='{{ ort_traffic_ops_username }}' --traffic-ops-password='{{ ort_traffic_ops_password }}' &> /tmp/trafficcontrol-cache-config/syncds.log"
        reval:
          schedule: '1-19,21-39,41-59 * * * *'
          user: root
          job:  "t3c apply --run-mode=revalidate --log-location-warning=stderr --log-location-error=stderr --traffic-ops-url='{{ ort_traffic_ops_url }}' --traffic-ops-user='{{ ort_traffic_ops_username }}' --traffic-ops-password='{{ ort_traffic_ops_password }}' &> /tmp/trafficcontrol-cache-config/reval.log"
```

License
-------

Apache 2.0

Author Information
------------------

Apache Traffic Control
