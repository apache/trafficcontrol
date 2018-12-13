..
..
.. Licensed under the Apache License, Version 2.0 (the "License");
.. you may not use this file except in compliance with the License.
.. You may obtain a copy of the License at
..
..     http://www.apache.org/licenses/LICENSE-2.0
..
.. Unless required by applicable law or agreed to in writing, software
.. distributed under the License is distributed on an "AS IS" BASIS,
.. WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
.. See the License for the specific language governing permissions and
.. limitations under the License.
..

.. _traffic_vault_util:

******************
Traffic Vault Util
******************
The ``traffic_vault_util`` tool is used to view and modify the contents of a :ref:`Traffic Vault` (i.e. Riak) cluster. The tool contains basic list_* operations to display the buckets, keys and values stored within TV.

``traffic_vault_util`` also has a small converter utility to perform a one-off conversion of key formats within the SSL bucket. This conversion is useful when moving from an older version of Traffic Ops to the current version. In the older version, SSL records were indexed by Delivery Service database ID. Currently, SSL records are indexed by Delivery Service XML ID. 

Usage
=====
Usage of ./traffic_vault_util:
  -dry_run
    	Do not perform writes
  -vault_action string
    	Action: list_buckets|list_keys|list_values|convert_ssl_to_xmlid
  -vault_ip string
    	IP/Hostname of Vault
  -vault_password string
    	Riak Password
  -vault_port uint
    	Protobuffers port of Vault (default 8087)
  -vault_user string
    	Riak Username

