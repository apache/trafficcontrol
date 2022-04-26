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

.. _traffic_vault_admin:

****************************
Traffic Vault Administration
****************************

Currently, the supported backends for Traffic Vault are PostgreSQL and Riak, but Riak support is deprecated and may be removed in a future release. More backends may be supported in the future.

.. _traffic_vault_postgresql_backend:

PostgreSQL
==========

In order to use the PostgreSQL backend for Traffic Vault, you will need to set the ``traffic_vault_backend`` option to ``"postgres"`` and include the necessary configuration in the ``traffic_vault_config`` section in :file:`cdn.conf`. The ``traffic_vault_config`` options for the PostgreSQL backend are as follows:

:dbname:                    The name of the database to use
:hostname:                  The hostname of the database server to connect to
:password:                  The password to use when connecting to the database
:port:                      The port number that the database listens for new connections on (NOTE: the PostgreSQL default is 5432)
:user:                      The username to use when connecting to the database
:aes_key_location:          The location on-disk for a base64-encoded AES key used to encrypt secrets before they are stored. It is highly recommended to backup this key to a safe, secure storage location, because if it is lost, you will lose access to all your Traffic Vault data. Either this option or ``hashicorp_vault`` must be used.
:hashicorp_vault:           This group of configuration options is for fetching the base64-encoded AES key from `HashiCorp Vault <https://www.vaultproject.io/>`_. This uses the `AppRole authentication method <https://learn.hashicorp.com/tutorials/vault/approle>`_.

	:address:     The address of the HashiCorp Vault server, e.g. http://localhost:8200
	:role_id:     The RoleID of the AppRole.
	:secret_id:   The SecretID issued against the AppRole.
	:secret_path: The URI path where the secret AES key is located, e.g. /v1/secret/data/trafficvault. The secret should be stored using the `KV Secrets Engine <https://www.vaultproject.io/docs/secrets/kv>`_ with a key of ``traffic_vault_key`` and value of a base64-encoded AES key, e.g. ``traffic_vault_key='WoFc86CisM1aXo8D5GvDnq2h9kjULuIP4upaqX15SRc='``.
	:login_path:  Optional. The URI path used to login with the AppRole method. Default: /v1/auth/approle/login
	:timeout_sec: Optional. The timeout (in seconds) for requests. Default: 30
	:insecure:    Optional. Disable server certificate verification. This should only be used for testing purposes. Default: false

:conn_max_lifetime_seconds: Optional. The maximum amount of time (in seconds) a connection may be reused. If negative, connections are not closed due to a connection's age. If 0 or unset, the default of 60 is used.
:max_connections:           Optional. The maximum number of open connections to the database. Default: 0 (unlimited)
:max_idle_connections:      Optional. The maximum number of connections in the idle connection pool. If negative, no idle connections are retained. If 0 or unset, the default of 30 is used.
:query_timeout_seconds:     Optional. The duration (in seconds) after which database queries will time out and be cancelled. Default: 30
:ssl:                       Optional. Whether or not to use SSL to connect to the database. Default: false

Example cdn.conf snippet:
-------------------------

.. code-block:: json

	{
		"traffic_ops_golang": {
			"traffic_vault_backend": "postgres",
			"traffic_vault_config": {
				"dbname": "tv_development",
				"hostname": "localhost",
				"user": "traffic_vault",
				"password": "twelve",
				"port": 5432,
				"ssl": false,
				"conn_max_lifetime_seconds": 60,
				"max_connections": 500,
				"max_idle_connections": 30,
				"query_timeout_seconds": 10,
				"aes_key_location": "/opt/traffic_ops/app/conf/tv.key"
			}
		}
	}

Administration of the PostgreSQL database for Traffic Vault
-----------------------------------------------------------

Similar to administering the Traffic Ops database, the :ref:`admin <database-management>` tool should be used for administering the PostgreSQL Traffic Vault backend.

.. program:: reencrypt

app/db/reencrypt/reencrypt
--------------------------
The :program:`reencrypt` binary is used to re-encrypt all data in the Postgres Traffic Vault with a new base64-encoded AES key.

.. note:: For proper resolution of configuration files, it's recommended that this binary be run from the ``app/db/reencrypt`` directory.

Usage
"""""
``./reencrypt [options]``

Options and Arguments
"""""""""""""""""""""
.. option:: --new-key NEW_KEY

	(Optional) The file path for the new base64-encoded AES key. Default is ``/opt/traffic_ops/app/conf/new.key``.

.. option:: --previous-key PREVIOUS_KEY

	(Optional) The file path for the previous base64-encoded AES key. Default is ``/opt/traffic_ops/app/conf/aes.key``.

.. option:: --cfg CONFIG_FILE

	(Optional) The path for the configuration file. Default is ``./reencrypt.conf``.

