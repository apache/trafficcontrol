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

.. _content-invalidation:

****************************
Forcing Content Invalidation
****************************
Invalidating content on the CDN is sometimes necessary when the :term:`Origin` was mis-configured and something is cached in the CDN  that needs to be removed.

.. impl-detail:: Given the size of a typical Traffic Control CDN and the amount of content that can be cached in it, removing the content from all the caches may take a long time. To speed up content invalidation, Traffic Control does not try to remove the content from the caches, but it makes the content inaccessible using the `regex_revalidate plugin for Apache Traffic Server <https://docs.trafficserver.apache.org/en/8.0.x/admin-guide/plugins/regex_revalidate.en.html>`_. This forces a "re-validation" of the content.

To invalidate content for a specific :term:`Delivery Service`, follow these steps:

#. Select the desired :term:`Delivery Service` from the :ref:`tp-services-delivery-service` view of Traffic Portal

	.. figure:: content_invalidation/01.png
		:align: center
		:alt: The Traffic Portal Delivery Services view

		The Traffic Portal Delivery Services view

#. From the :guilabel:`More` drop-down menu, select :menuselection:`Manage Invalidation Requests`

	.. figure:: content_invalidation/02.png
		:align: center
		:alt: The 'Manage Invalidation Requests' option under 'More'

		Select 'Manage Invalidation Requests'

#. From the :guilabel:`More` drop-down menu on this page, select :menuselection:`Create Invalidation Request`

	.. figure:: content_invalidation/03.png
		:align: center
		:alt: The 'Create Invalidation Request' option under 'More'

		Select 'Create Invalidation Request'

#. Fill out this form. The "Path Regex" field should be a `PCRE <http://www.pcre.org/>`_-compatible regular expression that matches all content that must be invalidated - and should **not** match any content that must *not* be invalidated. "TTL (hours)" specifies the number of hours for which the invalidation should remain active. Best practice is to set this to the same as the content's cache lifetime (typically set in the :term:`Origin`'s ``Cache-Control`` response header). :ref:`job-invalidation-type` describes how content will be invalidated.

	.. figure:: content_invalidation/04.png
		:align: center
		:alt: The new content invalidation submission form

		The 'new content invalidation submission' Form

#. Click on the :guilabel:`Create` button to finalize the content invalidation.
