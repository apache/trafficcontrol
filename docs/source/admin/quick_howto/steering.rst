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

.. _steering-qht:

***********************************
Configure Delivery Service Steering
***********************************

#. Create two target :term:`Delivery Services` in Traffic Portal. They must both be HTTP :term:`Delivery Services` that are part of the same CDN.

	.. figure:: steering/01.png
		:width: 80%
		:align: center
		:alt: Table of Target Delivery Services

		Target :term:`Delivery Services`

#. Create a :term:`Delivery Service` with Type ``STEERING`` or ``CLIENT_STEERING`` in Traffic Ops.

	.. figure:: steering/02.png
		:width: 50%
		:align: center
		:alt: Delivery Service Creation Page for STEERING Delivery Service

		Creating a STEERING :term:`Delivery Service`

#. Click :menuselection:`More --> View Targets` and then use the blue :guilabel:`+` button to assign targets.

	.. figure:: steering/03.png
		:width: 50%
		:align: center
		:alt: Table of STEERING Targets

		STEERING Targets


#. If desired, a 'steering' :term:`Role` user can create filters for the target :term:`Delivery Services` using :ref:`to-api-steering-id-targets`

	.. note:: This is only available via the :ref:`to-api`; no functionality for manipulating steering targets is offered by Traffic Portal. This feature has been requested and is tracked by `GitHub Issue #2811 <https://github.com/apache/trafficcontrol/issues/2811>`_

#. Any requests to Traffic Router for the steering :term:`Delivery Service` should now be routed to target :term:`Delivery Services` based on configured weight or order.

.. note:: This example assumes that the Traffic Ops instance is running at ``to.cdn.local`` and the administrative username and password are ``admin`` and ``twelve``, respectively. This is *not* recommended in production, but merely meant to replicate the default :ref:`ciab` environment!
