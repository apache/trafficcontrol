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

.. _traffic_vault_migrate:

#########################
Traffic Vault Migrate
#########################
The ``traffic_vault_migrate`` tool - located at :file:`tools/traffic_vault_migrate/traffic_vault_migrate.go` in the `Apache Traffic Control repository <https://github.com/apache/trafficcontrol>`_ -
is used to transfer TV keys between database servers. It interfaces directly with each backend so Traffic Ops/Vault being available is not a requirement.
The tool assumes that the schema for each backend is already setup as according to the :ref:`admin setup <traffic_vault_admin>`.

.. program:: traffic_vault_migrate

Usage
===========
``traffic_vault_migrate -from_cfg CFG -to_cfg CFG -from_type TYP -to_type TYP [-confirm] [-compare] [-dry] [-dump]``

.. option:: -compare

		Compare 'to' and 'from' backend keys. Will fetch keys from the dbs of both 'to' and 'from', sorts them by cdn/ds/version and does a deep comparison.

.. option:: -confirm

		Requires confirmation before inserting records (default true)

.. option:: -dry

		Do not perform writes. Will do a basic output of the keys on the 'from' backend.

.. option:: -dump

		Write keys (from 'from' server) to disk in the folder 'dump' with the unix permissions 0640.

		.. warning:: This can write potentially sensitive information to disk, use with care.

.. option:: -from_cfg

		From server config file (default "riak.json")

.. option:: -from_type

		From server types (Riak|PG) (default "Riak")

.. option:: -to_cfg

		To server config file (default "pg.json")

.. option:: -to_type

		From server types (Riak|PG) (default "PG")

Riak
----------

riak.json
""""""""""

 :user: The username used to log into the Riak server.

 :password: The password used to log into the Riak server.

 :host: The hostname for the Riak server.

 :port: The port for which the Riak server is listening for protobuf connections.

 :tls: (Optional) Determines whether to verify insecure certificates.

 :tlsVersion: (Optional) Max TLS version supported. Valid values are  "10", "11", "12", "13".


Postgres
---------
:program:`traffic_vault_migrate` will properly handle both encryption and decryption of postgres data as that is done on the client side.

pg.json
"""""""""

 :user: The username used to log into the PG server.

 :password: The password for the user to log into the PG server.

 :database: The database to connect to.

 :port: The port on which the PG server is listening.

 :host: The hostname of the PG server.

 :sslmode: The ssl settings for the client connection, `explanation here <https://www.postgresql.org/docs/9.1/libpq-ssl.html#LIBPQ-SSL-SSLMODE-STATEMENTS>`_. Options are 'disable', 'allow', 'prefer', 'require', 'verify-ca' and 'verify-full'

 :aesKey: The base64 encoding of a 16, 24, or 32 bit AES key.

Development
=============
To add a plugin, implement the traffic_vault_migrate.go:TVBackend interface and add the backend to the returned values in supportedBackends
