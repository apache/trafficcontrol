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
Traffic_portal_v2
=========

At the primary user interface of Apache Traffic Control is the component Traffic Portal V2 which is installed with this role.

Requirements
------------

* A valid RPM in an available yum repository.

Role Variables
--------------

Refer to the defaults/main.yml for most information.

tpv2_version: This is an optional string that can be provided to specify a particular version of Traffic Portal to install.  It should be something like `3.0.0-10063.5db80eca.el7`.  The absence of this variable entails automatically using the latest version available to yum at the time of initial installation.

config.json: Dictionary to merge into/atop the default traffic portal configuration values to file

Dependencies
------------

None

Example Playbook
----------------
```yaml
  - name: Deploy Traffic Portal V2
      import_role:
        name: traffic_portal_v2
      vars:
        install_traffic_portal: true
        tpv2_useSSL: true
        tpv2_http_port: 80
        tpv2_sslPort: 443
```

License
-------

Apache 2.0

Author Information
------------------

Apache Traffic Control
