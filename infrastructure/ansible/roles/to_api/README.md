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
TO_api
=========

There are certain tasks performed by ATC operators which involve multiple interactions with the Traffic Ops API and may include certain safety checks handled by the client.  This role wraps around several of these tasks.

Requirements
------------

* For snapshot:
  * [jq](https://stedolan.github.io/jq/)
  * [jd](https://github.com/josephburnett/jd)

Role Variables
--------------

Refer to the defaults/main.yml for most information.

Dependencies
------------

None

Example Playbook
----------------
```yaml
---
- name: Take caches out of service
  hosts: localhost
  gather_facts: no

  tasks:
    - name: Verify there are no pending changes for traffic router
      include_role:
        name: to_api
        tasks_from: snapshot.yml
      vars:
        to_api_target_cdn: "MKGA"
        to_api_commit_snapshot: false
        to_api_assert_clean_snapshot: true

    - name: Set cacheing server state
        include_role:
          name: to_api
          tasks_from: set_server_status.yml
        vars:
          to_api_target_host: "{{ item }}"
          to_api_desired_state: 'ADMIN_DOWN'
          to_api_down_reason: "Planned Maintinance"
        with_items: "{{ groups['target_systems'] }}"

    - name: Commit the server state changes to traffic router
      include_role:
        name: to_api
        tasks_from: snapshot.yml
      vars:
        to_api_target_cdn: "MKGA"
        to_api_commit_snapshot: true
        to_api_assert_clean_snapshot: false
```

License
-------

Apache 2.0

Author Information
------------------

Apache Traffic Control
