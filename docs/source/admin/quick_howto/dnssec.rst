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

.. _dnssec-qht:

****************
Configure DNSSEC
****************
In order to support `DNSSEC <https://en.wikipedia.org/wiki/Domain_Name_System_Security_Extensions>`_ in Traffic Router, Traffic Portal provides some actions for managing DNSSEC keys for a CDN and associated :term:`Delivery Services`. DNSSEC keys consist of a :abbr:`KSK (Key Signing Key)` which is used to sign other DNSKEY records as well as a :abbr:`ZSK (Zone Signing Keys)` which is used to sign other records. DNSSEC keys are stored in Traffic Vault and should only be accessible to Traffic Ops. Other applications needing access to this data, such as Traffic Router, must use the :ref:`to-api` to retrieve this information.

.. seealso:: :ref:`tr-dnssec`

.. note:: :abbr:`DNSSEC (DNS Security Extensions)` is only supported for DNS-:ref:`Routed <ds-types>` :term:`Delivery Services`.

Creating Keys
=============
.. Note:: In order for Traffic Ops to successfully store keys in Traffic Vault, at least one Traffic Vault server needs to be configured in Traffic Ops.


#. Go to :guilabel:`CDNs` and click on the desired CDN.

	.. figure:: dnssec/00.png
		:width: 60%
		:align: center
		:alt: Screenshot of the Traffic Portal UI depicting the CDNs page

		CDNs Page

#. Click on :menuselection:`More --> Manage DNSSEC Keys`.

	.. figure:: dnssec/01.png
		:width: 60%
		:align: center
		:alt: Screenshot of the Traffic Portal UI depicting the CDN details page

		CDN Details Page

#. Click on the :guilabel:`Generate DNSSEC Keys` button.

	.. figure:: dnssec/02.png
		:width: 60%
		:align: center
		:alt: Screenshot of the Traffic Portal UI depicting the CDN DNSSEC Key Management page

		DNSSEC Key Management Page

#. A modal will pop up asking you to confirm that you want to proceed.

	.. figure:: dnssec/03.png
		:width: 30%
		:align: center
		:alt: Screenshot of the Traffic Portal UI depicting the CDN DNSSEC Key Generation confirmation modal

		Confirmation Modal

#. Input the required information (reasonable defaults should be generated for you). When done, click on the green :guilabel:`Generate` button.

	.. note:: Depending upon the number of :term:`Delivery Services` in the CDN, generating DNSSEC keys may take several seconds.

	.. figure:: dnssec/04.png
		:width: 50%
		:align: center
		:alt: Screenshot of the Traffic Portal UI depicting the CDN DNSSEC Key Generation page

		DNSSEC Key Generation Page

#. You will be prompted to confirm the changes by typing the name of the CDN into a text box. After doing so, click on the red :guilabel:`Confirm` button.

	.. figure:: dnssec/05.png
		:width: 30%
		:align: center
		:alt: Screenshot of the Traffic Portal UI depicting the confirmation modal for committing changes to DNSSEC Keys.

		DNSSEC Key Change Confirmation

#. In order for :abbr:`DNSSEC (DNS Security Extensions)` to work properly, the :abbr:`DS (Delegation of Signing)` Record information needs to be added to the parent zone of the CDN's domain (e.g. If the CDN's domain is 'ciab.cdn.local' the parent zone is 'cdn.local'). If you control your parent zone you can enter this information yourself, otherwise you will need to work with your DNS team to get the :abbr:`DS (Delegation of Signing)` Record added to the parent zone.

Enabling and Disabling DNSSEC on a CDN
======================================
#. Once :abbr:`DS (Delegation of Signing)` Record information has been added to the parent zone, DNSSEC needs to be activated for the CDN so that Traffic Router will sign responses. Go to the CDN details page for this CDN, and set the 'DNSSEC Enabled' field to 'true' (or 'false' to disable DNSSEC), then click the green :guilabel:`Update` button.

	.. figure:: dnssec/06.png
		:width: 60%
		:align: center
		:alt: Screenshot of the Traffic Portal UI depicting the details page for a CDN when changing its 'DNSSEC Enabled' field

		Change 'DNSSEC Enabled' to 'true'

#. :abbr:`DNSSEC (DNS Security Extensions)` should now be active (or inactive, if disabled) on your CDN and Traffic Router should be signing responses. This should be tested e.g. with this :manpage:`dig(1)` command: ``dig edge.cdn.local. +dnssec``.

Regenerating KSKs
=================
When :abbr:`KSK (Key-Signing Key)` expiration is approaching (default 365 days), it is necessary to manually generate a new :abbr:`KSK (Key Signing Key)` for the :abbr:`TLD (Top Level Domain)` and add the :abbr:`DS (Delegation of Signing)` Record to the parent zone. In order to avoid signing errors, it is suggested that an effective date is chosen which allows time for the :abbr:`DS (Delegation of Signing)` Record to be added to the parent zone before the new :abbr:`KSK (Key-Signing Key)` becomes active.

Regenerating Keys
=================
It's also possible to re-do the generation of DNSSEC keys once they have already been generated. This can be done at any point after `Creating Keys`_ by clicking on the then-enabled :guilabel:`Regenerate Keys` button on the CDN's details page.
