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

.. _cdni_admin:

****************************
CDNi Administration
****************************

:abbr:`CDNi (Content Delivery Network Interconnect)` specifications define the standards for interoperability within the :abbr:`CDN (Content Delivery Network)` and open caching ecosystems set forth by the :abbr:`IETF (Internet Engineering Task Force)`. This integration utilizes the :abbr:`APIs (Application Programming Interfaces)` defined by the :abbr:`SVA (Streaming Video Alliance)`.

.. seealso:: For complete details on CDNi and the related API specifications see :rfc:`8006`, :rfc:`8007`, :rfc:`8008`, and the :abbr:`SVA (Streaming Video Alliance)` documents titled `Footprint and Capabilities Interface: Open Caching API`, `Open Caching API Implementation Guidelines`, `Configuration Interface: Part 1 Specification - Overview & Architecture`, `Configuration Interface: Part 2 Specification – CDNi Metadata Model Extensions`, and `Configuration Interface: Part 3 Specification – Publishing Layer APIs`.

In short, these documents describe the :abbr:`CDNi (Content Delivery Network Interconnect)` metadata interface that enables interconnected :abbr:`CDNs (Content Delivery Networks)` to exchange content distribution metadata to enable content acquisition and delivery. These define the interfaces through which a :abbr:`uCDN (Upstream Content Delivery Network)` and a :abbr:`dCDN (Downstream Content Delivery Network)` can communicate configuration and capacity information.

For our use case, it is assumed that :abbr:`ATC (Apache Traffic Control)` is the :abbr:`dCDN (Downstream Content Delivery Network)`.

	..  Note:: This is currently under construction and will be for a while. This document will be updated as new features are supported.

/OC/FCI/advertisement
=====================
.. seealso:: :ref:`to-api-oc-fci-advertisement`

The advertisement response is unique for the :abbr:`uCDN (Upstream Content Delivery Network)` and contains the complete footprint and capabilities information structure the :abbr:`dCDN (Downstream Content Delivery Network)` wants to expose. This endpoint will return an array of generic :abbr:`FCI (Footprint and Capabilities Advertisement Interface)` base objects, including type, value, and footprint for each. Currently supported base object types are `FCI.Capacitiy` and `FCI.Telemetry` but these will be expanded in the future.

/OC/CI/configuration
====================
.. seealso:: :ref:`to-api-oc-fci-configuration`

An endpoint that is used to push (``PUT``), fetch (``GET``), or delete (``DELETE``) the entire metadata set for a given :abbr:`uCDN (Upstream Content Delivery Network)` from a :abbr:`JWT (JSON Web Token)`. This puts the requested change into a queue to be reviewed later and returns an endpoint to view the asynchronous status updates.

.. Note:: This is under construction. Currently only ``PUT`` is supported and in a very limited sense.

/OC/CI/configuration/{{host}}
=============================
.. seealso:: :ref:`to-api-oc-fci-configuration-host`

An endpoint that is used to push (``PUT``), fetch (``GET``), or delete (``DELETE``) the metadata set that is attached to host name for a given :abbr:`uCDN (Upstream Content Delivery Network)` from a :abbr:`JWT (JSON Web Token)`. This puts the requested change into a queue to be reviewed later and returns an endpoint to view the asynchronous status updates.

.. Note:: This is under construction. Currently only ``PUT`` is supported and in a very limited sense.

/OC/CI/configuration/request/{{id}}/{{approved}}
================================================
.. seealso:: :ref:`to-api-oc-fci-configuration-request-id-approved`

This endpoint allows a user to approve or deny a queued update request from the previous endpoints. A denial will result in the removal from the queue and a ``FAILED`` status update. An approval will result in the changes being made to the configuration and a ``SUCCEEDED`` status update.

.. Note:: This is under construction and only supports very limited metadata field and limited configuration updates.
