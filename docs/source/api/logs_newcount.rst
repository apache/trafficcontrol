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


.. _to-api-logs-newcount:

*****************
``logs/newcount``
*****************

``GET``
=======
Gets the number of new changes made to the Traffic Control system - "new" being defined as the last time the client requested either :ref:`to-api-logs` or :ref:`to-api-logs-days-days`.

.. note:: This endpoint's functionality is implemented by the :ref:`to-api-logs` and :ref:`to-api-logs-days-days` endpoints' responses setting cookies for the client to use when requesting _this_ endpoint. Take care that your client respects cookies!

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Object

Request Structure
-----------------
No parameters available

Response Structure
------------------
:newLogcount: The integer number of new changes

.. code-block:: json
	:caption: Response Example

	{ "response": {
		"newLogcount": 4
	}}
