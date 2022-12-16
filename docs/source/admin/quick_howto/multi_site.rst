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

.. _multi-site-origin-qht:

***************************
Configure Multi-Site Origin
***************************

The following steps will take you through the procedure of setting up an :abbr:`MSO (Multi-Site Origin)`.
#. Create :term:`Cache Groups` for the origin locations, and assign the appropriate parent-child relationship between the Mid-tier :term:`Cache Group`\ (s) and origin :term:`Cache Groups`. Each Mid-tier :term:`Cache Group` can be assigned a primary and secondary origin parent :term:`Cache Group`. When the :term:`Cache Group` parent configuration is generated, origins in the primary :term:`Cache Groups` will be listed first, followed by origins in the secondary :term:`Cache Group`. Origin servers assigned to the :term:`Delivery Service` that are assigned to neither the primary nor secondary :term:`Cache Groups` will be listed last.

	.. figure:: multi_site/00.png
		:scale: 100%
		:align: center

#. Create a profile to assign to each of the origins:

	.. figure:: multi_site/01.png
		:scale: 100%
		:align: center

#. Create server entries for the origination vips:

	.. figure:: multi_site/02.png
		:scale: 100%
		:align: center

#. Check the multi-site check box in the :term:`Delivery Service` screen:

	.. figure:: multi_site/mso-enable.png
		:scale: 100%
		:align: center

#. Assign the org servers to the :term:`Delivery Service` that will have the multi site feature. Origin servers assigned to a :term:`Delivery Service` with multi-site checked will be assigned to be the origin servers for this :term:`Delivery Service`.

	.. figure:: multi_site/03.png
		:scale: 100%
		:align: center

	.. Note:: “Origin Server Base URL” uniqueness: In order to enable Mid-tier :term:`Cache Group` to distinguish :term:`Delivery Services` by different :abbr:`MSO (Multi-Site Origin)` algorithms while performing parent fail-over, it requires that :abbr:`OSBU (Origin Server Base URL)` for each :abbr:`MSO (Multi-Site Origin)`-enabled :term:`Delivery Service` is unique. This means that the :abbr:`OSBU (Origin Server Base URL)` of an :abbr:`MSO (Multi-Site Origin)`-enabled :term:`Delivery Service` should be different from the :abbr:`OSBU (Origin Server Base URL)`\ s of any other :term:`Delivery Service`, regardless of whether they are :abbr:`MSO (Multi-Site Origin)`-enabled or not. The exceptions to this rule are:

		- If there are multiple CDNs created on the same Traffic Ops, :term:`Delivery Services` across different CDNs may have the same :abbr:`OSBU (Origin Server Base URL)` configured.
		- If several :term:`Delivery Services` in the same CDN have the same :abbr:`MSO (Multi-Site Origin)` algorithm configured, they may share the same :abbr:`OSBU (Origin Server Base URL)`.
		- If delivery services are assigned with different Mid-tier :term:`Cache Groups` respectively, they can share the same :abbr:`OSBU (Origin Server Base URL)`.
		- This :abbr:`OSBU (Origin Server Base URL)` must be valid - :abbr:`ATS (Apache Traffic Server)` will perform a DNS lookup on this :abbr:`FQDN (Fully Qualified Domain Name)` even if IPs, not DNS, are used in the :file:`parent.config`.
		- The :abbr:`OSBU (Origin Server Base URL)` entered as the "Origin Server Base URL" will be sent to the origins as a host header. All origins must be configured to respond to this host.

#. Create a delivery service profile. This must be done to set the :abbr:`MSO (Multi-Site Origin)` algorithm. Also, as of :abbr:`ATS (Apache Traffic Server)` 6.x, multi-site options must be set as parameters within the :file:`parent.config`. Header rewrite parameters will be ignored. See `ATS parent.config <https://docs.trafficserver.apache.org/en/6.2.x/admin-guide/files/parent.config.en.html>`_ for more details. These :term:`Parameters` are now handled by the creation of a :term:`Delivery Service` :term:`Profile`.

	a) Create a :term:`Profile` of the :ref:`profile-type` ``DS_PROFILE`` for the :term:`Delivery Service` in question.

		.. figure:: multi_site/ds_profile.png
			:scale: 50%
			:align: center

	#) Click :guilabel:`Show profile parameters` to bring up the :term:`Parameters` screen for the :term:`Profile`. Create the following :term:`Parameters`:

		+-----------------------------------------+------------------+--------------------------+-------------------------+
		| Parameter Name                          | Config File Name | Value                    | ATS parent.config value |
		+=========================================+==================+==========================+=========================+
		| last.algorithm                          | parent.config    | true, false, strict,     | round_robin             |
		|                                         |                  | consistent_hash          |                         |
		+-----------------------------------------+------------------+--------------------------+-------------------------+
		| last.parent_retry                       | parent.config    | simple_retry, both,      | parent_retry            |
		|                                         |                  | unavailable_server_retry | (deprecated)            |
		+-----------------------------------------+------------------+--------------------------+-------------------------+
		| last.unavailable_server_retry_responses | parent.config    | list of server response  | defaults to the value   |
		|                                         |                  | codes, eg "500,502,503"  | in records.config       |
		|                                         |                  |                          | when unused.            |
		+-----------------------------------------+------------------+--------------------------+-------------------------+
		| last.max_unavailable_server_retries     | parent.config    | Number of retries made   | defaults to the value   |
		|                                         |                  | after specified errors   | in records.config       |
		|                                         |                  |                          | when unused.            |
		+-----------------------------------------+------------------+--------------------------+-------------------------+
		| last.simple_server_retry_responses      | parent.config    | list of server response  | defaults to the value   |
		|                                         |                  | codes, eg "404"          | in records.config       |
		|                                         |                  |                          | when unused.            |
		+-----------------------------------------+------------------+--------------------------+-------------------------+
		| last.max_simple_retries                 | parent.config    | Nubmer of retries made   | defaults to the value   |
		|                                         |                  | after specified errors   | in records.config       |
		|                                         |                  |                          | when unused.            |
		+-----------------------------------------+------------------+--------------------------+-------------------------+

		.. figure:: multi_site/ds_profile_parameters.png
			:scale: 100%
			:align: center

    .. deprecated:: ATC 6.2

	#) In the :term:`Delivery Service` page, select the newly created ``DS_PROFILE`` and save the :term:`Delivery Service`.

#. Turn on parent_proxy_routing in the MID :term:`Profile`.

.. Note:: Support for multisite configurations with single-layer CDNs is now available. If a :term:`Cache Groups` defined parents are either blank or of the type ``ORG_LOC``, that :term:`cache server`'s ``parent.config`` will be generated as a top layer cache, even if it is an edge. In the past, ``parent.config`` generation was strictly determined by cache type. The new method examines the parent :term:`Cache Group` definitions and generates the :file:`parent.config` accordingly.
