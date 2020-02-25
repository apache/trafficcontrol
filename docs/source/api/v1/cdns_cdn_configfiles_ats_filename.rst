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

.. _to-api-v1-cdns-cdn-configfiles-ats-filename:

*********************************************
``cdns/{{cdn}}/configfiles/ats/{{filename}}``
*********************************************
.. seealso:: The :ref:`to-api-v1-servers-server-configfiles-ats` endpoint

``GET``
=======
Gets the configuration file ``filename`` from the CDN ``cdn`` (by either name or ID).

:Auth. Required: Yes
:Roles Required: "operations"
:Response Type:  **NOT PRESENT** - endpoint returns custom text/plain response (represents the contents of the requested configuration file)

Request Structure
-----------------
.. table:: Request Path Parameters

	+-----------+-------------------+----------+----------------------------------------------------------+
	| Parameter | Type              | Required | Description                                              |
	+===========+===================+==========+==========================================================+
	| cdn       | string or integer | yes      | Either the name or integral, unique, identifier of a CDN |
	+-----------+-------------------+----------+----------------------------------------------------------+
	| filename  | string            | yes      | The name of a configuration file used by ``cdn``         |
	+-----------+-------------------+----------+----------------------------------------------------------+

Response Structure
------------------
.. note:: If the file identified by ``filename`` does exist, but is not used by the entire CDN, a JSON response will be returned and the ``alerts`` array will contain a ``"level": "error"`` node which identifies the correct scope of the configuration file.

.. code-block:: squid
	:caption: Response Example

	# DO NOT EDIT - Generated for CDN-in-a-Box by Traffic Ops (https://trafficops.infra.ciab.test:443/) on Thu Oct 25 13:26:31 UTC 2018
	cond %{REMAP_PSEUDO_HOOK}
	set-conn-dscp 8 [L]


