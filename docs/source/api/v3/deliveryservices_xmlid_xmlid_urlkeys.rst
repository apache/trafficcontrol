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

.. _to-api-v3-deliveryservices-xmlid-xmlid-urlkeys:

********************************************
``deliveryservices/xmlId/{{xmlid}}/urlkeys``
********************************************

``GET``
=======
.. seealso:: :ref:`to-api-v3-deliveryservices-id-urlkeys`

Retrieves URL signing keys for a :term:`Delivery Service`.

.. caution:: This method will return the :term:`Delivery Service`'s **PRIVATE** URL signing keys! Be wary of using this endpoint and **NEVER** share the output with anyone who would be unable to see it on their own.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+-------+------------------------------------------------------+
	|  Name |              Description                             |
	+=======+======================================================+
	| xmlid | The 'xml_id' of the desired :term:`Delivery Service` |
	+-------+------------------------------------------------------+

Response Structure
------------------
:key<N>: The private URL signing key for this :term:`Delivery Service` as a base-64-encoded string, where ``<N>`` is the "generation" of the key e.g. the first key will always be named ``"key0"``. Up to 16 concurrent generations are retained at any time (``<N>`` is always on the interval [0,15])

.. code-block:: json
	:caption: Response Example

	{ "response": {
		"key9":"ZvVQNYpPVQWQV8tjQnUl6osm4y7xK4zD",
		"key6":"JhGdpw5X9o8TqHfgezCm0bqb9SQPASWL",
		"key8":"ySXdp1T8IeDEE1OCMftzZb9EIw_20wwq",
		"key0":"D4AYzJ1AE2nYisA9MxMtY03TPDCHji9C",
		"key3":"W90YHlGc_kYlYw5_I0LrkpV9JOzSIneI",
		"key12":"ZbtMb3mrKqfS8hnx9_xWBIP_OPWlUpzc",
		"key2":"0qgEoDO7sUsugIQemZbwmMt0tNCwB1sf",
		"key4":"aFJ2Gb7atmxVB8uv7T9S6OaDml3ycpGf",
		"key1":"wnWNR1mCz1O4C7EFPtcqHd0xUMQyNFhA",
		"key11":"k6HMzlBH1x6htKkypRFfWQhAndQqe50e",
		"key10":"zYONfdD7fGYKj4kLvIj4U0918csuZO0d",
		"key15":"3360cGaIip_layZMc_0hI2teJbazxTQh",
		"key5":"SIwv3GOhWN7EE9wSwPFj18qE4M07sFxN",
		"key13":"SqQKBR6LqEOzp8AewZUCVtBcW_8YFc1g",
		"key14":"DtXsu8nsw04YhT0kNoKBhu2G3P9WRpQJ",
		"key7":"cmKoIIxXGAxUMdCsWvnGLoIMGmNiuT5I"
	}}
