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
Grove
=========

Grove is an experimental caching proxy similar to Apache Traffic Server or NGINX optimized for linear video traffic.  This role deploys grove and it's associated ORT alternative named grovetccfg.

Requirements
------------

A valid pair of RPMs in an available yum repository.

Role Variables
--------------

Refer to the defaults/main.yml for most information.

grove_version: This is an optional string that can be provided to specify a particular version of grove to install.  It should be something like `0.2-10063.5db80eca`.  The absense of this variable entails automatically using the latest version available to yum at the time of initial installation.

grovetccfg_version: This is an optional string that can be provided to specify a particular version of grovetccfg to install.  It should be something like `0.2-10063.5db80eca`.  The absense of this variable entails automatically using the latest version available to yum at the time of initial installation.

additional_yum_repos: An optional list of additional yum repositories to enable specifically when installing this component.  This could be used to enable non-production ready rpms in a separate repository and not supplying the specific RPM version to automatically use the latest available.

Dependencies
------------

None

Example Playbook
----------------
```yaml
  - name: Deploy grove
    include_role:
      name: grove
    vars:
      install_grove: true
      grovetccfg_traffic_ops_url: https://to.kabletown.invalid
      grovetccfg_traffic_ops_username: username
      grovetccfg_traffic_ops_password: "{{ grove_passwd }}"
      grove_port: 80
      grove_https_port: 443
      grove_ssl_cert_path: /etc/pki/tls/certs/server.crt
      grove_ssl_key_path: /etc/pki/tls/private/server.key.pem
```

License
-------

Apache 2.0

Author Information
------------------

Apache Traffic Control
