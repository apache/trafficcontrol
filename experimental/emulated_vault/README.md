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

# Emulated Vault - Background

!!! Deprecated
  Since this tool is meant specifically to emulate the deprecated Riak backend
  for Traffic Vault, it will be removed when support for that backend is.

The emulated_vault module supplies a HTTP server mimicking RIAK behavior for usage as traffic-control vault.
It may be used in order to replace RIAK traffic_vault, as it is much more simple to install.
The server may use different type of persistent storage (e.g. file-system), using the proper adapter.
The resiliency of the stored keys is derived from the resiliency of the underlying storage.

# Installation

Basic requirements: Centos ver >= 7; Python >= 2.7

In order to install the module on a server please:
1. Copy the module files to the server's root
2. Add the certificate and key to your favorite path
3. Adjust /opt/emulated_vault/conf/cfg.json - pointing at your certificate and key
4. "systemctl enable" the service

Logs may be found under /opt/emulated_vault/var/log

# Developer's Notes

If you just want to play around with the module, you may of course run the server script on its own.
Before doing that, you would probably need to adjust the opt/emulated_vault/conf/cfg.json:
1. Changing the db-path to one you have access to
2. Disable ssl (just to make it easier)

Additionally, the vault-debug script is also available to work against the DB with command line.
It is mostly useful when developing a new adapter.

# Contact

For additional information, questions or assistance, please approach [Nir B. Sopher](mailto:nir@apache.org)
