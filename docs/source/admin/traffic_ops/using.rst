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

.. _to-using:

*******************
Traffic Ops - Using
*******************

.. deprecated:: 3.0
	The Traffic Ops UI is deprecated, and will be removed entirely in the next major release (4.0). A much better way to interact with the CDN is to :ref:`use Traffic Portal <usingtrafficportal>`, which is the the only UI that will be receiving updates for the foreseeable future.

The Traffic Ops Menu
====================
.. figure:: images/12m.png
	:align: center
	:alt: The Traffic Ops Landing Page

	The Traffic Ops Landing Page

The following tabs are available in the menu at the top of the Traffic Ops user interface.

.. index::
	Change Log

ChangeLog
---------
The Changelog table displays the changes that are being made to the Traffic Ops database through the Traffic Ops user interface. This tab will show the number of changes since you last visited this tab in (brackets) since the last time you visited this tab. There are currently no sub menus for this tab.


Help
----
Help for Traffic Ops and Traffic Control. Hover over this tab to get the following options:

+---------------+---------------------------------------------------------------------+
|     Option    |                             Description                             |
+===============+=====================================================================+
| About         | Traffic Ops information, such as version, database information, etc |
+---------------+---------------------------------------------------------------------+
| Release Notes | Release notes for the most recent releases of Traffic Ops           |
+---------------+---------------------------------------------------------------------+
| Logout        | Logout from Traffic Ops                                             |
+---------------+---------------------------------------------------------------------+

.. index::
	Invalidate Content
	Purge

.. _purge:

Invalidate Content
==================
Invalidating content on the CDN is sometimes necessary when the origin was mis-configured and something is cached in the CDN  that needs to be removed. Given the size of a typical Traffic Control CDN and the amount of content that can be cached in it, removing the content from all the caches may take a long time. To speed up content invalidation, Traffic Ops will not try to remove the content from the caches, but it makes the content inaccessible using the *regex_revalidate* ATS plugin. This forces a *revalidation* of the content, rather than a new get.

.. Note:: This method forces a HTTP *revalidation* of the content, and not a new *GET* - the origin needs to support revalidation according to the HTTP/1.1 specification, and send a ``200 OK`` or ``304 Not Modified`` as applicable.

To invalidate content:

#. Click **Tools > Invalidate Content**
#. Fill out the form fields:

	- Select the *:term:`Delivery Service`**
	- Enter the **Path Regex** - this should be a `PCRE <http://www.pcre.org/>`_ compatible regular expression for the path to match for forcing the revalidation. Be careful to only match on the content you need to remove - revalidation is an expensive operation for many origins, and a simple ``/.*`` can cause an overload condition of the origin.
	- Enter the **Time To Live** - this is how long the revalidation rule will be active for. It usually makes sense to make this the same as the ``Cache-Control`` header from the origin which sets the object time to live in cache (by ``max-age`` or ``Expires``). Entering a longer TTL here will make the caches do unnecessary work.
	- Enter the **Start Time** - this is the start time when the revalidation rule will be made active. It is pre-populated with the current time, leave as is to schedule ASAP.

#. Click the **Submit** button.
