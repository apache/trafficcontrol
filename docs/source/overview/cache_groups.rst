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

.. _cache-groups:

************
Cache Groups
************
A :dfn:`Cache Group` is - ostensibly - exactly what it sounds like it is: a group of :term:`cache servers`. More specifically, every server in a Traffic Control CDN must be in a Cache Group (even if they are not actually :term:`cache servers`). Typically a Cache Group is representative of the available :term:`cache servers` within a specific geographical location. Despite that :term:`cache servers` have their own :term:`Physical Locations`, when :term:`cache servers` are chosen to serve content to a client based on geographic location the geographic location actually used for comparisons is that for the Cache Group that contains it, not the geographic location of the :term:`cache server` itself.

The most typical :ref:`Types <cache-group-type>` of Cache Groups are EDGE_LOC_ which contain :term:`Edge-tier cache servers` and MID_LOC_ which contain :term:`Mid-tier cache servers`. The latter are each designated as a Parent_ of one or more of the former to fill out the two-tiered caching hierarchy of an :abbr:`ATC (Apache Traffic Control)` CDN.

Consider the example CDN in :numref:`fig-cg_hierarchy`. Here some country/province/region has been divided into quarters: Northeast, Southeast, Northwest, and Southwest. The arrows in the diagram indicate the flow of requests. If a client in the Northwest, for example, were to make a request to the :term:`Delivery Service`, it would first be directed to some :term:`cache server` in the "Northwest" Edge-tier :dfn:`Cache Group`. Should the requested content not be in cache, the Edge-tier server will select a parent from the "West" :dfn:`Cache Group` and pass the request up, caching the result for future use. All Mid-tier :dfn:`Cache Groups` (usually) answer to a single :term:`Origin` that provides canonical content. If requested content is not in the Mid-tier cache, then the request will be passed up to the :term:`Origin` and the result cached.

.. _fig-cg_hierarchy:

.. figure:: images/cg_hierarchy.*
	:align: center
	:width: 60%
	:alt: An illustration of Cache Group hierarchy

	An example CDN that shows the hierarchy between four Edge-tier :dfn:`Cache Groups`, two Mid-tier :dfn:`Cache Groups`, and one Origin


Regions, Divisions, and Locations
=================================
In addition to being in a Cache Group, all servers have to have a :term:`Physical Location`, which defines their geographic latitude and longitude. Each :term:`Physical Location` is part of a :term:`Region`, and each :term:`Region` is part of a :term:`Division`. For example, ``Denver`` could be the name of a :term:`Physical Location` in the ``Mile High`` :term:`Region` and that :term:`Region` could be part of the ``West`` :term:`Division`. The hierarchy between these terms is illustrated graphically in :ref:`topography-hierarchy`.

.. _topography-hierarchy:

.. figure:: images/topography.*
	:align: center
	:alt: A graphic illustrating the hierarchy exhibited by topological groupings
	:figwidth: 25%

	Topography Hierarchy

To create these structures in Traffic Portal, first make at least one :term:`Division` under :menuselection:`Topology --> Divisions`. Next enter the desired :term:`Region`\ (s) in :menuselection:`Topology --> Regions`, referencing the earlier-entered :term:`Division`\ (s). Finally, enter the desired :term:`Physical Location`\ (s) in :menuselection:`Topology --> Phys Locations`, referencing the earlier-entered :term:`Region`\ (s).

A Cache Group is a logical grouping of cache servers, that don't have to be in the same :term:`Physical Location` (in fact, usually a Cache Group is spread across minimally two :term:`Physical Locations` for redundancy purposes), but share geographical coordinates for content routing purposes. There is no strict requirement that :term:`cache servers` in a Cache Group share a :term:`Physical Location`, :term:`Region`, or :term:`Division`. This may be confusing at first as there are a few places in code, interfaces, or even documentation where Cache Groups are referred to as "Cache Locations" or even erroneously as "Physical Locations".

Properties
==========
Cache Groups are modeled several times over, in the Traffic Ops database, in Traffic Portal forms and tables, and several times for various :ref:`to-api` versions in the new Go Traffic Ops codebase. Go-specific data structures can be found at :atc-godoc:`lib/go-tc.CacheGroupNullable`. Rather than application-specific definitions, what follows is an attempt at consolidating all of the different properties and names of properties of Cache Group objects throughout the :abbr:`ATC (Apache Traffic Control)` suite. The names of these fields are typically chosen as the most human-readable and/or most commonly-used names for the fields, and when reading please note that in many cases these names will appear camelCased or snake_cased to be machine-readable. Any aliases of these fields that are not merely case transformations of the indicated, canonical names will be noted in a table of aliases.

