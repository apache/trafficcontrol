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

.. _anonymous_blocking-qht:

****************************
Configure Anonymous Blocking
****************************

.. Note:: Anonymous Blocking is only supported for HTTP delivery services. You will need access to a database that provides anonymous IP statistics (`Maxmind's database <https://www.maxmind.com/en/solutions/geoip2-enterprise-product-suite/anonymous-ip-database>`_ is recommended, as this functionality was built specifically to work with it.)

#. Prepare the Anonymous Blocking configuration file. Anonymous Blocking uses a configuration file in JSON format to define blocking rules for :term:`Delivery Services`. The file needs to be put on an HTTP server accessible to Traffic Router.

	.. code-block:: json
		:caption: Example Configuration JSON

		{
			"customer": "YourCompany",
			"version": "1",
			"date" : "2017-05-23 03:28:25",
			"name": "Anonymous IP Blocking Policy",

			"anonymousIp": { "blockAnonymousVPN": true,
				"blockHostingProvider": true,
				"blockPublicProxy": true,
				"blockTorExitNode": true},

			"ip4Whitelist": ["192.168.30.0/24", "10.0.2.0/24", "10.1.1.1/32"],
			"ip6Whitelist": ["2001:550:90a::/48", "::1/128"],
			"redirectUrl": "http://youvebeenblocked.com"
		}

	anonymousIp
		Contains the types of IPs which can be checked against the Anonymous IP Database. There are 4 types of IPs which can be checked: :abbr:`VPN (Virtual Private Network)`\ s, Hosting Providers, Public Proxies, and :abbr:`TOR (The Onion Ring)` "Exit Nodes". Each type of IP can be enabled or disabled. If the value is true, IPs matching this type will be blocked when the feature is enabled in the :term:`Delivery Service`. If the value is false, IPs which match this type will not be blocked. If an IP matches more than 1 type and any type is enabled, the IP will be blocked.
	redirectUrl
		The URL that will be returned to the blocked clients. Without a :dfn:`redirectUrl`, the clients will receive an HTTP response code ``403 Forbidden``. With a :dfn:`redirectUrl`, the clients will be redirected with an HTTP response code ``302 Found``.
	ipWhiteList
		An optional element. It includes a list of :abbr:`CIDR (Classless Inter-Domain Routing)` blocks indicating the IPv4 and IPv6 subnets that are allowed by the rule. If this list exists and the value is not ``null``, client IPs will be matched against the :abbr:`CIDR (Classless Inter-Domain Routing)` list, and if there is any match, the request will be allowed. If there is no match in the white list, further anonymous blocking logic will continue.


#. Add the following three Anonymous Blocking :ref:`Parameters` in Traffic Portal with the "CRConfig.json" :ref:`parameter-config-file`, and ensure they are assigned to all of the Traffic Routers that should perform Anonymous Blocking:

	``anonymousip.policy.configuration``
		The URL of the Anonymous Blocking configuration file. Traffic Router will fetch the file from this URL.
	``anonymousip.polling.url``
		The URL of the Anonymous IP Database. Traffic Router will fetch the file from this URL.
	``anonymousip.polling.interval``
		The interval that Traffic Router polls the Anonymous Blocking configuration file and Anonymous IP Database.

	.. figure:: anonymous_blocking/01.png
		:width: 40%
		:align: center

#. Enable Anonymous Blocking for a :term:`Delivery Service` using the :ref:`Delivery Services view in Traffic Portal <tp-services-delivery-service>` (don't forget to save changes!)

	.. figure:: anonymous_blocking/02.png
		:width: 40%
		:align: center

#. Go to :ref:`the Traffic Portal CDNs view <tp-cdns>`, click on :guilabel:`Diff CDN Config Snapshot`, and click :guilabel:`Perform Snapshot`.

	.. figure:: anonymous_blocking/03.png
		:width: 40%
		:align: center


Traffic Router Access Log
=========================
Anonymous Blocking extends the field of ``rtype`` and adds a new field ``ANON_BLOCK`` in the Traffic Router ``access.log`` file to help monitor this feature. If the ``rtype`` in an access log is ``ANON_BLOCK`` then the client's IP was found in the Anonymous IP Database and was blocked.

.. seealso:: :ref:`tr-logs`
