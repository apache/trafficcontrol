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

.. _to-api-v4-oc-fci-configuration-request-id-approved:

***************************************************
``OC/CI/configuration/request/{{id}}/{{approved}}``
***************************************************

``PUT``
=======
Triggers an asynchronous task to update the configuration for the :abbr:`uCDN (Upstream Content Delivery Network)` and the specified host by adding the request to a queue to be reviewed later. This returns a 202 Accepted status and an endpoint to be used for status updates.

:Auth. Required: Yes
:Roles Required: "admin"
:Permissions Required: CDNI-ADMIN:READ, CDNI-ADMIN:UPDATE
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+-----------+----------------------------------------------------------------------------------------+
	| Name      |                 Description                                                            |
	+===========+========================================================================================+
	|  id       | The integral identifier for the configuration update request to be approved or denied. |
	+-----------+----------------------------------------------------------------------------------------+
	|  approved | A boolean for whether to approve a configuration change request or not.                |
	+-----------+----------------------------------------------------------------------------------------+

Response Structure
------------------

.. code-block:: http
	:caption: Response Example For Approved Change

	HTTP/1.1 200 OK
	Content-Type: application/json

	{ "response": "Successfully updated configuration." }

.. code-block:: http
	:caption: Response Example For Denied Change

	HTTP/1.1 200 OK
	Content-Type: application/json

	{ "response": "Successfully denied configuration update request." }