.. seealso:: The API reference for Cache Group-related endpoints such as :ref:`to-api-cachegroups` contains definitions of the Cache Group object(s) returned and/or accepted by those endpoints.

.. _cache-group-asns:

ASNs
----
A Cache group can have zero or more :abbr:`ASNs (Autonomous System Numbers)` assigned to it, which is used to classify traffic that passes through a CDN. These are typically not represented on a Cache Group object itself, but rather as a separate object indicating the relationship, e.g. in the requests and responses of the :ref:`to-api-asns` endpoint.

.. seealso:: `The Autonomous System Wikipedia page <https://en.wikipedia.org/wiki/Autonomous_system_%28Internet%29>`_ for an explanation of what an :abbr:`ASN (Autonomous System Number)` actually is.

.. _cache-group-coordinate:

Coordinate
----------
.. tip:: Normally, one need not interact with this. In most contexts, this property of a Cache Group is not even exposed, but instead the Cache Group's Latitude_ and Longitude_ are exposed and should be directly manipulated.

The :dfn:`Coordinate` of a Cache Group defines the geographic coordinates of a Cache Group that is used for routing clients based on geographic location. It is also used to determine the "closest" Cache Group to another for the purposes of `Fallback to Closest`_.

Typically, this is expressed as an integral, unique identifier for the "Coordinate" object bound to the Cache Group that defines its geographic location, but occasionally it may appear as the name of that "Coordinate" object.

.. note:: When a new Cache Group is created, it is not necessary to first create a "Coordinate" object where it may reside. Instead, "Coordinates" are created automatically to reflect the Latitude_ and Longitude_ given to the newly created Cache Group. The name of the generated "Coordinate" will conform to the pattern :samp:`from_cachegroup_{Name}` where ``Name`` is the Cache Group's Name_. Because of this, creating new Cache Groups will fail if a "Coordinate" with a name matching that pattern already exists.

.. _cache-group-fallbacks:

Fallbacks
---------
:dfn:`Fallbacks` are a group of zero or more Cache Groups to be considered for routing when a Cache Group becomes unavailable due to high load or excessive maintenance. These are normally represented by an array of each Cache Group's ID_, but may occasionally appear as the Name_ or `Short Name`_ of each Cache Group.

This set is consulted before `Fallback to Closest`_ is taken into consideration.

.. seealso:: :ref:`health-proto`

.. table:: Aliases

	+-----------------------+-------------------------------------------+---------------------------------------------------------------------------------------------------------------+
	| Name                  | Use(s)                                    | Type(s)                                                                                                       |
	+=======================+===========================================+===============================================================================================================+
	| Failover Cache Groups | Traffic Portal forms - but **not** tables | List or array of :ref:`Names <cache-group-name>` as strings                                                   |
	+-----------------------+-------------------------------------------+---------------------------------------------------------------------------------------------------------------+
	| backupLocations       | In CDN :term:`Snapshots`                  | A sub-object called "list" which is a list or array of Cache Group :ref:`Names <cache-group-name>` as strings |
	+-----------------------+-------------------------------------------+---------------------------------------------------------------------------------------------------------------+
	| BackupCacheGroups     | Traffic Router source code                | A List of strings that are the :ref:`Names <cache-group-name>` of Cache Groups                                |
	+-----------------------+-------------------------------------------+---------------------------------------------------------------------------------------------------------------+

.. _cache-group-fallback-to-closest:

Fallback to Closest
-------------------
This is a boolean field which, when ``true`` (``True``, ``TRUE`` etc.) causes routing to "fall back" on the nearest Cache Group - geographically - when this Cache Group becomes unavailable due to high load and/or excessive maintenance.

When this *is* a "true" value, the closest Cache Group will be chosen if and only if any set of Fallbacks_ configured on the Cache Group has already been exhausted and no available Cache Groups were found..

.. seealso:: :ref:`health-proto`

