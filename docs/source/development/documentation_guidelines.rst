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

.. _docs-guide:

************************
Documentation Guidelines
************************
The Apache Traffic Control documentation is written in :abbr:`RST (reStructuredText)` and uses the `Sphinx documentation build system <http://www.sphinx-doc.org/en/master/>`_ to convert these into the desired output document format. This collection of guidelines does **not** aim to be a primer in :abbr:`RST (reStructuredText)`, but merely a style guide regarding how the components of the document ought to be formatted and structured. It may also point out some features of the markup language of which a writer may not be aware.

.. seealso:: `The docutils RST reference <http://docutils.sourceforge.net/rst.html>`_.

.. _docs-build:

Building
========
This documentation uses the `Sphinx documentation build system <http://www.sphinx-doc.org/en/master/>`_, and as such requires a Python3 version that is at least 3.4.1. It also has dependency on Sphinx, and Sphinx extensions and themes. All of these can be easily installed using `pip` by referencing the requirements file like so:

.. code-block:: shell
	:caption: Run from the Repository's Root Directory

	python3 -m pip install --user -r docs/source/requirements.txt

Once all dependencies have been satisfied, build using the Makefile at ``docs/Makefile``.

Alternatively, it is also possible to :ref:`pkg` or to :ref:`build-with-dc`, both of which will output a documentation "tarball" to ``dist/``.

Writing
=======
When writing documentation, the most important things to remember are:

- Spell Check. Most text editors have this built-in (e.g. :kbd:`F6` in Sublime Text) or have plugins that will do this for you.
- Proof-Read. Spell-checkers won't catch grammatical errors or poor wording, so it's very important to actually proof-read all documentation *before* submitting it in a Pull Request.
- Make Sure the Documentation Actually Builds. Please actually verify the documentation not only builds, but builds *correctly*. That means there probably shouldn't be any warnings, no malformed tables etc. and it also means that new documentation is actually accessible once built. It's not enough to create a new :file:`{something}.rst` file, that file must actually be linked to from some other, already included document. Some warnings may be considered acceptable, but do be prepared to defend them.
- Traffic Ops UI is Dead. Do not ever create documentation that references or includes images of the Traffic Ops UI. That is officially dead now, and if the documentation being created is best made with references to a user-friendly UI, such references, examples and/or images should all be to/of Traffic Portal.

Formatting
----------
Whenever possible, avoid specifying manual line breaks, except as required by :abbr:`RST (reStructuredText)` syntax. Extremely long lines will be wrapped by the user-agent, compiler, or output format as necessary. A single blank line may be used to separate paragraphs. This means that the 'flow break' character should never need to be used, i.e. no line in the documentation should ever match the regular expression :regexp:`^\\|$`.

Abbreviations
"""""""""""""
When using an abbreviation, acronym or initialism for the first time on a page, it **must** be named fully and followed by the abbreviation in parentheses e.g. "Fully Qualified Domain Name (FQDN)". Strictly speaking, the *best* way to create an abbreviation is to always fully name it in parentheses immediately following the abbreviation using the ``:abbr:`` :abbr:`RST (reStructuredText)` text role e.g. ``:abbr:`FQDN (Fully Qualified Domain Name)```, but it's not reasonable to expect that of everyone. Some abbreviations can be assumed to be understood by the documentation's target audience, and do not need full naming; they are general, basic networking and computing terms including (though not strictly limited to):

- API
- CSS
- DNS
- HTML
- HTTP
- HTTPS
- IP/IPv4/IPv6
- ISO
- JPG
- JSON
- PDF
- PNG
- RPM
- SQL
- SSL
- SVG
- TCP
- TLS
- UDP
- URL
- URI
- XML
- YAML

Please do **not** abbreviate Traffic Control terms e.g. :term:`Cache Group`, :term:`Delivery Service`. See `Terms`_ for the proper way to use these terms.

