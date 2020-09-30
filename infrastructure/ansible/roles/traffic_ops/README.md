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
Traffic_ops
=========

At the center of Apache Traffic Control is the database API known as Traffic Ops.  The TrafficOps installation also contains relevant initialization for both the Postgresql database as well as Riak SOLR indexing; all of which are handled by this role.

Requirements
------------

* A valid RPM in an available yum repository.
* An internet connection to obtain external dependencies
* A functional Postgresql 9.6 instance
* A functional Traffic Vault (Riak) instance with SOLR

Role Variables
--------------

Refer to the defaults/main.yml for most information.

to_version: This is an optional string that can be provided to specify a particular version of Traffic Ops to install.  It should be something like `3.0.0-10063.5db80eca.el7`.  The absense of this variable entails automatically using the latest version available to yum at the time of initial installation.

additional_yum_repos: An optional list of additional yum repositories to enable specifically when installing this component.  This could be used to enable non-production ready rpms in a separate repository and not supplying the specific RPM version to automatically use the latest available.

Dependencies
------------

None

Example Playbook
----------------
```yaml
  - name: Deploy Traffic_ops
    import_role:
      name: traffic_ops
    vars:
      install_traffic_ops: true
      to_certs_cert: /etc/pki/tls/certs/server.crt
      to_certs_key: /etc/pki/tls/private/server.key.pem
      to_certs_ca: /etc/pki/ca-trust/source/anchors/lab.rootca.crt
```

License
-------

Apache 2.0

Author Information
------------------

Apache Traffic Control
