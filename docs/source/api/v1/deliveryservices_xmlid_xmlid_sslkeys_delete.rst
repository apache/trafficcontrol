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

.. _to-api-v1-deliveryservices-xmlid-xmlid-sslkeys-delete:

***************************************************
``deliveryservices/xmlId/{{xmlid}}/sslkeys/delete``
***************************************************
.. deprecated:: ATCv4
	Use the ``DELETE`` method of :ref:`to-api-deliveryservices-xmlid-xmlid-sslkeys` instead.

``GET``
=======
:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object (string)

Request Structure
-----------------
.. table:: Request Path Parameters

	+-------+----------+-------------------------------------------------------------+
	| Name  | Required | Description                                                 |
	+=======+==========+=============================================================+
	| xmlId | yes      | The :ref:`ds-xmlid` of the desired :term:`Delivery Service` |
	+-------+----------+-------------------------------------------------------------+

.. table:: Request Query Parameters

	+---------+----------+------------------------------------------------------------+
	|   Name  | Required |          Description                                       |
	+=========+==========+============================================================+
	| version | no       | The version number of the SSL keys that shall be retrieved |
	+---------+----------+------------------------------------------------------------+

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Wed, 18 Mar 2020 17:36:10 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: Pj+zCoOXg19nGNxcSkjib2iDjG062Y3RcEEV+OYnwbGIsLcpa0BKZleY/qJOKT5DkSoX2qQkckUxUqdDxjVorQ==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 18 Mar 2020 16:36:10 GMT
	Content-Length: 173

	{
		"alerts": [
			{
				"text": "This endpoint is deprecated, please use DELETE /deliveryservices/xmlId/:xmlid/sslkeys instead",
				"level": "warning"
			}
		],
		"response": "Successfully deleted ssl keys for demo1"
	}