.. table:: Aliases

	+--------------------------+----------------------+----------------------------------------+
	| Name                     | Use(s)               | Type(s)                                |
	+==========================+======================+========================================+
	| Fallback to Geo Failover | Traffic Portal forms | Unchanged (``bool``, ``Boolean`` etc.) |
	+--------------------------+----------------------+----------------------------------------+

.. _cache-group-id:

ID
--
All Cache Groups have an integral, unique identifier that is mainly used to reference it in the :ref:`to-api`.

Despite that a Cache Group's Name_ must be unique, this is the identifier most commonly used to represent a unique Cache Group in most contexts throughout :abbr:`ATC (Apache Traffic Control)`. One notable exception is in CDN :term:`Snapshots` and in routing configuration used by Traffic Router.

.. _cache-group-latitude:

Latitude
--------
The Cache Group's geomagnetic latitude for use in routing and for the purposes of `Fallback to Closest`_.

.. table:: Aliases

	+-----------------------+----------------------+----------------------------------------+
	| Name                  | Use(s)               | Type(s)                                |
	+=======================+======================+========================================+
	| Geo Magnetic Latitude | Traffic Portal forms | Unchanged (``number``, ``float`` etc.) |
	+-----------------------+----------------------+----------------------------------------+

.. _cache-group-localization-methods:

Localization Methods
--------------------
The :dfn:`Localization Methods` of a Cache Group define the methods by which Traffic Router is allowed to route clients to :term:`cache servers` within this Cache Group. This is a collection of the allowed methods, and the values in the collection are restricted to the following.

- "Coverage Zone File" (alias ``CZ`` in source code, database entries, and :ref:`to-api` requests/responses) allows Traffic Router to direct clients to this Cache Group if they were assigned a geographic location by looking up their IP address in the :term:`Coverage Zone File`.
- "Deep Coverage Zone File" (alias ``DEEP_CZ`` in source code, database entries, and :ref:`to-api` requests/responses) was intended to allow Traffic Router to direct clients to this Cache Group if they were assigned a geographic location by looking up their IP addresses in the :term:`Deep Coverage Zone File`. However, it **has no effect at all**. This option therefore will not appear in Traffic Portal forms.

	.. warning:: In order to make use of "deep caching" for a :term:`Delivery Service`, all that is required is that :term:`Delivery Service` has :ref:`ds-deep-caching` enabled. If that is done and a :term:`cache server` appears in the :term:`Deep Coverage Zone File` then clients can and will be routed using that method. There is no way to disable this behavior on a Cache Group (or otherwise) basis, and the precensce or absence of ``DEEP_CZ`` in a Cache Group's Localization Methods has no meaning.

- "Geo-IP Database" (alias ``GEO`` in source code, database entries, and :ref:`to-api` requests/responses) allows Traffic Router direct clients to this  Cache Group if the client's IP was looked up in a provided IP address-to-geographic location mapping database to provide their geographic location.

If none of these localization methods are in the set of allowed methods on a Cache Group, it is assumed that *clients should be allowed to be routed to that Cache Group regardless of the method used to determine their geographic location*.

This property only has meaning for Cache Groups containing :term:`Edge-tier cache servers`. Which is to say (one would hope) that it only has meaning for EDGE_LOC_ Cache Groups.

.. table:: Aliases

	+------------------------------+--------------------------------------------------+-----------------------------------------------------+
	| Name                         | Use(s)                                           | Type(s)                                             |
	+==============================+==================================================+=====================================================+
	| Enabled Localization Methods | Traffic Portal forms, Traffic Router source code | Unchanged (``Set<String>``, ``Array<string>`` etc.) |
	+------------------------------+--------------------------------------------------+-----------------------------------------------------+

.. _cache-group-longitude:

Longitude
---------
The Cache Group's geomagnetic longitude for use in routing and for the purposes of `Fallback to Closest`_.

.. table:: Aliases

	+------------------------+----------------------+----------------------------------------+
	| Name                   | Use(s)               | Type(s)                                |
	+========================+======================+========================================+
	| Geo Magnetic Longitude | Traffic Portal forms | Unchanged (``number``, ``float`` etc.) |
	+------------------------+----------------------+----------------------------------------+


.. _cache-group-name:

Name
----
A unique, human-friendly name for the Cache Group, with no special character restrictions or length limit.

