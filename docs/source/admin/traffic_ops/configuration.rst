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

.. _`regex_revalidate plugin`: https://docs.trafficserver.apache.org/en/7.1.x/admin-guide/plugins/regex_revalidate.en.html

*************************
Traffic Ops - Configuring
*************************
Follow the steps below to configure the newly installed Traffic Ops Instance.

Installing the SSL Certificate
==============================
By default, Traffic Ops runs as an SSL web server (that is, over HTTPS), and a certificate needs to be installed.

Self-signed Certificate (Development)
-------------------------------------
.. code-block:: console
	:caption: Example Procedure

	$ openssl genrsa -des3 -passout pass:x -out localhost.pass.key 2048
	Generating RSA private key, 2048 bit long modulus
	...
	$ openssl rsa -passin pass:x -in localhost.pass.key -out localhost.key
	writing RSA key
	$ rm localhost.pass.key

	$ openssl req -new -key localhost.key -out localhost.csr
	You are about to be asked to enter information that will be incorporated
	into your certificate request.
	What you are about to enter is what is called a Distinguished Name or a DN.
	There are quite a few fields but you can leave some blank
	For some fields there will be a default value,
	If you enter '.', the field will be left blank.
	-----
	Country Name (2 letter code) [XX]:US<enter>
	State or Province Name (full name) []:CO<enter>
	Locality Name (eg, city) [Default City]:Denver<enter>
	Organization Name (eg, company) [Default Company Ltd]: <enter>
	Organizational Unit Name (eg, section) []: <enter>
	Common Name (eg, your name or your server's hostname) []: <enter>
	Email Address []: <enter>

	Please enter the following 'extra' attributes
	to be sent with your certificate request
	A challenge password []: pass<enter>
	An optional company name []: <enter>
	$ openssl x509 -req -sha256 -days 365 -in localhost.csr -signkey localhost.key -out localhost.crt
	Signature ok
	subject=/C=US/ST=CO/L=Denver/O=Default Company Ltd
	Getting Private key
	$ sudo cp localhost.crt /etc/pki/tls/certs
	$ sudo cp localhost.key /etc/pki/tls/private
	$ sudo chown trafops:trafops /etc/pki/tls/certs/localhost.crt
	$ sudo chown trafops:trafops /etc/pki/tls/private/localhost.key

Certificate from Certificate Authority (Production)
---------------------------------------------------

.. Note:: You will need to know the appropriate answers when generating the certificate request file :file:`trafficopss.csr` below.

Example Procedure
"""""""""""""""""
.. code-block:: console
	:caption: Example Procedure

	$ openssl genrsa -des3 -passout pass:x -out trafficops.pass.key 2048
	Generating RSA private key, 2048 bit long modulus
	...
	$ openssl rsa -passin pass:x -in trafficops.pass.key -out trafficops.key
	writing RSA key
	$ rm localhost.pass.key

Generate the :abbr:`CSR (Certificate Signing Request)` file needed for :abbr:`CA (Certificate Authority)` request

.. code-block:: console
	:caption: Example Certificate Signing Request File Generation

	$ openssl req -new -key trafficops.key -out trafficops.csr
	You are about to be asked to enter information that will be incorporated
	into your certificate request.
	What you are about to enter is what is called a Distinguished Name or a DN.
	There are quite a few fields but you can leave some blank
	For some fields there will be a default value,
	If you enter '.', the field will be left blank.
	-----
	Country Name (2 letter code) [XX]: <enter country code>
	State or Province Name (full name) []: <enter state or province>
	Locality Name (eg, city) [Default City]: <enter locality name>
	Organization Name (eg, company) [Default Company Ltd]: <enter organization name>
	Organizational Unit Name (eg, section) []: <enter organizational unit name>
	Common Name (eg, your name or your server's hostname) []: <enter server's hostname name>
	Email Address []: <enter e-mail address>

	Please enter the following 'extra' attributes
	to be sent with your certificate request
	A challenge password []: <enter challenge password>
	An optional company name []: <enter>
	$ sudo cp trafficops.key /etc/pki/tls/private
	$ sudo chown trafops:trafops /etc/pki/tls/private/trafficops.key

You must then take the output file :file:`trafficops.csr` and submit a request to your :abbr:`CA (Certificate Authority)`. Once you get approved and receive your :file:`trafficops.crt` file

.. code-block:: shell
	:caption: Certificate Installation

	sudo cp trafficops.crt /etc/pki/tls/certs
	sudo chown trafops:trafops /etc/pki/tls/certs/trafficops.crt

If necessary, install the :abbr:`CA (Certificate Authority) certificate's ``.pem`` and ``.crt`` files in ``/etc/pki/tls/certs``.

You will need to update the file :file:`/opt/traffic_ops/app/conf/cdn.conf` with the any necessary changes.

.. code-block:: text
	:caption: Sample 'listen' Line When Path to ``trafficops.crt`` and ``trafficops.key`` are Known

	'hypnotoad' => ...
	    'listen' => 'https://[::]:443?cert=/etc/pki/tls/certs/trafficops.crt&key=/etc/pki/tls/private/trafficops.key&ca=/etc/pki/tls/certs/localhost.ca&verify=0x00&ciphers=AES128-GCM-SHA256:HIGH:!RC4:!MD5:!aNULL:!EDH:!ED'
		 ...


Regions, Locations and Cache Groups
===================================
All servers have to have a :term:`Physical Location`, which defines their geographic latitude and longitude. Each :term:`Physical Location` is part of a :term:`Region`, and each :term:`Region` is part of a :term:`Division`. For example, ``Denver`` could be the name of a :term:`Physical Location` in the ``Mile High`` :term:`Region` and that :term:`Region` could be part of the ``West`` :term:`Division`. The hierarchy between these terms is illustrated graphically in :ref:`topography-hierarchy`.

.. _topography-hierarchy:
.. figure:: images/topography.*
	:align: center
	:alt: A graphic illustrating the hierarchy exhibited by topological groupings
	:figwidth: 25%

	Topography Hierarchy

To create these structures in Traffic Portal, first make at least one :term:`Division` under :menuselection:`Topology --> Divisions`. Next enter the desired :term:`Region`\ (s) in :menuselection:`Topology --> Regions`, referencing the earlier-entered :term:`Division`\ (s). Finally, enter the desired :term:`Physical Location`\ (s) in :menuselection:`Topology --> Phys Locations`, referencing the earlier-entered :term:`Region`\ (s).

All servers also have to be part of a :term:`Cache Group`. A :term:`Cache Group` is a logical grouping of cache servers, that don't have to be in the same :term:`Physical Location` (in fact, usually a :term:`Cache Group` is spread across minimally two :term:`Physical Location`\ s for redundancy purposes), but share geographical coordinates for content routing purposes.

Configuring Content Purge
=========================
Purging cached content using :abbr:`ATS (Apache Traffic Server)` is not simple; there is no file system from which to delete files and/or directories, and in large caches it can be hard to delete content matching a simple regular expression from the cache. This is why Traffic Control uses the `Regex Revalidate Plugin <https://docs.trafficserver.apache.org/en/7.1.x/admin-guide/plugins/regex_revalidate.en.html>`_ to purge content from the cache. The cached content is not actually removed, instead a check that runs before each request on each cache server is serviced to see if this request matches a list of regular expressions. If it does, the cache server is forced to send the request upstream to its parents (possibly other caches, possibly the origin) without checking for the response in its cache. The Regex Revalidate Plugin will monitor its configuration file, and will pick up changes to it without needing to alert :abbr:`ATS (Apache Traffic Server). Changes to this file need to be distributed to the highest tier (Mid-tier) cache servers in the CDN before they are distributed to the lower tiers, to prevent filling the lower tiers with the content that should be purged from the higher tiers without hitting the origin. This is why the :term:`ORT` script will - by default - push out configuration changes to Mid-tier cache servers first, confirm that they have all been updated, and then push out the changes to the lower tiers. In large CDNs, this can make the distribution and time to activation of the purge too long, and because of that there is the option to not distribute the ``regex_revalidate.config`` file using the :term:`ORT` script, but to do this using other means. By default, Traffic Ops will use :term:`ORT` to distribute the ``regex_revalidate.config`` file.

.. _Creating-CentOS-Kickstart:

Creating the CentOS Kickstart File
==================================
The Kickstart file is a text file, containing a list of items, each identified by a keyword. This file can be generated using the `Red Hat Kickstart Configurator application <https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/5/html/installation_guide/ch-redhat-config-kickstart>`_, or it can be written from scratch. The Red Hat Enterprise Linux installation program also creates a sample Kickstart file based on the options selected during installation. It is written to the file :file:`/root/anaconda-ks.cfg` in this case. This file is editable using most text editors.

Generating a System Image
-------------------------
#. Create a Kickstart file.
#. Create a boot media with the Kickstart file or make the Kickstart file available on the network.
#. Make the installation tree available.
#. Start the Kickstart installation.

.. code-block:: shell
	:caption: Creating a New System Image Definition Tree from an Existing One

	# Starting from the Kickstart root directory (`/var/www/files` by default)
	mkdir newdir
	cd newdir/

	# In this example, the pre-existing system image definition tree is for CentOS 7.4 located in `centos74`
	cp -r ../centos74/* .
	vim ks.src
	vim isolinux/isolinux.cfg
	cd ..
	vim osversions.cfg

:file:`ks.src` is a standard, Kickstart-formatted file that the will be used to create the Kickstart (ks.cfg) file for the install whenever a system image is generated from the source tree. :file:`ks.src` is a template - it will be overwritten by any information set in the form submitted from :menuselection:`Tools --> Generate ISO` in Traffic Portal. Ultimately, the two are combined to create the final Kickstart file (:file:`ks.cfg`).

.. Note:: It is highly recommended for ease of use that the system image source trees be kept under 1GB in size.

.. seealso:: For in-depth instructions, please see `Kickstart Installation <https://access.redhat.com/documentation/en-US/Red_Hat_Enterprise_Linux/6/html/Installation_Guide/s1-kickstart2-howuse.html>`_ in the Red Hat documentation.


Configuring the Go Application
==============================
Traffic Ops is in the process of migrating from Perl to Go, and currently runs as two applications. The Go application serves all endpoints which have been rewritten in the Go language, and transparently proxies all other requests to the old Perl application. Both applications are installed by the RPM, and both run as a single :manpage:`systemd(1)` service. When the project has fully migrated to Go, the Perl application will be removed, and the RPM and service will consist solely of the Go application.

By default, the :program:`postinstall` script configures the Go application to behave and transparently serve as the old Perl Traffic Ops did in previous versions. This includes reading the old ``cdn.conf`` and ``database.conf`` config files, and logging to the old ``access.log`` location. However, the Go Traffic Ops application may be customized by passing the command-line flag, ``-oldcfg=false``. By default, it will then look for a configuration file at :file:`/opt/traffic_ops/conf/traffic_ops_golang.config`. The new configuration file location may also be customized via the ``-cfg`` flag. A sample configuration file is installed by the RPM at :file:`/opt/traffic_ops/conf/traffic_ops_golang.config`. The new Go Traffic Ops application as a :manpage:`systemd(1)` service with a new configuration file, the ``-oldcfg=false`` and  ``-cfg`` flags may be added to the ``start`` function in the service file, located by default at :file:`/etc/init.d/traffic_ops`.
