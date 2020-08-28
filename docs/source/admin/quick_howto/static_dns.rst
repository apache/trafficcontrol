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

.. _static-dns-qht:

******************************
Configuring Static DNS Entries
******************************
Static DNS records (historically "entries") can be configured within the DNS subdomain of a given :term:`Delivery Service`. In a typical scenario, the :term:`Delivery Service` will have DNS records automatically generated based on its "xml_id" and "routing name", and the name and subdomain of the CDN to which it belongs. For example, in the :ref:`ciab` default environment, the "demo1" :term:`Delivery Service` has an automatically generated DNS record for ``video.demo1.mycdn.ciab.test``. Configuring a static DNS record allows for further extension of this, for example, one could create an ``A`` record that enforces lookups of the name ``foo.demo1.mycdn.ciab.test`` resolve to the IPv4 address ``192.0.2.1``.

.. note:: It's only possible to create static DNS records within a :term:`Delivery Service`'s subdomain. That is, one could not create an ``A`` record for ``foo.bar.mycdn.ciab.test`` on the :ref:`ciab` :term:`Delivery Service` "demo1", since "demo1"'s subdomain is ``demo1.mycdn.ciab.test``.

.. seealso:: This guide covers how to set up static DNS records using Traffic Portal. It's also possible to do so directly using the :ref:`to-api` endpoint :ref:`to-api-staticdnsentries`.

Example
=======
To set up the aforementioned rule, follow these steps.

#. In Traffic Portal, expand the :ref:`tp-services` sidebar menu and select :guilabel:`Delivery Services`.
#. From the now-displayed table of :term:`Delivery Services`, select the desired one for static DNS record configuration.
#. From the :guilabel:`More` drop-down menu, select :guilabel:`Static DNS Entries`. The displayed table will probably be empty.

	.. figure:: static_dns/00.png
		:alt: The static DNS entries table page
		:align: center

		The Static DNS Entries Table Page

#. Click on the blue :guilabel:`+` button to add a new static DNS Entry
#. Fill in all of the fields.

	Host
		This is the lowest-level DNS label that will be used in the DNS record. In the :ref:`ciab` scenario, for example, entering ``foo`` here will result in a full DNS name of ``foo.demo1.mycdn.ciab.test``.
	Type
		Indicates the type of DNS record that will be created. The available types are

			* A
			* AAAA
			* CNAME
			* TXT

	TTL
		The :abbr:`TTL (Time To Live)` of the DNS record, after which clients will be expected to re-request name resolution.
	Address
		The meaning of this field depends on the value of the "Type" field.

			* If the "Type" is ``A``, this must be a valid IPv4 address
			* If the "Type" is ``AAAA``, this must be a valid IPv6 address
			* If the "Type" is ``CNAME``, this must be a valid DNS name - **not** an IP address at all
			* If the "Type" is ``TXT``, no restrictions are placed on the content whatsoever

	.. figure:: static_dns/01.png
		:alt: An example static DNS entry form
		:align: center

		An Example Static DNS Entry Form

#. Click on the green :guilabel:`Create` button to finalize the changes.
#. At this point, although the static DNS record has been created, it will have no effect until a new CDN :term:`Snapshot` is taken. Once that is done (and enough time has passed for Traffic Router to poll for the changes), the new DNS record should be usable through the CDN's designated Traffic Router.

	.. code-block:: console
		:caption: Example DNS Query to Test a New Static DNS Entry within :ref:`ciab`

		$ docker exec cdninabox_enroller_1 dig +noall +answer foo.demo1.mycdn.ciab.test
		foo.demo1.mycdn.ciab.test. 42	IN	A	192.0.2.1
