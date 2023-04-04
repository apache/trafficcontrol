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

***************
pr-to-update-go
***************

Opens a PR if a new minor Go revision is available.

For example, if the ``GO_VERSION`` contains ``1.14.7`` but Go versions 1.15.1 and 1.14.8 are available, it will

1. Create a branch named ``go-1.14.8`` to update the repo's Go version to 1.14.8
2. Updates all golang.org/x/ dependencies of the project, since these are meant to be updated with the Go compiler.
3. Open a PR targeting the ``master`` branch from branch ``go-1.14.8``

Other behavior in this scenario:

- If a branch named ``go-1.14.8`` already exists, no additional branch is created.
- If a PR titled *Update Go version to 1.14.8* already exists, no additional PR is opened.

Environment Variables
=====================

+----------------------------+----------------------------------------------------------------------------------+
| Environment Variable Name  | Value                                                                            |
+============================+==================================================================================+
| ``GIT_AUTHOR_NAME``        | Optional. The username to associate with the commit that updates the Go version. |
+----------------------------+----------------------------------------------------------------------------------+
| ``GITHUB_TOKEN``           | Required. ``${{ github.token }}`` or ``${{ secrets.GITHUB_TOKEN }}``             |
+----------------------------+----------------------------------------------------------------------------------+
| ``PR_GITHUB_TOKEN``        | Optional. A GitHub token other than GitHub Actions so that Actions will run on   |
|                            | the generated Pull Request                                                       |
+----------------------------+----------------------------------------------------------------------------------+
| ``GO_VERSION_FILE``        | Required. The file in the repo containing the version of Go used by the repo.    |
+----------------------------+----------------------------------------------------------------------------------+


Outputs
=======

``exit-code``
-------------

Exit code is 0 unless an error was encountered.

Example usage
=============

.. code-block:: yaml

	- name: PR to Update Go
	  run: python3 -m pr_to_update_go
	  env:
	    GIT_AUTHOR_NAME: asf-ci-trafficcontrol
	    GITHUB_TOKEN: ${{ github.token }}
	    GO_VERSION_FILE: GO_VERSION

Tests
=====

To run the unit tests:

.. code-block:: shell

	python3 -m unittest discover ./tests

To run the doctests:

.. code-block:: shell

	python3 ./pr_to_update_go/go_pr_maker.py
