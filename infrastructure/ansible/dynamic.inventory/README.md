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
# Ansible Dynamic Inventory

This python script uses the Traffic Ops Python Native Client to expose much of the TO dataset as an ansible inventory on demand as Ansible patterns.

## Requirements
You will need to ensure the Traffic Ops Python Native Client is available to the python env shared by Ansible.

Due to limitations in the way parameters are passed in Ansible Dynamic Inventory scripts, the following environment variables must be defined:
```bash
export TO_USERNAME=<my.to.username>
export TO_PASSWORD=<my.to.password>
export TO_URL=to.kabletown.invalid
```
Failure to set login credentials will result in a valid, but empty response.

If you find yourself debugging this or are curious what's available, the following commands are handy:
```bash
ansible-inventory -i infrastructure/ansible/dynamic.inventory/TO.py --graph --vars > ansible.inventory.txt
python infrastructure/ansible/dynamic.inventory/TO.py --list -to <to.username> <to.password> -url to.kabletown.invalid
```

## Example Usage
Use ansible ad-hoc to test connectivity to all offline edge caches belonging to the Kabletown2.0 CDN and having "den" in the fqdn somewhere, but not with "aurora" in their fqdn.
```bash
ansible -i infrastructure/ansible/dynamic.inventory/TO.py 'server_status|OFFLINE:&server_type|EDGE:&server_cdnName|Kabletown2.0:&*den*:!*aurora*' -m ping
```
