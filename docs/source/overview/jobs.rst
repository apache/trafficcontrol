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

.. _jobs:

*************************
Content Invalidation Jobs
*************************
:dfn:`Content Invalidation Jobs`, or simply "jobs" as they are sometimes known, are ways of forcing :term:`cache servers` to treat content as no longer valid, bypassing their normal caching policies.

In general, this *should* be unnecessary, because a well-behaved :term:`Origin` *should* be setting its HTTP caching headers properly, so that content is only considered valid for some appropriate time intervals. Occasionally, however, an :term:`Origin` will be too optimistic with its caching instructions, and when content needs to be updated, :term:`cache servers` need to be informed that they must check back with the :term:`Origin`. Content Invalidation Jobs allow this to be done for specific patterns of assets, so that :term:`cache servers` will check back in with the :term:`Origin` and verify that the content they have cached is still valid.

The model for Content Invalidation Job as API objects is given in :ref:`jobs-model`.

.. _jobs-model:

.. code-block:: typescript
	:caption: Content Invalidation Job as a Typescript interface.

	/** This is the form used to create a new Content Invalidation Job */
	interface ContentInvalidationJobCreationRequest {
		deliveryService: string;
		invalidationType: "REFRESH" | "REFETCH";
		regex: `/${string}` | `\\/${string}`; // must also be a valid RegExp
		startTime: Date; // RFC3339 string
		ttlHours: number;
	}

	/**
	 * This is the form used to return representations of Content Invalidation
	 * Requests to clients.
	 */
	interface ContentInvalidationJob {
		assetUrl: string;
		createdBy: string;
		deliveryService: string;
		id: number;
		invalidationType: "REFRESH" | "REFETCH";
		startTime: Date; // RFC3339 string
		ttlHours: number;
	}

.. _job-asset-url:

Asset URL
---------
This property only appears in responses from the :ref:`to-api` (and in the bodies of ``PUT`` requests to :ref:`to-api-jobs`, where the scheme and host/authority sections of the URL is held immutable). The :dfn:`Asset URL` is constructed from the `Regular Expression`_ used in the creation of a Content Invalidation Job and the :ref:`ds-origin-url` of the :term:`Delivery Service` for which it was created. It is a URL that has a valid regular expression as its path (and may not be "percent-encoded" where a normal URL typically would be). Requests from CDN clients for content that matches this pattern will trigger Content Invalidation behavior.

.. _job-created-by:

Created By
----------
The username of the user who created the Content Invalidation Job is stored as the :dfn:`Created By` property of the Content Invalidation Job.

.. _job-ds:

Delivery Service
----------------
A Content Invalidation Job can only act on content for a single :term:`Delivery Service` - invalidating content for multiple :term:`Delivery Services` requires multiple Content Invalidation Jobs. The :dfn:`Delivery Service` property of a Content Invalidation Job holds the :ref:`ds-xmlid` of the :term:`Delivery Service` on which it operates.

.. versionchanged:: 4.0
	In earlier API versions, this property was allowed to be either the integral, unique identifier of the target :term:`Delivery Service`, *or* its :ref:`ds-xmlid` - this is no longer the case, but it should always be safe to use the :ref:`ds-xmlid` in any case.

.. _job-id:

ID
--
The integral, unique identifier for the Content Invalidation Job, assigned to it upon its creation.

.. _job-invalidation-type:

Invalidation Type
-----------------
:dfn:`Invalidation Type` defines how a :term:`cache server` should go about ensuring that its cache is valid.

The normal operating mode for a Content Invalidation Job is to force the :term:`cache server` to send a request to the :term:`Origin` to verify that its cache is valid. If that is the case, no extra work is done and business as usual resumes. However, some :term:`Origins` are misconfigured and do not respond as required by HTTP specification. In this case, it is strongly advised to fix the :term:`Origin` so that it properly implements HTTP. However, if an :term:`Origin` is sending cache-able responses to requests, and cannot be trusted to verify the validity of cached content based on cache-controlling HTTP headers (e.g. :mailheader:`If-Modified-Since`) instead returning responses like ``304 Not Modified`` *even when the content has in fact been modified*, **and** if correcting this behavior is not an option, then the :term:`cache server` may be forced to pretend that the content it has was actually invalidated by the :term:`Origin` and must be completely re-fetched.

