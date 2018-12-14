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

****************
Traffic Ops - Go
****************
Traffic Ops is currently in the process of being re-written in Go, which has its own dependencies, build process and source tree. To develop for the Go version of Traffic Ops you **must have cloned the Traffic Control repository into** ``$GOPATH/src/github.com/apache/trafficcontrol/`` **or linked your downloaded repository to that path. Either way, you must build from there**.

Building
========
To fetch dependencies, it should suffice to run ``go get -v`` from the ``traffic_ops/traffic_ops/golang`` directory. Then simply run ``go build .`` from the same directory to actually run the build\ [1]_.

.. [1] This should also run the ``go get`` command if this has not already been done, but it is recommended to do this in two different steps for ease of debugging.
