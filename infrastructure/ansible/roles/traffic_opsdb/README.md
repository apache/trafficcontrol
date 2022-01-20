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
Traffic_opsdb
=========

This is an optional role to facilitate getting Postgresql9.6 installed and initialized for Traffic Ops.

Requirements
------------

* An internet connection to obtain external dependencies

Role Variables
--------------

Refer to the defaults/main.yml for most information.

additional_yum_repos: An optional list of additional yum repositories to enable specifically when installing this component.  This could be used to enable non-production ready rpms in a separate repository and not supplying the specific RPM version to automatically use the latest available.

Dependencies
------------

None

Example Playbook
----------------
```yaml
---
- hosts: traffic_opsdb
  gather_facts: yes
  become: yes

  tasks:
    - name: Load environment specific vars
      include_vars:
        file: "{{ lookup('env', 'PWD') }}/ansible/vars.json"
      no_log: true

    - name: Deploy Traffic_opsdb
      import_role:
        name: traffic_opsdb
      vars:
        install_traffic_opsdb: true

- hosts: traffic_opsdb-primary
  gather_facts: yes
  become: yes

  tasks:
    - name: Load environment specific vars
      include_vars:
        file: "{{ lookup('env', 'PWD') }}/ansible/vars.json"
      no_log: true

    - name: Initialize Traffic_opsdb
      import_role:
        name: traffic_opsdb
      vars:
        initialize_traffic_opsdb: true

```

License
-------

Apache 2.0

Author Information
------------------

Apache Traffic Control
