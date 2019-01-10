..
..
.. _cachegroup-fallback-qht:

******************************
Configure CacheGroup Fallbacks
******************************

.. seealso:: :ref:`tp-configure-cache-groups`

#. Go to 'Topology', click on Cache Groups, and click on your desired Cache Group or click the + button to create a new Cache Group.

	.. figure:: cachegroup_fallback_qht/00.png
		:width: 60%
		:align: center
		:alt: Screenshot of the Traffic Portal UI depicting the Cache Groups page

		Cache Groups Page

#. Verify that the Cache Group is of type EDGE_LOC.  Cache Group Failovers only apply to EDGE_LOC Cache Groups.

	.. figure:: cachegroup_fallback_qht/01.png
		:width: 60%
		:align: center
		:alt: Screenshot of the Traffic Portal UI depicting the Cache Group details page

		Cache Group Details Page

#. Once EDGE_LOC is selected, the Failover Locations section will appear at the bottom of the page.  If you are editing an existing Cache Group, then the current Failovers will be listed.  If creating a new Cache Group, the Fallback to Geo Failover box will default to be checked.

	.. figure:: cachegroup_fallback_qht/02.png
		:width: 60%
		:align: center
		:alt: Screenshot of the Traffic Portal UI depicting the Failover Locations section of the Cache Group details page

		Failover Locations Section of Cache Group Details Page

#. To add a new Failover to the list, select the "Add Failover Cache Group" drop down and choose which Cache Group you would like.  While in the drop down, you can also type in order to search.

	.. figure:: cachegroup_fallback_qht/03.png
		:width: 60%
		:align: center
		:alt: Screenshot of the Traffic Portal UI depicting the Failover Locations section of the Cache Group details page

		Failover Locations Section of Cache Group Details Page

#. The order of the Failovers is important.  If you want to reorder the Failovers, you can drag and drop them into a new position.  A red line will appear to show where the Failover will be dropped.

	.. figure:: cachegroup_fallback_qht/04.png
		:width: 60%
		:align: center
		:alt: Screenshot of the Traffic Portal UI depicting the Failover Locations Drag and Drop of the Cache Group details page

		Failover Locations Section Drag and Drop Functionality

#. To remove a Failover, click the trash can symbol on the right hand side of the list.

	.. figure:: cachegroup_fallback_qht/05.png
		:width: 60%
		:align: center
		:alt: Screenshot of the Traffic Portal UI depicting the Failover Locations Delete of the Cache Group details page

		Failover Locations Delete

#. Click the Update button (if editing existing Cache Group) or the Create button (if creating new Cache Group) in order to save the Failovers to the Cache Group.