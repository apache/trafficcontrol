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

The following steps will take you through the procedure of setting up a Multi-Site Origin (MSO).

#. Create Cache Groups for the origin locations, and assign the appropriate parent-child relationship between the mid and origin Cache Groups. Each mid Cache Group can be assigned a primary and secondary origin parent Cache Group. When the mid cache parent configuration is generated, origins in the primary Cache Groups will be listed first, followed by origins in the secondary Cache Group. Origin servers assigned to the Delivery Service that are assigned to neither the primary nor secondary Cache Groups will be listed last.

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

#. Check the multi-site check box in the delivery service screen:

	.. figure:: multi_site/mso-enable.png
		:scale: 100%
		:align: center

#. Assign the org servers to the delivery service that will have the multi site feature. Org servers assigned to a delivery service with multi-site checked will be assigned to be the origin servers for this DS.

	.. figure:: multi_site/03.png
		:scale: 100%
		:align: center

	.. Note:: “Origin Server Base URL” uniqueness: In order to enable MID caches to distinguish delivery services by different MSO algorithms while performing parent fail-over, it requires that “Origin Server Base URL” (OSBU) for each MSO-enabled Delivery Service is unique. This means that the OSBU of a MSO-enabled Delivery Service should be different with the OSBUs of any other Delivery Service, regardless of whether they are MSO-enabled or not. The exceptions to this rule are:

		- If there are multiple CDNs created on the same Traffic Ops, delivery services across different CDNs may have the same OSBU configured.
		- If several delivery services in the same CDN have the same MSO algorithm configured, they may share the same OSBU.
		- If delivery services are assigned with different MID cache groups respectively, they can share the same OSBU.
		- This OSBU must be valid - ATS will perform a DNS lookup on this FQDN even if IPs, not DNS, are used in the parent.config.
		- The OSBU entered as the "Origin Server Base URL" will be sent to the origins as a host header. All origins must be configured to respond to this host.

#. Create a delivery service profile. This must be done to set the MSO algorithm. Also, as of ATS 6.x, multi-site options must be set as parameters within the parent.config. Header rewrite parameters will be ignored. See `ATS parent.config <https://docs.trafficserver.apache.org/en/6.2.x/admin-guide/files/parent.config.en.html>`_ for more details. These parameters are now handled by the creation of a delivery service profile.

	a) Create a profile of the type DS_PROFILE for the delivery service in question.

		.. figure:: multi_site/ds_profile.png
			:scale: 50%
			:align: center

	#) Click "Show profile parameters" to bring up the parameters screen for the profile. Create parameters for the following:

		+----------------------------------------+------------------+--------------------------+-------------------------+
		| Parameter Name                         | Config File Name | Value                    | ATS parent.config value |
		+========================================+==================+==========================+=========================+
		| mso.algorithm                          | parent.config    | true, false, strict,     | round_robin             |
		|                                        |                  | consistent_hash          |                         |
		+----------------------------------------+------------------+--------------------------+-------------------------+
		| mso.parent_retry                       | parent.config    | simple_retry, both,      | parent_retry            |
		|                                        |                  | unavailable_server_retry |                         |
		+----------------------------------------+------------------+--------------------------+-------------------------+
		| mso.unavailable_server_retry_responses | parent.config    | list of server response  | defaults to the value   |
		|                                        |                  | codes, eg "500,502,503"  | in records.config       |
		|                                        |                  |                          | when unused.            |
		+----------------------------------------+------------------+--------------------------+-------------------------+
		| mso.max_simple_retries                 | parent.config    | Nubmer of retries made   | defaults to the value   |
		|                                        |                  | after a 4xx error        | in records.config       |
		|                                        |                  |                          | when unused.            |
		+----------------------------------------+------------------+--------------------------+-------------------------+
		| mso.max_unavailable_server_retries     | parent.config    | Nubmer of retries made   | defaults to the value   |
		|                                        |                  | after a 5xx error        | in records.config       |
		|                                        |                  |                          | when unused.            |
		+----------------------------------------+------------------+--------------------------+-------------------------+


		.. figure:: multi_site/ds_profile_parameters.png
			:scale: 100%
			:align: center

	#) In the delivery service page, select the newly created DS_PROFILE and save the delivery service.

#. Turn on parent_proxy_routing in the MID profile.
