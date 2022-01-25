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
Dataset_loader
=========

This role binds together topological information from ansible and integrates that into an initial dataset for Apache Traffic Control.

For assistance with converting an existing TO server profile to a yaml-fied version mostly suitable for diff and merge to this dataset, see [this extra documentation](profile.parameter.conversion.md).

Requirements
------------

A freshly installed and functional Traffic Ops.  Users of this playbook should also be familiar with the Traffic Ops datamodel and be able to interpret/research Apache Traffic Server configuration.

Today this dataset loader assumes mostly pools of resources such as caches, cachegroups, regions, etc.  It will perform a round-robin assignment within the resource pools automatically.  If you need a different arrangement, you'll want to craft your own supplamentory modification playbook or bypass the dataset loader with a static dataset.

A web server to host the static files for CZF and GeoIP needs to be responding with the required content.

An ATS RPM that matches the version described by `ats_version` to be available in an accessible yum repo to ORT.

It is also pesumed that you're using a working directory that contains `./out/ssl/lab.intermediateca.crt` and `./out/ssl/lab.intermediateca.pem` in order to sign Delivery Service SSL CSRs

Role Variables
--------------

Refer to the defaults/main.yml for most information.

ats_version: This is a string to specify a particular version of ATS to install via ORT.  It should be something like `7.1.4-2.el7`.

Dependencies
------------

All hosts in the inventory except those in the mso_parent_alias group must be reachable in order to retrieve facts to become part of the dataset.  The inventory should include the groups, groupvars, and hostvars as outlined in the ATC provisioning layer contract documentation.

Example Playbook
----------------
```yaml
- hosts: all:!mso_parent_alias
  gather_facts: true

  tasks:

- hosts: localhost
  connection: local
  gather_facts: false


  tasks:
    - name: Load environment specific vars
      include_vars:
        file: "{{ lookup('env', 'PWD') }}/ansible/vars.json"
      no_log: true

    - name: Deploy the Initial ATC Dataset
      import_role:
        name: dataset_loader
      vars:
        load_dataset: true
        dl_shallow_czf_url: http://examplewebserver.local/czf.json
        dl_allow_ip4: "{{ (['127.0.0.1','192.168.100.0/24'] + (dl_hosts_tm | map('extract', hostvars, ['ansible_default_ipv4', 'address']) | list)) | join(',') }}"
        dl_allow_ip6: "{{ (['::1'] + (dl_hosts_tm | map('extract', hostvars, ['ansible_default_ipv6', 'address']) | list)) | join(',') }}"
        dl_to_url: "{{ to_url }}"
        dl_ds_merged_cdns: "{{ dl_ds_default_cdns | combine(dl_ds_test_cdns) }}"
        dl_ds_test_cdns:
          AnotherCDN:
            name: Another.Test.CDN
            dnssecEnabled: false
        dl_ds_merged_types: "{{ dl_ds_default_types + dl_ds_private_types }}"
        dl_ds_private_types:
          - name: INFLUX_RELAY
            description: Influx Relay Server type
            useInTable: server
```

License
-------

Apache 2.0

Author Information
------------------

Apache Traffic Control