The two values allowed for a Content Invalidation Job's Invalidation Type are:

REFRESH
	A :dfn:`REFRESH` Content Invalidation Job instructs :term:`cache servers` to behave normally - when matching content is requested, send an upstream request to (eventually) the :term:`Origin` with cache-controlling HTTP headers, and trust the :term:`Origin`'s response. The vast majority of all Content Invalidation Jobs should most likely use this Invalidation Type.
REFETCH
	Rather than treating the cached content as "stale", the :term:`cache servers` processing a :dfn:`REFETCH` Content Invalidation Job should fetch the cached content again, regardless of what the :term:`Origin` has to say about the validity of their caches. These types of Content Invalidation Jobs cannot be created without a proper "semi-global" :ref:`refetch_enabled Parameter <parameter-name-refetch_enabled>`.

.. caution:: A "REFETCH" Content Invalidation Job should be used **only** when the :term:`Origin` is not properly configured to support HTTP caching, and will return invalid or incorrect responses to conditional requests  as described in section 4.3.2 of :rfc:`7234`. In any other case, this will cause undo load on both the :term:`Origin` and the requesting :term:`cache servers`, and "REFRESH" should be used instead.

.. _job-regex:

Regular Expression
------------------
The :dfn:`Regular Expression` of a Content Invalidation Job defines the content on which it acts. It is used to match URL *paths* (including the query string - but **not** including document fragments, which are not sent in HTTP requests) of content to be invalidated, and is combined with the :ref:`ds-origin-url` of the :term:`Delivery Service` for which the Content Invalidation Job was created to obtain a final pattern that is made available as the `Asset URL`_.

.. note:: While the :ref:`to-api` and :ref:`tp-overview` both require the Regular Expression to begin with ``/`` (so that it matches URL paths), the :ref:`to-api` allows optionally escaping this leading character with a "backslash" :kbd:`\\`, while :ref:`tp-overview` does not. As ``/`` is not syntactically important to regular expressions, the use of a leading :kbd:`\\` should be avoided where possible, and is only allowed for legacy compatibility reasons.

.. table:: Aliases/Synonyms

	+------------+--------------------------------------------------------------------------------+-------------------------------+
	| Name       | Use(s)                                                                         | Type                          |
	+============+================================================================================+===============================+
	| Path Regex | In Traffic Portal forms                                                        | unchanged (String, str, etc.) |
	+------------+--------------------------------------------------------------------------------+-------------------------------+
	| regex      | In raw :ref:`to-api` requests and responses, internally in multiple components | unchanged (String, str, etc.) |
	+------------+--------------------------------------------------------------------------------+-------------------------------+

.. _job-start-time:

Start Time
----------
Content Invalidation Jobs are planned in advance, by setting their :dfn:`Start Time` to some point in the future (the :ref:`to-api` will refuse to create Content Invalidation Jobs with a Start Time in the past). Content Invalidation Jobs will have no effect until their Start Time.

.. _job-ttl:

TTL
---
The :dfn:`TTL` of a Content Invalidation Job defines how long a Content Invalidation Job should remain in effect. This is generally expressed as an integer number of hours.

.. table:: Aliases/Synonyms

	+------------+-----------------------------------------+----------------------------------------------------------------------+
	| Name       | Use(s)                                  | Type                                                                 |
	+============+=========================================+======================================================================+
	| parameters | In legacy :ref:`to-api` versions        | A string, containing the TTL in the format :samp:`TTL:{Actual TTL}h` |
	+------------+-----------------------------------------+----------------------------------------------------------------------+
	| ttlHours   | In :ref:`to-api` requests and responses | Unchanged (unsigned integer number of hours)                         |
	+------------+-----------------------------------------+----------------------------------------------------------------------+
