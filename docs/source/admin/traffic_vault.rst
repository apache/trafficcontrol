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

****************************
Traffic Vault Administration
****************************
Installing Traffic Vault
========================
In order to successfully store private keys you will need to install Riak. The latest version of Riak can be downloaded on `the Riak website <https://docs.riak.com/riak/latest/downloads/>`_. The installation instructions for Riak can be found `here <https://docs.riak.com/riak/kv/latest/setup/installing/index.html>`__. Based on experience, version 2.0.5 of Riak is recommended, but the latest version should suffice.

Configuring Traffic Vault
=========================
The following steps were taken to configure Riak in Comcast production environments.

Self Signed Certificate configuration
-------------------------------------
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
-----------------------
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
"""""""""""""""""""""""""""

#. Add a line at the bottom of the :file:`riak.conf` for TLSv1.1 by setting ``tls_protocols.tlsv1.1 = on``
#. Once the configuration file has been updated restart Riak
#. Consult the `Riak documentation <https://docs.riak.com/riak/kv/latest/setup/installing/verify/>`_ for instructions on how to verify the installed service

``riak-admin`` Configuration
----------------------------
``riak-admin`` is a command line utility used to configure Riak that needs to be run as root on a server in the Riak cluster.

.. seealso:: `The riak-admin documentation <https://docs.riak.com/riak/kv/latest/using/admin/riak-admin/>`_

.. code-block:: shell
	:caption: Traffic Vault Setup with ``riak-admin``

	# This script need only be run on any *one* Riak server in the cluster

	# Enable security and secure access groups
	riak-admin security enable
	riak-admin security add-group admins
	riak-admin security add-group keysusers

	# User name and password should be stored in
	# /opt/traffic_ops/app/conf/<environment>/riak.conf on the Traffic Ops
	# server
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
-------------------------
Before a fully set-up Traffic Vault instance may be used, it must be added as a server to Traffic Ops. The easiest way to accomplish this is via Traffic Portal at :menuselection:`Configure --> Servers`, though :ref:`to-api-servers` may also be used by low-level tools and/or scripts. The Traffic Ops configuration file :file:`/opt/traffic_ops/app/conf/{environment}/riak.conf` for the appropriate environment must also be updated to reflect the correct username and password for accessing the Riak database.

Configuring Riak Search
=======================
In order to more effectively support retrieval of SSL certificates by Traffic Router and :term:`ORT`, Traffic Vault uses `Riak search <https://docs.riak.com/riak/kv/latest/using/reference/search/>`_. Riak Search uses `Apache Solr <https://lucene.apache.org/solr>`_ for indexing and searching of records. This section explains how to enable, configure, and validate Riak Search.

Riak Configuration
------------------
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
""""""""""""""""""""""
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
		curl -Lvs -H "Cookie: $COOKIE" https://trafficops.infra.ciab.test/api/2.0/cdns/name/mycdn/sslkeys