.. option:: --help

	(Optional) Print usage information and exit.

.. code-block:: bash
	:caption: Example Usage

	./reencrypt --new-key ~/exampleNewKey.txt --previous-key ~/exampleOldKey.txt

reencrypt.conf
""""""""""""""
This file deals with configuration of the Traffic Vault Database to be used with the :program:`reencrypt` tool.

:dbname: The name of the PostgreSQL database used.
:hostname: The hostname (:abbr:`FQDN (Fully Qualified Domain Name)`) of the server that runs the Traffic Vault Database.
:password: The password to use when authenticating with the Traffic Vault database.
:port: The port number on which the Traffic Vault Database is listening for incoming connections (NOTE: the PostgreSQL default is 5432).
:ssl: A boolean that sets whether or not the Traffic Vault Database encrypts its connections with SSL.
:user: The name of the user as whom to connect to the database.


.. _traffic_vault_riak_backend:

Riak (deprecated)
=================

.. deprecated:: ATCv6
	The Riak Traffic Vault backend is deprecated and support may be removed in a future release. It is highly recommended to use the PostgreSQL Traffic Vault backend instead.

In order to use the Riak backend for Traffic Vault, you will need to set the ``traffic_vault_backend`` option to ``"riak"`` and include the necessary configuration in the ``traffic_vault_config`` section in :file:`cdn.conf`. The ``traffic_vault_config`` options for the Riak backend are as follows:

:password:      The password to use when authenticating with Riak
:user:          The username to use when authenticating with Riak
:port:          The Riak protobuf port to connect to. Default: 8087
:tlsConfig:     Optional. Certain TLS options from `the tls.Config struct options <https://golang.org/pkg/crypto/tls/#Config>`_ may be included here, such as ``insecureSkipVerify: true`` to disable certificate validation in order to use self-signed certificates for test/development purposes.
:MaxTLSVersion: Optional. This is the highest TLS version that Traffic Ops is allowed to use to connect to Traffic Vault. Valid values are "1.0", "1.1", "1.2", and "1.3". The default is "1.1".

.. note:: Enabling TLS 1.1 in Riak itself is required for Traffic Ops to communicate with Riak. See :ref:`Enabling TLS 1.1 <tv-admin-enable-tlsv1.1>` for details.

Example cdn.conf snippet:
-------------------------

.. code-block:: json

	{
		"traffic_ops_golang": {
			"traffic_vault_backend": "riak",
			"traffic_vault_config": {
				"user": "riakuser",
				"password": "password",
				"MaxTLSVersion": "1.1",
				"port": 8087
			}
		}
	}

Installing the Riak backend for Traffic Vault
---------------------------------------------
In order to successfully store private keys you will need to install Riak. The latest version of Riak can be downloaded on `the Riak website <https://docs.riak.com/riak/latest/downloads/>`_. The installation instructions for Riak can be found `here <https://docs.riak.com/riak/kv/latest/setup/installing/index.html>`__. Based on experience, version 2.0.5 of Riak is recommended, but the latest version should suffice.

Configuring Riak
----------------
Follow these steps to configure Riak in a production environment.

Self Signed Certificate configuration
"""""""""""""""""""""""""""""""""""""
.. note:: Self-signed certificates are not recommended for production use. Intended for development or learning purposes only. Modify subject as necessary.

.. code-block:: shell
	:caption: Self-Signed Certificate Configuration

	cd ~
	mkdir certs
	cd certs
	openssl genrsa -out ca-bundle.key 2048
	openssl req -new -key ca-bundle.key -out ca-bundle.csr -subj "/C=US/ST=CO/L=DEN/O=somecompany/OU=CDN/CN=somecompany.net/emailAddress=someuser@somecompany.net"
	openssl x509 -req -days 365 -in ca-bundle.csr -signkey ca-bundle.key -out ca-bundle.crt
	openssl genrsa -out server.key 2048
	openssl req -new -key server.key -out server.csr -subj "/C=US/ST=CO/L=DEN/O=somecompany/OU=CDN/CN=somecompany.net/emailAddress=someuser@somecompany.net"
	openssl x509 -req -days 365 -in server.csr -CA ca-bundle.crt -CAkey ca-bundle.key -CAcreateserial -out server.crt
	mkdir /etc/riak/certs
	mv -f server.crt /etc/riak/certs/.
	mv -f server.key /etc/riak/certs/.
	mv -f ca-bundle.crt /etc/pki/tls/certs/.


