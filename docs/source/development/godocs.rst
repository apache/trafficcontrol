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

******
Godocs
******

Linking to Godocs
=================

As mentioned in :ref:`docs-guide`, you can use the ``:godoc:`` role, the ``:atc-godoc:`` role, and the ``:to-godoc:`` role to link to Godocs from the ATC documentation.

Keeping pkg.go.dev from hiding the Apache Traffic Control Godocs
================================================================
If less than 75% of ATC's :atc-file:`LICENSE` file contains OSS licenses according to `pkg.go.dev <https://pkg.go.dev/>`_ license detection system pkg.go.dev will hide all of ATC's Godocs for that particular ATC version, saying "Documentation not displayed due to license restrictions.". Example: https://pkg.go.dev/github.com/apache/trafficcontrol@v5.1.0+incompatible/lib/go-tc

When listing a dependency in the LICENSE file as part of `bundling a new dependency <https://infra.apache.org/licensing-howto.html#permissive-deps>`_, make sure that the license "pointer" adheres to this format:

.. code-block:: cpsa
	:caption: The ``atc-dependency.lre`` license exception from :godoc:`golang.org/x/pkgsite`

	This product bundles __4__, which
	(( is || are ))
	available under
	(( a || an ))
	(( Apache-2.0 || BSD-2-Clause || BSD-3-Clause || MIT ))
	license.
	__15__
	(( /* || .css || .js || .scss ))
	(( ./licenses/__4__ || ./vendor/__16__/LICENSE
	(( .libyaml || .md || .txt ))??
	))
	Refer to the above license for the full text.

Example:

.. code-block::
	:caption: An example of a bundled dependency license pointer

	This product bundles go-acme/lego, which is available under an MIT license.
	@vendor/github.com/go-acme/lego/*
	./vendor/github.com/go-acme/lego/LICENSE
	Refer to the above license for the full text.

The ATC repository includes a GitHub Actions workflorkflow (:atc-file:`.github/workflows/license-file-coverage.yml`) to verify that changes to the LICENSE file will not result in pkg.go.dev hiding the Apache Traffic Control Godocs.
