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

.. _dev-building:

************************
Building Traffic Control
************************
The build steps for Traffic Control components are all pretty much the same, despite that they are written in a variety of different languages and frameworks. This is accomplished by using Docker.

.. note:: Currently, both listed methods of building Traffic Control components will produce ``*.rpm`` files, meaning that the support of these components is limited to RedHat-based distributions - and none of them are currently tested (or guaranteed to work) outside of CentOS7, specifically.

.. _pkg:

Build using ``pkg``
===================
This is the easiest way to build all the components of Traffic Control; all requirements are automatically loaded into the image used to build each component.

Requirements
------------
-  ``docker`` (https://docs.docker.com/engine/installation/)
-  ``docker-compose`` (https://docs.docker.com/compose/install/)
   (optional, but recommended)

If ``docker-compose`` is not available, the ``pkg`` script will automatically download and run it in a container. This is noticeably slower than running it natively.

Usage
-----
``./pkg [options] [projects]``

Options:

-q      Quiet mode. Suppresses output.
-v      Verbose mode. Lists all build output.
-l      List available projects.

If the ``projects`` argument is given, it should be one or more project names. If it is not given, all projects will be built. Valid projects:

- traffic_portal_build
- traffic_router_build
- traffic_monitor_build
- source
- traffic_ops_build
- traffic_stats_build

Output ``*.rpm`` files, build logs and source tarballs will be output to the ``dist`` directory at the root of the Traffic Control repository directory.

Build Using ``docker-compose``
==============================
If the ``pkg`` script fails, ``docker-compose`` can still be used to build the projects directly. The compose file can be found at ``infrastructure/docker/build/docker-compose.yml`` under the repository's root directory. It is recommended that between builds ``docker-compose down -v`` is run to prevent caching of old build steps. The service names are the same as the project names described above in `Usage`_, and similar to the ``pkg`` script, the build results, logs and source tarballs may all be found in the ``dist`` directory after completion.

Building Individual Components
==============================
Each Traffic Control component can be individually built, and the instructions for doing so may be found in their respective component's development documentation.

Building This Documentation
---------------------------
This documentation uses the `Sphinx documentation build system <http://www.sphinx-doc.org/en/master/>`_, and as such requires a Python3 version that is at least 3.4.1, but no greater than 3.6.5\ [1]_. It also has dependency on Sphinx, and Sphinx extensions and themes. All of these can be easily installed using `pip` by referencing the requirements file like so:

.. code-block:: shell
	:caption: Run from the Repository's Root Directory

	python3 -m pip install --user -r docs/source/requirements.txt

Once all dependencies have been satisfied, build using the Makefile at ``docs/Makefile``.

Alternatively, it is also possible to :ref:`pkg`.

.. [1] A bug in `pygments` prevents Python 3.7.x from working with Sphinx
