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
.. _cachegroup-fallback-qht:

*******************************
Configure Cache Group Fallbacks
*******************************

.. seealso:: :ref:`tp-configure-cache-groups`

#. Go to 'Topology', click on :term:`Cache Groups`, and click on your desired :term:`Cache Group` or click the :guilabel:`+` button to create a new :term:`Cache Group`.

	.. figure:: cachegroup_fallback/00.png
		:width: 60%
		:align: center
		:alt: Screenshot of the Traffic Portal UI depicting the Cache Groups page

		Cache Groups Page

#. Verify that the :term:`Cache Group` is of type EDGE_LOC. :term:`Cache Group` Failovers only apply to EDGE_LOC :term:`Cache Groups`.

	.. figure:: cachegroup_fallback/01.png
		:width: 60%
		:align: center
		:alt: Screenshot of the Traffic Portal UI depicting the Cache Group details page

		Cache Group Details Page

#. Once EDGE_LOC is selected, the Failover Cache Groups section will appear at the bottom of the page. If you are editing an existing :term:`Cache Group`, then the current Failovers will be listed. If creating a new :term:`Cache Group`, the Fallback to Geo Failover box will default to be checked.

	.. figure:: cachegroup_fallback/02.png
		:width: 60%
		:align: center
		:alt: Screenshot of the Traffic Portal UI depicting the Failover Cache Groups section of the Cache Group details page

		Failover Cache Groups Section of Cache Group Details Page

#. To add a new Failover to the list, select the "Add Failover :term:`Cache Group`" drop down and choose which :term:`Cache Group` you would like. While in the drop down, you can also type in order to search.

	.. figure:: cachegroup_fallback/03.png
		:width: 60%
		:align: center
		:alt: Screenshot of the Traffic Portal UI depicting the Failover Cache Groups section of the Cache Group details page

		Add New Failover Cache Group Section of Cache Group Details Page

#. The order of the Failovers is important. If you want to reorder the Failovers, you can drag and drop them into a new position. A red line will appear to show where the Failover will be dropped.

	.. figure:: cachegroup_fallback/04.png
		:width: 60%
		:align: center
		:alt: Screenshot of the Traffic Portal UI depicting the Failover Cache Groups Drag and Drop of the Cache Group details page

		Failover Cache Groups Section Drag and Drop Functionality

#. To remove a Failover, click the trash can symbol on the right hand side of the list.

	.. figure:: cachegroup_fallback/05.png
		:width: 60%
		:align: center
		:alt: Screenshot of the Traffic Portal UI depicting the Failover Cache Groups Delete of the Cache Group details page

		Failover Cache Groups Delete

#. Click the :guilabel:`Update` button (if editing existing :term:`Cache Group`) or the :guilabel:`Create` button (if creating new :term:`Cache Group`) in order to save the Failovers to the :term:`Cache Group`.
