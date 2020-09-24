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
Traffic_vault
=========

Traffic Vault is the component of Apache Traffic Control which stores sensitive delivery service data using the opensource Basho Technologies NoSQL database [Riak_KV](https://riak.com/products/riak-kv/index.html).  This role handles the installation of Riak.

Requirements
------------

* A valid RPM in an available yum repository.

Role Variables
--------------

Refer to the defaults/main.yml for most information.  For information regarding Riak configuration please consult their [documentation](https://docs.riak.com/riak/kv/latest/configuring/basic/index.html).

additional_yum_repos: An optional list of additional yum repositories to enable specifically when installing this component.  This could be used to enable non-production ready rpms in a separate repository and not supplying the specific RPM version to automatically use the latest available.

Dependencies
------------

None

Example Playbook
----------------
```yaml
  - name: Deploy Traffic Vault
    import_role:
      name: traffic-vault
    vars:
      install_traffic_vault: true
      riak_nodename: "riak@{{ ansible_default_ipv4.address }}"
      riak_erlang_max_ports: 65536
      riak_listener_protobuf_internal: "{{ ansible_default_ipv4.address }}:8087"
      riak_protobuf_backlog: 4096
      riak_listener_https_internal: "{{ ansible_default_ipv4.address }}:8088"
      riak_ringleader: "riak@{{ hostvars[groups['riak'] | first].ansible_default_ipv4.address }}"
      riak_ssl_certfile: /etc/pki/tls/certs/server.crt
      riak_ssl_keyfile: /etc/pki/tls/private/server.key.pem
      riak_generate_ssl: false
```

License
-------

Apache 2.0

Author Information
------------------

Apache Traffic Control
