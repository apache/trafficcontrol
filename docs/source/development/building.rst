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

Building Traffic Control
========================


Build using pkg
---------------

This is the easiest way to build all the components of Traffic Control;
all requirements are automatically loaded into the image used to build
each component.

Requirements
~~~~~~~~~~~~

-  ``docker`` (https://docs.docker.com/engine/installation/)
-  ``docker-compose`` (https://docs.docker.com/compose/install/)
   (optional, but recommended)

If ``docker-compose`` is not available, the ``pkg`` script will
automatically download and run it in a container. This is noticeably
slower than running it natively.

Usage
~~~~~

::

    $ ./pkg -?
    Usage: ./pkg [options] [projects]
        -q      Quiet mode. Supresses output.
        -v      Verbose mode. Lists all build output.
        -l      List available projects.

        If no projects are listed, all projects will be packaged.
        Valid projects:
                - traffic_portal_build
                - traffic_router_build
                - traffic_monitor_build
                - source
                - traffic_ops_build
                - traffic_stats_build


If any project names are provided on the command line, only those will be built.
Otherwise, all projects are built.

All artifacts (rpms, logs, source tar ball) are copied to ``dist`` at the top level of the
``incubator-trafficcontrol`` directory.

Example
~~~~~~~

::

    $ ./pkg -v source traffic_ops_build
    Building source.
    Building traffic_ops_build.

Build using docker-compose
--------------------------

If the ``pkg`` script fails, ``docker-compose`` can still be used directly.

Usage
~~~~~

::

    $ docker-compose -f ./infrastructure/docker/build/docker-compose.yml down -v
    $ docker-compose -f ./infrastructure/docker/build/docker-compose.yml up --build source traffic_ops_build
    $ ls -1 dist/
    build-traffic_ops.log
    traffic_ops-2.1.0-6396.07033d6d.el7.src.rpm
    traffic_ops-2.1.0-6396.07033d6d.el7.x86_64.rpm
    traffic_ops_ort-2.1.0-6396.07033d6d.el7.src.rpm
    traffic_ops_ort-2.1.0-6396.07033d6d.el7.x86_64.rpm
    trafficcontrol-incubating-2.1.0.tar.gz
