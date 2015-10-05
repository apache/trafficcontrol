.. 
.. Copyright 2015 Comcast Cable Communications Management, LLC
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
In order to successfully store private keys you will need to install Riak. 
The latest version of Riak can be downloaded on the Riak `website <http://docs.basho.com/riak/latest/downloads/>`_.
The installation instructions for Riak can be found `here <http://docs.basho.com/riak/latest/ops/building/installing/>`_.

Production is currently running version 2.0.5 of Riak, but the latest version should suffice.


Configuring Traffic Vault
=========================
The following steps were taken to configure Riak in our environments.

Riak configuration file configuration
-------------------------------------

The following steps need to be performed on each Riak server in the cluster:

* Log into riak server as root

* cd to /etc/riak/

* Update the following in riak.conf to reflect your IP:
	- nodename = riak@a-host.sys.kabletown.net
	- listener.http.internal = a-host.sys.kabletown.net:8098 (can be 80 - This endpoint will not work with sec enabled)
	- listener.protobuf.internal = a-host.sys.kabletown.net:8087 (can be different port if you want)
	- listener.https.internal = a-host.sys.kabletown.net:8088 (can be 443)

* Updated the following conf file to point to your cert files
	- ssl.certfile = /etc/riak/certs/server.crt
	- ssl.keyfile = /etc/riak/certs/server.key
	- ssl.cacertfile = /etc/pki/tls/certs/ca-bundle.crt

* Add a line at the bottom of the config for tlsv1
	- tls_protocols.tlsv1 = on

* Once the config file has been updated restart riak
	- ``/etc/init.d/riak restart``

* Validate server is running by going to the following URL:
 	- https://<serverHostname>:8088/ping

Riak-admin configuration
-------------------------

Riak-admin is a command line utility that needs to be run as root on a server in the riak cluster.

Assumptions:
	* Riak 2.0.2 or greater is installed
	* SSL Certificates have been generated (signed or self-signed)
	* Root access to riak servers

Add admin user and riakuser to riak
	* Admin user will be a super user
	* Riakuser will be the application user

Login to one of the riak servers in the cluster as root (any will do) 
	
	1. Enable security
	
		``riak-admin security enable``
	
	2. Add groups
	
		``riak-admin security add-group admins``

		``riak-admin security add-group keysusers``  
	3. Add users 
	 
	 .. Note:: username and password should be stored in /opt/traffic_ops/app/conf/<environment>/riak.conf
	 ..

		``riak-admin security add-user admin password=<AdminPassword> groups=admins``

		``riak-admin security add-user riakuser password=<RiakUserPassword> groups=keysusers`` 

	4. Grant access for admin and riakuser

		``riak-admin security add-source riakuser 0.0.0.0/0 password``

		``riak-admin security add-source admin 0.0.0.0/0 password``

	5. Grant privs to admins for everything

		``riak-admin security grant riak_kv.list_buckets,riak_kv.list_keys,riak_kv.get,riak_kv.put,riak_kv.delete on any to admins``

	6. Grant privs to keysuser for ssl, dnssec, and url_sig_keys buckets only

		``riak-admin security grant riak_kv.get,riak_kv.put,riak_kv.delete on default ssl to keysusers``

		``riak-admin security grant riak_kv.get,riak_kv.put,riak_kv.delete on default dnssec to keysusers``

		``riak-admin security grant riak_kv.get,riak_kv.put,riak_kv.delete on default url_sig_keys to keysusers``

.. seealso:: For more information on security in Riak, see the `Riak Security documentation <http://docs.basho.com/riak/2.0.4/ops/advanced/security/>`_.
.. seealso:: For more information on authentication and authorization in Riak, see the `Riak Authentication and Authorization documentation <http://docs.basho.com/riak/2.0.4/ops/running/authz/>`_.	


Traffic Ops Configuration
-------------------------

There are a couple conifgurations that are necessary in Traffic Ops.

1. Database Updates
	* A new profile for Riak needs to be added to the profile table
	* A new type of Riak needs to be added to the type table
	* The servers in the Riak cluster need to be added to the server table

	 .. Note:: profile and type data should be pre-loaded by seeds sql script.
	 ..

2. Configuration updates
	* /opt/traffic_ops/app/conf/<environment>/riak.conf needs to be updated to reflect the correct username and password for accessing riak.





