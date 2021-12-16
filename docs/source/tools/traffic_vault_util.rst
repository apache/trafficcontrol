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

*************************
Traffic Vault Util (Riak)
*************************
.. deprecated:: ATCv6
	When support for the Riak backend is removed, support for this tool will
	also be dropped.

The ``traffic_vault_util`` tool - located at :file:`tools/traffic_vault_util.go` in the `Apache Traffic Control repository <https://github.com/apache/trafficcontrol>`_ - is used to view and modify the contents of a Traffic Vault Riak cluster. The tool contains basic operations to display the buckets, keys and values stored within Riak.

.. note:: This tool does not apply to the PostgreSQL Traffic Vault backend.

``traffic_vault_util`` also has a small converter utility to perform a one-off conversion of key formats within the SSL bucket. This conversion is useful when moving from an older version of Traffic Ops to the current version. In the older version, SSL records were indexed by :term:`Delivery Service` database ID. Currently, SSL records are indexed by :term:`Delivery Service` ``xml_id``.

.. program:: traffic_vault_util

Usage
=====
``traffic_vault_util [--dry_run] --vault_ip IP --vault_action ACTION [--vault_user USER] [--vault_password PASSWD] [--vault_port PORT] [--insecure]``

.. option:: --dry_run

	An optional flag which, if given, will cause :program:`traffic_vault_util` to not write changes, but merely print what *would* be done in a real run.

.. option:: --vault_action ACTION

	Defines the action to be performed. Available actions are:

	list_buckets
		Lists the "buckets" in the Riak cluster used by Traffic Vault
	list_keys
		Lists all the keys in all the buckets in the Riak cluster used by Traffic Vault
	list_values
		Lists all the values of all the keys in all the buckets in the Riak cluster used by Traffic Vault
	convert_ssl_to_xmlid
		Changes the key of all records in all buckets that start with "ds" into the ``xml_id`` of the :term:`Delivery Service` for which we assume the record was created.

.. option:: --vault_ip IP

	Either the IP address or :abbr:`FQDN (Fully Qualified Domain Name)` of the Traffic Vault instance with which :program:`traffic_vault_util` will interact.

	.. warning:: If this IP address or :abbr:`FQDN (Fully Qualified Domain Name)` does not point to a real Riak cluster, :program:`traffic_vault_util` will print an error message to STDOUT, but *will* **not** *terminate*. Instead, it will try forever to query the server to which it failed to connect, consuming large amounts of CPU usage all the while\ [1]_.

.. option:: --vault_password PASSWD

	An optional flag used to specify the password of the user defined by :option:`--vault_user` when authenticating with Traffic Vault's Riak cluster.

	.. warning:: Although this flag is optional, the utility will not work without it. It will try, but it will fail\ [1]_.

.. option:: --vault_port PORT

	An optional flag which, if given, sets the port to which :program:`traffic_vault_util` will try to connect to Riak. Default: 8087

.. option:: --vault_user USER

	An optional flag which, if given, specifies the name of the user as whom to connect to Riak.

	.. warning:: Although this flag is optional, the utility will not work without it. It will try, but it will fail\ [1]_.

.. option:: --insecure

	An optional flag which, if given, specifies whether to utilize TLS certificate checks when establishing a connection. Defaults to false.

.. [1] These problems are all tracked by `GitHub Issue #3261 <https://github.com/apache/trafficcontrol/issues/3261>`_.
