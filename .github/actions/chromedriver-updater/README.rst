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

********************
chromedriver-updater
********************

Assigns the GitHub Triage role to non-committers who have fixed 2 Issues in the past month.

Environment Variables
=====================

+----------------------------+----------------------------------------------------------------------------------+
| Environment Variable Name  | Value                                                                            |
+============================+==================================================================================+
| ``GITHUB_TOKEN``           | Required. ``${{ github.token }}`` or ``${{ secrets.GITHUB_TOKEN }}``             |
+----------------------------+----------------------------------------------------------------------------------+
| ``PR_GITHUB_TOKEN``        | Required. ``${{ github.token }}`` or another token                               |
+----------------------------+----------------------------------------------------------------------------------+
| ``GIT_AUTHOR_NAME``        | Optional. The username to associate with the commit that updates the Go version. |
+----------------------------+----------------------------------------------------------------------------------+

Outputs
=======

``exit-code``
-------------

Exit code is 0 unless an error was encountered.

Inputs
======
Optionally takes a file that contains a line for each project update of the form: `{project path}:{old version},{new version}`

Example usage
=============

.. code-block:: yaml

    - name: Update Chromedriver Versions
      run: python3 -m chromedriver_updater
      env:
        GIT_AUTHOR_NAME: asf-ci-trafficcontrol
        GITHUB_TOKEN: ${{ github.token }}
