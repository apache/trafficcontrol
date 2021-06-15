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
``traffic_vault_migrate [-cdhmr] [-f value] [-g value] [-o value] [-t value]``

.. option:: -c, --compare

		Compare 'to' and 'from' backend keys. Will fetch keys from the dbs of both 'to' and 'from', sorts them by cdn/ds/version and does a deep comparison.

		.. note:: Mutually exclusive with :option:`-r`/:option:`--dry`

.. option:: -d, --dump

		Write keys (from 'from' server) to disk in the folder 'dump' with the unix permissions 0640.

		.. warning:: This can write potentially sensitive information to disk, use with care.

		.. note:: Mutually exclusive with :option:`-l`/:option:`--fill`

.. option:: -e LEVEL, --logLevel=LEVEL

		Print everything at above specified log level (error|warning|info|debug|event) [info]

		.. note:: Mutually exclusive with :option:`-l`/:option:`--logCfg`

.. option:: -f CFG, --fromCfgPath=CFG

		From server config file ["riak.json"]

.. option:: -g CFG, --toCfgPath=CFG

		To server config file ["pg.json"]

.. option:: -h, --help

		Displays usage information

.. option:: -i DIR, --fill DIR

		Insert data into `to` server with data this directory

		.. note:: Mutually exclusive with :option:`-d`/:option:`--dump`

.. option:: -l CFG, --logCfg CFG

		Log configuration file

		.. note:: Mutually exclusive with :option:`-e`/:option:`--logLevel`

.. option:: -o TYPE, --toType=TYPE

		From server types (Riak|PG) [PG]

.. option:: -m, --noConfirm

		Do not require confirmation before inserting records

.. option:: -r, --dry

		Do not perform writes. Will do a basic output of the keys on the 'from' backend.

		.. note:: Mutually exclusive with :option:`-c`/:option:`--compare`

.. option:: -t TYPE, --fromType=TYPE

		From server types (Riak|PG) [Riak]


Riak
----------

riak.json
""""""""""

 :user: The username used to log into the Riak server.

 :password: The password used to log into the Riak server.

 :host: The hostname for the Riak server.

 :port: The port for which the Riak server is listening for protobuf connections.

 :timeout: The number of seconds to wait for each operation.

 :insecure: (Optional) Determines whether to verify insecure certificates.

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

 :sslmode: The ssl settings for the client connection, `explanation here <https://www.postgresql.org/docs/13/libpq-ssl.html#LIBPQ-SSL-SSLMODE-STATEMENTS>`_. Options are 'disable', 'allow', 'prefer', 'require', 'verify-ca' and 'verify-full'

 :aesKey: The base64 encoding of a 16, 24, or 32 bit AES key.


Logging
----------

The log configuration file has the structure:

 :error_log: Where to output error messages (stderr|stdout|null)

 :warning_log: Where to output warning messages (stderr|stdout|null)

 :info_log: Where to output info messages (stderr|stdout|null)

 :debug_log: Where to output error messages (stderr|stdout|null)

 :event_log: Where to output error messages (stderr|stdout|null)

Development
=============
To add a plugin, implement the traffic_vault_migrate.go:TVBackend interface and add the backend to the returned values in :atc-godoc:`tools/traffic_vault_migrate.supportBackends`.