Riak Configuration File
"""""""""""""""""""""""
The following steps need to be performed on each Riak server in the cluster:

#. Log into Riak server as root
#. Update the following in :file:`riak.conf` to reflect your IP, hostname, and CDN domains and sub-domains:

	* ``nodename = riak@a-host.sys.kabletown.net``
	* ``listener.http.internal = a-host.sys.kabletown.net:8098`` (port can be 80 - This endpoint will not work over HTTPS)
	* ``listener.protobuf.internal = a-host.sys.kabletown.net:8087`` (can be different port if you want)
	* ``listener.https.internal = a-host.sys.kabletown.net:8088`` (port can be 443)

#. Update the following in :file:`riak.conf` file to point to your SSL certificate files

	- ``ssl.certfile = /etc/riak/certs/server.crt``
	- ``ssl.keyfile = /etc/riak/certs/server.key``
	- ``ssl.cacertfile = /etc/pki/tls/certs/ca-bundle.crt``

.. _tv-admin-enable-tlsv1.1:

Enabling TLS 1.1 (required)
'''''''''''''''''''''''''''

#. Add a line at the bottom of the :file:`riak.conf` for TLSv1.1 by setting ``tls_protocols.tlsv1.1 = on``
#. Once the configuration file has been updated restart Riak
#. Consult the `Riak documentation <https://docs.riak.com/riak/kv/latest/setup/installing/verify/>`_ for instructions on how to verify the installed service

``riak-admin`` Configuration
""""""""""""""""""""""""""""
``riak-admin`` is a command line utility used to configure Riak that needs to be run as root on a server in the Riak cluster.

.. seealso:: `The riak-admin documentation <https://docs.riak.com/riak/kv/latest/using/admin/riak-admin/>`_

.. code-block:: shell
	:caption: Traffic Vault Setup with ``riak-admin``

	# This script need only be run on any *one* Riak server in the cluster

	# Enable security and secure access groups
	riak-admin security enable
	riak-admin security add-group admins
	riak-admin security add-group keysusers

	# User name and password should be stored in the traffic_vault_config section in
	# /opt/traffic_ops/app/conf/cdn.conf on the Traffic Ops server (with traffic_vault_backend = riak)
	# In this example, we assume the usernames 'admin' and 'riakuser' with
	# respective passwords stored in the ADMIN_PASSWORD and RIAK_USER_PASSWORD
	# environment variables
	riak-admin security add-user admin password=$ADMIN_PASSWORD groups=admins
	riak-admin security add-user riakuser password=$RIAK_USER_PASSWORD groups=keysusers
	riak-admin security add-source riakuser 0.0.0.0/0 password
	riak-admin security add-source admin 0.0.0.0/0 password

	# Grant privileges to the admins group for everything
	riak-admin security grant riak_kv.list_buckets,riak_kv.list_keys,riak_kv.get,riak_kv.put,riak_kv.delete on any to admins

	# Grant privileges to keysusers group for SSL, DNSSEC, and url_sig_keys buckets only
	riak-admin security grant riak_kv.get,riak_kv.put,riak_kv.delete on default ssl to keysusers
	riak-admin security grant riak_kv.get,riak_kv.put,riak_kv.delete on default dnssec to keysusers
	riak-admin security grant riak_kv.get,riak_kv.put,riak_kv.delete on default url_sig_keys to keysusers
	riak-admin security grant riak_kv.get,riak_kv.put,riak_kv.delete on default cdn_uri_sig_keys to keysusers

.. seealso:: For more information on security in Riak, see the `Riak Security documentation <https://docs.riak.com/riak/kv/latest/using/security/index.html>`_.


Traffic Ops Configuration
"""""""""""""""""""""""""
Before a fully set-up Riak instance may be used as the Traffic Vault backend, it must be added as a server to Traffic Ops. The easiest way to accomplish this is via Traffic Portal at :menuselection:`Configure --> Servers`, though :ref:`to-api-servers` may also be used by low-level tools and/or scripts. The Traffic Ops configuration file :file:`/opt/traffic_ops/app/conf/cdn.conf` must be updated to set ``traffic_vault_backend`` to ``"riak"`` and the ``traffic_vault_config`` to include the correct username and password for accessing the Riak database.

Configuring Riak Search
-----------------------
In order to more effectively support retrieval of SSL certificates by Traffic Router and :term:`ORT`, the Riak backend for Traffic Vault uses `Riak search <https://docs.riak.com/riak/kv/latest/using/reference/search/>`_. Riak Search uses `Apache Solr <https://lucene.apache.org/solr>`_ for indexing and searching of records. This section explains how to enable, configure, and validate Riak Search.

Riak Configuration
""""""""""""""""""
On each Traffic Vault server follow these steps.

#. If Java (JDKv1.8+) is not already installed on your Riak server, install Java

	.. code-block:: shell
		:caption: Check if Java is Installed, Then Install if Needed

		# Ensure that this outputs a Java version that is at least 1.8
		java -version

		# If it didn't, or produced an error because `java` doesn't exist,
		# install the correct version
		# (OpenJDK is used here because of its permissive license, though OracleJDK
		# should work with some tinkering)

		# On CentOS/RedHat/Fedora (recommended)
		yum install -y java-1.8.0-openjdk java-1.8.0-openjdk-devel

		# On Ubuntu/Debian/Linux Mint
		apt install -y openjdk-8-jdk

		# Arch/Manjaro
		pacman -Sy jdk8-openjdk

