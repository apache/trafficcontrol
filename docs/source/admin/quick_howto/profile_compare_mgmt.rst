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

.. _profile_compare_mgmt:

***********************
Profile Parameters Diff
***********************
In Traffic Portal *all* users can diff the :term:`Parameters` of any 2 :term:`Profiles` side-by-side, and users with a higher level :term:`Role` ("operations" or "admin") can easily add or remove parameters from each profile as necessary.

The ability to diff 2 profiles can be found under :menuselection:`Configure --> Profiles --> More --> Diff Profile Parameters`

.. figure:: profile_compare_mgmt/compare_profiles_menu.png
	:align: center
	:alt: A screenshot of the "Diff Profile Parameters" menu item

	The "Diff Profile Parameters" menu item

Once you have selected the :guilabel:`Diff Profile Parameters` menu item, you will be asked to choose 2 profiles to diff.

.. figure:: profile_compare_mgmt/select_profiles_dialog.png
	:width: 60%
	:align: center
	:alt: A screenshot of the "Diff Profile Parameters" dialog

	The "Diff Profile Parameters" dialog

All parameters exclusively assigned to one profile but not the other will be displayed with their profile membership displayed side-by-side. In addition, by selecting the :guilabel:`Show All Params` link, the user can see a superset of all parameters across the 2 profiles. Both views provide users with higher level permissions ("operations" or "admin") the ability to easily remove or add parameters for each profile and persist the final state of both profiles (or restore the original state and discard changes). As the user makes changes, a blue shadow is added to all modified checkboxes.

.. figure:: profile_compare_mgmt/compare_profiles_table.png
	:align: center
	:alt: A screenshot of the "Diff Profile Parameters" table

	The "Diff Profile Parameters" table