Though this property must be unique for all Cache Groups, ID_ is more commonly used as a unique identifier for a Cache Group in most contexts.

.. table:: Aliases

	+------------+-----------------------+--------------------------------------+
	| Name       | Use(s)                | Type(s)                              |
	+============+=======================+======================================+
	| locationID | CDN :term:`Snapshots` | Unchanged (``str``, ``String`` etc.) |
	+------------+-----------------------+--------------------------------------+

.. _cache-group-parent:

Parent
------
A Cache Group can have a :dfn:`parent` Cache Group which has different meanings based on the Cache Group's Type_.

- An EDGE_LOC_ Cache Group's parent must be a MID_LOC_ Cache Group. When configuration files are generated for :term:`Edge-tier cache servers`, their :term:`parents` (different from a Cache Group parent) will be selected from the :term:`Mid-tier cache servers` contained within the Cache Group that is the parent of their containing Cache Group.
- A MID_LOC_ Cache Group's parent must be an ORG_LOC_ Cache Group. However, if any given MID_LOC_ *either* **doesn't** *have a parent, or* **does** *and it's an* ORG_LOC_, *then* **all** ORG_LOC_ *Cache Groups* - **even across CDNs** - *will be considered the parents of that* MID_LOC_ *Cache Group*. This parent relationship only has any meaning in the context of "multi-site-origins", as they are unnecessary in other scenarios.

	.. seealso:: :ref:`multi-site-origin-qht`

- For all other Cache Group :ref:`Types <cache-group-type>`, parent relationships have no meaningful semantics.

.. danger:: There is no safeguard in the data model or :ref:`to-api` that ensures these relationships hold. If they are violated, the resulting CDN behavior is undefined - and almost certainly undesirable.

:dfn:`Parents` are typically represented by their ID_, but may occasionally appear as their Name_.

.. seealso:: `The Apache Traffic Server documentation for the "parent.config" configuration file <https://docs.trafficserver.apache.org/en/7.1.x/admin-guide/files/parent.config.en.html>`_. :abbr:`ATC (Apache Traffic Control)` parentage relationships boil down to the ``parent`` field of that configuration file, which sets the :term:`parents` of :term:`cache servers`.

.. table:: Aliases

	+----------------------+---------------------------------------------------------------------------------+---------------------------------------------+
	| Name                 | Use(s)                                                                          | Type(s)                                     |
	+======================+=================================================================================+=============================================+
	| Parent Cache Group   | Traffic Portal forms and tables                                                 | Unchanged (``str``, ``String`` etc.)        |
	+----------------------+---------------------------------------------------------------------------------+---------------------------------------------+
	| ParentCacheGroupID   | Traffic Ops database, Traffic Ops Go code, :ref:`to-api` requests and responses | Positive integer (``int``, ``bigint`` etc.) |
	+----------------------+---------------------------------------------------------------------------------+---------------------------------------------+
	| ParentCacheGroupName | Traffic Ops database, Traffic Ops Go code, :ref:`to-api` requests and responses | Unchanged (``str``, ``String`` etc.)        |
	+----------------------+---------------------------------------------------------------------------------+---------------------------------------------+

.. _cache-group-parameters:

Parameters
----------
For unknown reasons, it's possible to assign :term:`Parameters` to Cache Groups. This has no attached semantics respected by :abbr:`ATC (Apache Traffic Server)`, and exists today only for compatibility purposes.

These are nearly always represented as a collection of :ref:`Parameter IDs <parameter-id>` but occasionally they can be expressed as full :term:`Parameter` objects.

.. table:: Aliases

	+----------------------+----------------------------------------------------------------+-------------------------------------------------------------------------------------+
	| Name                 | Use(s)                                                         | Type(s)                                                                             |
	+======================+================================================================+=====================================================================================+
	| cachegroupParameters | Certain :ref:`to-api` responses, sometimes in internal Go code | Usually a tuple of associative information like a Cache Group ID and a Parameter ID |
	+----------------------+----------------------------------------------------------------+-------------------------------------------------------------------------------------+

.. _cache-group-secondary-parent:

Secondary Parent
----------------
A :dfn:`secondary parent` of a Cache Group is used for fall-back purposes on a :term:`cache server`-to-:term:`cache server` basis after routing has already occurred (in contrast with Fallbacks_ and `Fallback to Closest`_ which operate at the routing step on a Cache Group-to-Cache Group basis).

