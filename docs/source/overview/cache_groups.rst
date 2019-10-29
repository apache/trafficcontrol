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

***********************************
Regions, Locations and Cache Groups
***********************************
All servers have to have a :term:`Physical Location`, which defines their geographic latitude and longitude. Each :term:`Physical Location` is part of a :term:`Region`, and each :term:`Region` is part of a :term:`Division`. For example, ``Denver`` could be the name of a :term:`Physical Location` in the ``Mile High`` :term:`Region` and that :term:`Region` could be part of the ``West`` :term:`Division`. The hierarchy between these terms is illustrated graphically in :ref:`topography-hierarchy`.

.. _topography-hierarchy:

.. figure:: images/topography.*
	:align: center
	:alt: A graphic illustrating the hierarchy exhibited by topological groupings
	:figwidth: 25%

	Topography Hierarchy

To create these structures in Traffic Portal, first make at least one :term:`Division` under :menuselection:`Topology --> Divisions`. Next enter the desired :term:`Region`\ (s) in :menuselection:`Topology --> Regions`, referencing the earlier-entered :term:`Division`\ (s). Finally, enter the desired :term:`Physical Location`\ (s) in :menuselection:`Topology --> Phys Locations`, referencing the earlier-entered :term:`Region`\ (s).

All servers also have to be part of a :term:`Cache Group`. A :term:`Cache Group` is a logical grouping of cache servers, that don't have to be in the same :term:`Physical Location` (in fact, usually a :term:`Cache Group` is spread across minimally two :term:`Physical Locations` for redundancy purposes), but share geographical coordinates for content routing purposes.