Floating Objects
""""""""""""""""
"Floating objects" are images, tables, source code listings, and equations. These may not be placed relative to other content exactly as shown in the source :abbr:`RST (reStructuredText)` document, as it may be necessary to move them for e.g. page breaks in PDF documents so they are not split across pages.

Figures
'''''''
Images should *always* be included inside of a ``.. figure`` directive. **Always** caption figures to make their purpose clear, as well as to make them directly link-able inside of the document and include them in figure listings. Though not syntactically required, figures should, in general be sized explicitly. The size should not be *absolute*, however; i.e. use ``70%`` not ``540px``. Figures should, in general, be centered on the page. When drawings, graphs, or diagrams are included they should ideally be provided in both SVG and PNG formats, and included using globbing as ``filename.*``. This will use the appropriate format for the output type.

Source Code Listings
''''''''''''''''''''
Do not ever use the double-colon (``::`` ) directive to mark a section of text as a source code listing. This not only doesn't support direct linking or provide a caption, but also uses the default "domain" - which is Python - for syntax highlighting. Instead, use ``.. code-block`` to include source code in the documentation. Source code must always be left-aligned, so do not provide any configuration options that would alter the default.

.. code-block:: rst
	:caption: Example Usage of the code-block Directive

	.. code-block:: syntax
		:caption: A short, meaningful caption
		:linenos:

		``:linenos:`` is an optional field which will include line numbers in the listing. 'syntax'
		should be the name of a valid Pygments syntax.

Tables
''''''
Tables should be included in ``.. table`` directive bodies, **never** as a floating, block-quoted tabular environment. This ensures that all tables will be captioned, which makes their purpose clear and makes them directly link-able in the output as well as includes them in table listings. Tables should avoid wrapping lines until they reach 215 characters in width in the source :abbr:`reStructuredText` document (including indention which should be counted as 4 characters per TAB). No table may ever exceed 215 characters in width. Tables should, in general be left-aligned (which is the default configuration). For the usage or command-line flags or options of a utility, use an "option list" or the ``.. program`` and ``.. option`` directives instead of a table.

Indentation
"""""""""""
Firstly, indentation of a text paragraph is rarely required. Doing so does not "link" the text with a heading in any way, text falls beneath a section or sub-section purely by being literally *beneath* that heading. When placing source code into a source code listing, any indentation may be used for the source code (typically should represent the actual preferred indentation of the code's respective project were it included in the repository), but to avoid ambiguity in indentation used in the documentation versus indentation used in the source code, all documentation indentation should be done using *one (1) TAB character* and **never spaces**.

Lists
"""""
When making a list, consider first what *kind* of list it is. A list only needs to be enumerated if the enumeration has meaning e.g. a list of hierarchically ordered data or a sequential list of steps to accomplish a task or desired state. When enumerating list items, use ``#.`` instead of literal numbers, as this will enumerate automatically which makes modification of the list much easier at a later date. Unordered lists may use ``-`` or ``*`` for each item. Lists do not need to be indented beyond the current paragraph level. If a list is an unordered list of terms and their definitions, use a "definition list" instead of any other kind of list. If a list is a list of fields in a document or object, use a "field list" instead of any other kind of list. If a list is a list of parameters, arguments, or flags used by a command-line utility, use an "option list" instead of any other kind of list.

Notes and Footnotes
"""""""""""""""""""
Instead of ``**NOTE**`` or similar, consider using the ``.. note`` directive, or one of the appropriate admonitions supported by :abbr:`RST (reStructuredText)`:

attention
	The default admonition that calls attention to something without any specific semantics or attached context. Use when none of the others seem appropriate.
caution
	Includes cautionary information. Should be used when advising the reader that the containing section includes instructions or information that frequently confuse people/trip people up, and how to avoid these pitfalls.
danger
	Advises the reader of potential security risks or system damage that could occur as a result of following instructions in the containing section, or as a result of making assumptions about and/or improperly utilizing information in the containing section.
error
	Denotes an error. This has limited uses in the :abbr:`ATC (Apache Traffic Control)` documentation.
hint
	Offers a hint to nudge readers in the direction of a solution. This has limited uses in the :abbr:`ATC (Apache Traffic Control)` documentation.
impl-detail
	Contains information describing the details of an implementation of a feature described in the containing section. For an example, see :ref:`the DSCP section of the Delivery Services page <ds-dscp>`.

	.. note:: This is an `extension`_ provided by the :abbr:`ATC (Apache Traffic Control)` documentation build system, and will not appear in the Sphinx project documentation nor the reStructuredText standard.

important
	Contains information that is important to consider while reading the containing section. Typically content that is related to a section's content in an important way should appear as content of that section, but if a section is in danger of readers "skimming" it for information this can be useful to catch their eye.
note
	Used to segregate content that is only tangentially related to the containing section, but is noteworthy nonetheless. Historically the most used admonition, containing caveats, exceptions etc.
tip
	Provides the reader with information that can be helpful - particularly to users/administrators/developers new to :abbr:`ATC (Apache Traffic Control)` - but not strictly necessary to understand the containing section.
warning
	Warns the reader about possible unintended consequences of following instructions/utilizing information in the containing section. If the behavior warned about constitutes a security risk and/or serious damage to one or more systems - including clients and origins - please use ``.. danger`` instead.

In a similar vein, instead of e.g. "(See also: some-link-or-reference)" please use the special ``.. seealso`` admonition. If the same admonition is required more than twice on the same page, it most likely ought to be a footnote instead. Footnotes should ideally use human-readable labels or, failing that, be labeled sequentially in the order of appearance. Footnotes should appear at the end of the major section in which they first or last appear. In practice, however, placement of the footnote is left to the writer's discretion.

Section Headings
""""""""""""""""
When deciding on the name of a section heading, it is important to select a name that both accurately reflects the content it contains and is suitable for reference later e.g. "Health Protocol" is good, but "Configuring Multi-Site Origin" as the title of a page which not only explains MSO configuration but also the concept is not good. Section headings follow a hierarchy, and for consistency's sake should follow this particular hierarchy:

#. Document title. There should only be one of these per page, and it should be the first heading on the page. This will also make the contained text the "Page Title" in HTML output (i.e. ``<title>Page Title</title>`` in the ``<head>``).

	.. code-block:: rst
		:caption: Document Title

		**************
		Document Title
		**************

#. Section header. This should represent a main topic of the page

	.. code-block:: rst
		:caption: Section Header

		Section Header
		==============

#. Subsection header. This should represent a key piece of a main topic on the page

	.. code-block:: rst
		:caption: Subsection Header

		Subsection Header
		-----------------

#. Sub-Subsection header. This represents a group of content logically separate from the rest of the subsection but still related to the content thereof. It is also acceptable to use this as an "aside" containing information only tangentially related to the subsection content.

	.. code-block:: rst
		:caption: Sub-Subsection Header

		Sub-Subsection Header
		"""""""""""""""""""""

#. Aside or Sub-Sub-Subsection header. This is the lowest denomination of header, and should almost always be used exclusively for "asides" which contain information only tangentially related to the sub-subsection content.

	.. code-block:: rst
		:caption: Aside

		Aside
		'''''

Section headings should *always* follow this order exactly, and **never** skip levels (which will generally cause a failure to compile properly). These can be thought of as the equivalents of the HTML tags ``<h1>`` through ``<h5>``, respectively. Sectioning should never require more specificity than can be provided by an "Aside". Please do not use **bold text** in lieu of a proper section heading. When referencing information in another section on the same page, please do not refer to the current placement of the referenced content relative to the referencing content. For example, instead of "as discussed below", use "as discussed in `Terms`_".

Terms
"""""
Please always spell out the entire name of any Traffic Control terms used in the definition. For example, a collection of :term:`cache servers` associated with a certain physical location is called a "Cache Group", not a "CG", "cachegroup", "cache location" etc. A subdomain and collection of :term:`cache servers` responsible collectively for routing traffic to a specific origin is called a :term:`Delivery Service`", not a "DS", "deliveryservice" etc. Similarly, always use *full* permissions role names e.g. "operations" not "oper". This will ensure the :ref:`glossary` is actually helpful. To link a term to the glossary, use the ``:term:`` role. This should be done for virtually every use of a Traffic Control term, e.g. ``:term:`Cache Group``` will render as: :term:`Cache Group`.
Generally speaking, be wary of using the word "cache". To most people that means the *actual* cache on a hard disk somewhere. This word is frequently confused with " :term:`cache server`", which - when accurate - is always preferred over "cache".

.. _api-doc-guidelines:

Documenting API Routes
----------------------
Follow all of the formatting conventions in `Formatting`_. Maintain the structural format of the API documentation as outlined in the :ref:`to-api` section. API routes that have variable paths e.g. :ref:`to-api-profiles-id` should use `mustache templates <https://mustache.github.io/mustache.5.html>`_ **not** the Mojolicious-specific ``:param`` syntax. This keeps the templates generic, familiar, and reflects the inability of a request path to contain procedural instructions or program logic. Please do not include the ``/api/2.x/`` part of the request path for Traffic Ops API endpoints. If an endpoint is unavailable prior to a specific version, use the ``.. versionadded`` directive to indicate that version. Likewise, do not make a new page for an endpoint when it changes across versions, instead call out the changes using the ``.. versionchanged`` directive. If an endpoint should not be used because newer endpoints provide the same functionality in a better way, use the ``.. deprecated`` directive to link to them and explain why they are better.

When documenting an API route, be sure to include *all* methods, request/response JSON payload fields, path parameters, and query parameters, whether they are optional or not. When describing a field in a JSON payload, remember that JSON does not have "hashes" it has "objects" or even "maps". When documenting path parameters such as :term:`Profile` :ref:`profile-id` in :ref:`to-api-profiles-id`, consider that the endpoint path cannot be formed without defining **all** path parameters, and so to label them as "required" is superfluous.

The "Response Example" must **always** exist. "TODO" is **not** an acceptable Response Example for new endpoints. The "Request Example" must only exist if the request requires data in the body (most commonly this will be for ``PATCH``, ``POST`` and ``PUT`` methods). It is, however, strongly advised that a request example be given if the endpoint takes Query Parameters or Path Parameters, and it is required if the Response Example is a response to a request that used a query or path parameter. If the Request Example *is* present, then the Response Example **must** be the appropriate response **to that request**. When generating Request/Response Examples, attempt to use the :ref:`ciab` environment whenever possible to provide a common basis and familiarity to new users who likely set up "CDN in a Box" as a primer for understanding CDNs/Traffic Control. Responses are sometimes hundreds of lines long, and in those cases only as much as is required for an understanding of the structure needs to be included in the example - along with a note mentioning that the output was trimmed. Also always attempt to place structure explanations before any example so that the content of the example can be understood by the reader (though in general the placement of a floating environment like a code listing is not known at compile-time). Whenever possible, the Request and Response examples should include the *complete HTTP stack*, which captures behavior like Query Parameters, Path Parameters and HTTP cookie operations like those used by e.g. :ref:`to-api-logs`. A few caveats to the "include all headers" rule:

- The ``Host`` header ought to reflect the actual hostname of the Traffic Ops server - which should be "trafficops.infra.ciab.test" for the CDN in a Box environment. This can be polluted when requests are made to a remotely running CDN in a Box on a different server.
- The "mojolicious" cookie is extremely long and potentially insecure to publicly show. As such, a placeholder should be used for its value, preferably "...".
- The ``Content-Type`` header sent by :manpage:`curl(1)` (and possibly others) is always ``application/x-www-form-urlencoded`` regardless of the actual content (unless overridden). Virtually all payloads accepted by the API must be JSON, so this should be modified to reflect that when appropriate e.g. ``application/json``.
- API output is often beautified by inserting line breaks and indentation, which will make the ``Content-Length`` header (if any) incorrect. Don't worry about fixing that - just try to leave the output as close as possible to what will actually be returned by leaving it the way it is.

File names should reflect the request path of the endpoint, e.g. a file for an endpoint of the Traffic Ops API ``/api/4.0/foo/{{fooID}}/bar/{{barID}}`` should be named ``foo_fooID_bar_barID.rst``. Similarly, reference labels linking to the document title for API route pages should follow the convention: ``<component>-api-<path>`` in all lowercase where ``<component>`` is an abbreviated Traffic Control component name e.g. ``to`` and ``<path>`` is the request path e.g. ``foo_bar``. So a label for an endpoint of the Traffic Ops API at ``/api/4.0/foo_bar/{{ID}}`` should be ``to-api-foo_bar-id``.

Extension
---------
The :abbr:`ATC (Apache Traffic Control)` documentation provides an extension to the standard roles and directives offered by Sphinx, located at :atc-file:`docs/source/_ext/atc.py`. It provides the following roles and directives:

impl-detail
	An admonition directive used to contain implementation-specific notes on a subject.

	.. code-block:: rst
		:caption: Example impl-detail usage

		.. impl-detail:: Implementation-specific information here.

	This example usage renders like so:

	.. impl-detail:: Implementation-specific information here.

atc-file
	Creates a link to the specified file on the ``master`` branch of the :abbr:`ATC (Apache Traffic Control)` repository. For example, "``:atc-file:`docs/source/development/documentation_guidelines```" renders as a link to the source of this documenting section like so: :atc-file:`docs/source/development/documentation_guidelines`. You can also link to directories as well as files.

issue
	A text role that can be used to easily link to GitHub Issues for the :abbr:`ATC (Apache Traffic Control)` repository. For example, "``:issue:`1```" renders as: :issue:`1`.
pr
	A text role that can be used to easily link to GitHub Pull Requests for the :abbr:`ATC (Apache Traffic Control)` repository. For example, "``:pr:`1```" renders as :pr:`1`.
pull-request
	A synonym for ``pr``

godoc
	A text role that can be used to easily link to the documentation for any Go package, type, or function/method (grouped constants/variables not supported). For example, "``:godoc:`net/http.HandlerFunc```" renders as :godoc:`net/http.HandlerFunc`.
atc-godoc
	This is provided for convenience, and is identical to ``:godoc:`` except that it is assumed to be relative to the Apache Traffic Control project. For example, ``:atc-godoc:`lib/go-rfc.MimeType.Quality``` renders as :atc-godoc:`lib/go-rfc.MimeType.Quality`.
to-godoc
	This is provided for convenience, and is identical to ``:godoc:`` except that it is assumed to be relative to the :atc-godoc:`traffic_ops/traffic_ops_golang` package. For example, ``:to-godoc:`api.Info``` renders as :to-godoc:`api.Info`.
