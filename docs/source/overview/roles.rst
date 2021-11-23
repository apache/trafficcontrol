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

.. _roles:

*****
Roles
*****
A :dfn:`Role` is a collection of Permissions that may be associated with zero or more users.

.. seealso:: The :atc-godoc:`lib/go-tc.RoleV40` Go structure documentation.

The model for a Role in the most recent version of the :ref:`to-api` is given in :ref:`role-interface` as a Typescript interface.

.. _role-interface:

.. code-block:: typescript
	:caption: Role as a Typescript Interface

	interface Role {
		description: string;
		lastUpdated?: Date; //RFC3339 timestamp string - only in responses
		name: string;
		permissions: Array<string>;
	}

Description
===========
The description of a Role is some human-readable text that describes what the
Role is for.

.. versionchanged:: ATCv6.1.0
	In versions of ATC earlier than 6.1.0, this property was allowed to be ``null`` or blank - this is no longer the case.

Last Updated
============
The date and time at which the Role was last modified.

Name
====
A unique name for the Role. The only Role Name that is treated specially is `The Admin Role`_.

.. versionchanged:: ATCv6.1.0
	Though it was undocumented, there used to be many Role names that were treated specially by ATC in versions earlier than 6.1.0, including (but not limited to), "disallowed", "operations", and "steering".

The Admin Role
--------------
A Role with a Name that is exactly "admin" is special in that it is always treated as having all Permissions_, regardless of what Permissions_ are actually assigned to it. For this reason, the "admin" Role may never be deleted or modified.

Permissions
===========
A Role's :dfn:`Permissions` is a set of the Permissions afforded to the Role that define the :ref:`to-api` interactions it is allowed to have. Usually - but not always - this maps directly to some HTTP method of some :ref:`to-api` endpoint.

.. table:: Aliases

	+--------------+-------------------------------------------------------+----------------------------------+
	| Name         | Use(s)                                                | Type(s)                          |
	+==============+=======================================================+==================================+
	| Capabilities | Legacy :ref:`to-api` versions' requests and responses | unchanged (array/set of strings) |
	+--------------+-------------------------------------------------------+----------------------------------+
