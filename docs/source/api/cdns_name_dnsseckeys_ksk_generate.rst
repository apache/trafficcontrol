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

.. _to-api-cdns-name-dnsseckeys-ksk-generate:

*****************************************
``cdns/{{name}}/dnsseckeys/ksk/generate``
*****************************************

.. versionadded:: 1.4

``POST``
========
Generates a new Key-Signing Key (KSK) for a specific CDN.

:Auth. Required: Yes
:Roles Required: "admin"
:Response Type:  Object (string)

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+----------+---------------------------------------------------------+
	| Name | Required | Description                                             |
	+======+==========+=========================================================+
	| name | yes      | The name of the CDN for which the KSK will be generated |
	+------+----------+---------------------------------------------------------+

.. table:: Request Data Parameters

	+--------------------+----------+---------+--------------------------------------------------------------------------------------------------------------------------------------------------------+
	| Name               | Required | Type    | Description                                                                                                                                            |
	+====================+==========+=========+========================================================================================================================================================+
	| ``expirationDays`` | yes      | integer | The number of days until the newly generated KSK expires                                                                                               |
	+--------------------+----------+---------+--------------------------------------------------------------------------------------------------------------------------------------------------------+
	| ``effectiveDate``  | no       | string  | The time at which the newly generated KSK becomes effective, in `RFC3339 <https://tools.ietf.org/html/rfc3339>`_ format - defaults to the current time |
	+--------------------+----------+---------+--------------------------------------------------------------------------------------------------------------------------------------------------------+

Response Structure
------------------
.. code-block:: json
	:caption: Response Example

	{ "response": "Successfully generated ksk dnssec keys for my-cdn-name" }