#. Enable search in :file:`riak.conf` by changing the ``search = off`` setting to ``search = on``
#. Restart Riak to propagate configuration changes

	.. code-block:: bash
		:caption: Restarting Riak on :manpage:`systemd(1)` Systems

		systemctl restart riak

One-time Configuration
''''''''''''''''''''''
After Riak has been configured to use Riak Search, permissions still need need to be updated to allow users to utilize this feature. Unlike actually setting up Riak Search, the permissions step need only be done on any *one* of the Riak servers in the cluster.

#. Use ``riak-admin`` to grant ``search.admin`` permissions to the "admin" user and ``search.query`` permissions to **both** the "admin" user and the "riakuser" user. The "admin" user will also require ``search.admin`` permissions on the ``schema`` (in addition to ``index``) and ``riak_core.set_bucket`` permissions on ``any``.

	.. code-block:: bash
		:caption: Setting up Riak Search Permissions

		riak-admin security grant search.admin on schema to admin
		riak-admin security grant search.admin on index to admin
		riak-admin security grant search.query on index to admin
		riak-admin security grant search.query on index sslkeys to admin
		riak-admin security grant search.query on index to riakuser
		riak-admin security grant search.query on index sslkeys to riakuser
		riak-admin security grant riak_core.set_bucket on any to admin

#. Add the search schema to Riak. This schema is a simple Apache Solr configuration file which will index all records on CDN, hostname, and :term:`Delivery Service`. The file can be found at :file:`traffic_ops/app/config/misc/riak_search/sslkeys.xml` in the Traffic Control repository.

	.. code-block:: bash
		:caption: Adding the GitHub-hosted Search Schema to Riak

		# Obtain the configuration file - in this example by downloading it from GitHub
		wget https://raw.githubusercontent.com/apache/trafficcontrol/master/traffic_ops/app/conf/misc/riak_search/sslkeys.xml

		# Upload the schema to the Riak server using its API
		# Note that the assumptions made here are that the "admin" user's password is "pass"
		# and the server is accessible at port 8088 on the hostname "trafficvault.infra.ciab.test"
		curl --tlsv1.1 --tls-max 1.1 -kvsX PUT "https://admin:pass@trafficvault.infra.ciab.test:8088/search/schema/sslkeys" -H "Content-Type: application/xml" -d @sslkeys.xml

#. Add the search index to Riak.

	.. code-block:: bash
		:caption: Adding the Search Index to Riak Via its API

		# Note that the assumptions made here are that the "admin" user's password is "pass"
		# and the server is accessible at port 8088 on the hostname "trafficvault.infra.ciab.test"
		curl --tlsv1.1 --tls-max 1.1 -kvsX PUT "https://admin:pass@trafficvault.infra.ciab.test:8088/search/index/sslkeys" -H 'Content-Type: application/json' -d '{"schema":"sslkeys"}'

#. Associate the ``sslkeys`` index to the ``ssl`` bucket in Riak

	.. code-block:: bash
		:caption: Using the Riak API to Create an Index-to-Bucket Association for ``sslkeys``

		# Note that the assumptions made here are that the "admin" user's password is "pass"
		# and the server is accessible at port 8088 on the hostname "trafficvault.infra.ciab.test"
		curl --tlsv1.1 --tls-max 1.1 -kvs -XPUT "https://admin:pass@trafficvault.infra.ciab.test:8088/buckets/ssl/props" -H'content-type:application/json' -d'{"props":{"search_index":"sslkeys"}}'

#. To validate the search is working run a query against the Riak database server, or use the Traffic Ops API endpoint: :ref:`to-api-cdns-name-name-sslkeys`

	.. code-block:: bash
		:caption: Validate Riak Search is Working

		# Note that the assumptions made here are that the "admin" user's password is
		# "pass", the Traffic Vault server's Riak database is accessible at port 8088 on
		# the hostname "trafficvault.infra.ciab.test", $COOKIE contains a valid
		# Mojolicious cookie for a Traffic Ops user with proper permissions, and the
		# Traffic Ops server is available at the hostname "trafficops.infra.ciab.test"

		# Verify by querying Riak directly
		curl --tlsv1.1 --tls-max 1.1 -kvs "https://admin:password@trafficvault.infra.ciab.test:8088/search/query/sslkeys?wt=json&q=cdn:CDN-in-a-Box"

		# Verify using the Traffic Ops API
		curl -Lvs -H "Cookie: $COOKIE" https://trafficops.infra.ciab.test/api/4.0/cdns/name/mycdn/sslkeys
