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

#. Create two target Delivery Services in Traffic Ops. They must both be HTTP Delivery Services that are part of the same CDN.

	.. figure:: steering/01.png
		:width: 80%
		:align: center
		:alt: Table of Target Delivery Services

		Target Delivery Services

#. Create a Delivery Service with Type STEERING or CLIENT_STEERING in Traffic Ops.

	.. figure:: steering/02.png
		:width: 50%
		:align: center
		:alt: Delivery Service Creation Page for STEERING Delivery Service

		Creating a STEERING Delivery Service

#. In the 'More' drop-down menu, click 'View Targets' and then use the blue '+' to assign targets.

	.. figure:: steering/03.png
		:width: 50%
		:align: center
		:alt: Table of STEERING Targets

		STEERING Targets


#. If desired, a 'steering' user can create filters for the target Delivery Services (only available via API). Sample JSON request body:

	.. code-block:: json

		{
			"filters": [
			 {
				 "pattern": ".*\\gototarget1\\..*",
				 "deliveryService": "target-deliveryservice-1"
			 }
			],
			"targets": [
			 {
				 "weight": "1000",
				 "deliveryService": "target-deliveryservice-1"
			 },
			 {
				 "weight": "9000",
				 "deliveryService": "target-deliveryservice-2"
			 }
			 {
				 "order": -1,
				 "deliveryService": "target-deliveryservice-3"
			 }
			 {
				 "order": 3,
				 "deliveryService": "target-deliveryservice-4"
			 }
			]
		}

	Sample script of ``curl`` commands to accomplish this, given the above request body is saved as ``/tmp/steering.json`` [1]_:

	.. code-block:: shell

		curl -sc cookie.jar https://to.cdn.local/api/1.2/user/login -d '{"u":"admin","p":"twelve"}'
		curl -sb cookie.jar -XPUT "https://to.cdn.local/internal/api/1.2/steering/steering-ds" -d @/tmp/steering.json

#. Any requests to Traffic Router for the steering Delivery Service should now be routed to target Delivery Services based on configured weight or order.

.. [1] This example also assumes that the Traffic Ops instance is running at ``to.cdn.local`` and the administrative username and password are ``admin`` and ``twelve``, respectively. This is *not* recommended in production, but merely meant to replicate the default 'CDN-in-a-box' environment!