For an explanation of what it means for one Cache Group to be the Parent_ of another, refer to the Parent_ section.

.. seealso:: `The Apache Traffic Server documentation for the "parent.config" configuration file <https://docs.trafficserver.apache.org/en/7.1.x/admin-guide/files/parent.config.en.html>`_. :abbr:`ATC (Apache Traffic Control)` secondary parentage relationships boil down to the ``secondary_parent`` field of that configuration file, which sets the "secondary parents" of :term:`cache servers`.

.. table:: Aliases

	+--------------------------------+---------------------------------------------------------------------------------+---------------------------------------------+
	| Name                           | Use(s)                                                                          | Type(s)                                     |
	+================================+=================================================================================+=============================================+
	| Secondary Parent Cache Group   | Traffic Portal forms and tables                                                 | Unchanged (``str``, ``String`` etc.)        |
	+--------------------------------+---------------------------------------------------------------------------------+---------------------------------------------+
	| SecondaryParentCacheGroupID    | Traffic Ops database, Traffic Ops Go code, :ref:`to-api` requests and responses | Positive integer (``int``, ``bigint`` etc.) |
	+--------------------------------+---------------------------------------------------------------------------------+---------------------------------------------+
	| SecondaryParentCacheGroupName  | Traffic Ops database, Traffic Ops Go code, :ref:`to-api` requests and responses | Unchanged (``str``, ``String`` etc.)        |
	+--------------------------------+---------------------------------------------------------------------------------+---------------------------------------------+

.. _cache-group-servers:

Servers
-------
The primary purpose of a Cache Group is to contain servers. In most scenarios it is implied or assumed that the servers are :term:`cache servers`, but this is not required and in fact it is not by any means uncommon for the contained servers to be of some other arbitrary :term:`Type`.

A Cache Group can have zero or more assigned servers, and each server can belong to at most one Cache Group.

.. _cache-group-short-name:

Short Name
----------
This is typically an abbreviation of the Cache Group's Name_. The main difference is that it isn't required to be unique.

.. _cache-group-type:

Type
----
A Cache Group's :term:`Type` determines what kind of servers it contains. Note that there's no real restriction on the kinds of servers that a Cache Group can contain, and this :term:`Type` serves as more of a guide in certain contexts. The :term:`Types` available by default are described in this section.

.. tip:: Because :term:`Types` are mutable, the actual :term:`Types` that can describe Cache Groups cannot be completely or precisely defined. However, there are no good reasons of which this author can think to modify or delete the default :term:`Types` herein described, and honestly the good reasons to even merely add to their ranks are likely few.

.. table:: Aliases

	+-----------+-----------------------------------------------------------------------------+----------------------------------------------+
	| Name      | Use(s)                                                                      | Type(s)                                      |
	+===========+=============================================================================+==============================================+
	| Type ID   | Traffic Ops client and server Go code, :ref:`to-api` requests and responses | positive integer (``int``, ``bigint``, etc.) |
	+-----------+-----------------------------------------------------------------------------+----------------------------------------------+
	| Type Name | :ref:`to-api` requests and responses, Traffic Ops database                  | unchanged (``string``, ``String`` etc.)      |
	+-----------+-----------------------------------------------------------------------------+----------------------------------------------+


EDGE_LOC
""""""""
This :term:`Type` of Cache Group contains :term:`Edge-tier cache servers`.

MID_LOC
"""""""
This :term:`Type` of Cache Group contains :term:`Mid-tier cache servers`

ORG_LOC
"""""""
This :term:`Type` of Cache Group contains :term:`Origins`. The primary purpose of these is to group :term:`Origins` for the purposes of "multi-site-origins", and it's suggested that if that doesn't meet your use-case that these be mostly avoided. In general, it's not strictly necessary to create :term:`Origin` *servers* in ATC at all, unless you have to support "multi-site-origins".

.. seealso:: :ref:`multi-site-origin-qht`

TC_LOC
""""""
A catch-all :term:`Type` of Cache Group that's meant to house infrastructure servers that gain no special semantics based on the Cache Group containing them, e.g. Traffic Portal instances.

TR_LOC
""""""
This :term:`Type` of Cache Group is meant specifically to contain Traffic Router instances.
