<!--
    Licensed to the Apache Software Foundation (ASF) under one
    or more contributor license agreements.  See the NOTICE file
    distributed with this work for additional information
    regarding copyright ownership.  The ASF licenses this file
    to you under the Apache License, Version 2.0 (the
    "License"); you may not use this file except in compliance
    with the License.  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing,
    software distributed under the License is distributed on an
    "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
    KIND, either express or implied.  See the License for the
    specific language governing permissions and limitations
    under the License.
-->


# Rpm Build Instructions

##  Using `docker-compose`

We are moving toward using a docker-based build system.  This eliminates the need to maintain a local installation with all the
build tools as well as ensuring that you are using the same versions as we use for testing.

These are the versions of these tools we are using:
* docker 1.12.2
* docker-compose 1.8.1

You can build from any repository/branch combination,  but that repository must be available to `git clone ...`.  Note that in the
following `docker-compose` commands, you can limit building to one or more sub-projects by supplying arguments at the end.  If none
are supplied,  then *all* will be run.

Starting at the top-level of your trafficcontrol git clone (e.g. `~/src/trafficcontrol`):

> cd infrastructure/docker/build
> docker-compose build traffic_ops_build traffic_monitor_build ...
> GITREPO=https://github.com/username/trafficcontrol BRANCH=mybranch docker-compose up traffic_ops_build traffic_monitor_build ...

The resulting `.rpm` files will be created in the `artifacts` directory.


## Building the old-fashioned way

rpm files for all sub-projects can be built using the file `build/build.sh`.  If this script is given parameters, it will build only
those projects specified on the command line, e.g.  `$ ./build/build.sh traffic_ops`.  The prerequisites for each sub-project are
listed below.

These build scripts depend on the text in the __VERSION__ file along with the __BUILD_NUMBER__ described below to name each rpm.

The build scripts use environment variables to control how the build is done.  These have sensible defaults listed below, and it is
recommended to not override them:
* __WORKSPACE__
   - defaults to the top level of the traffic_control directory.  The _dist_ and _rpmbuild_ directories are created in this
     directory during the rpm build process.
* __BUILD_NUMBER__
   - generates build number from the number of commits on the current git branch followed by the 8 character short commit hash of
     the last commit on the branch.This number is used to create the rpm version, e.g. _traffic_ops.1.2.0.1723.a18e2bb7_.

At the conclusion of the build,  all rpms are copied into the __$WORKSPACE/dist__ directory.

## Prerequisites for building:

### all sub-projects

* CentOS 6.x
* rpmbuild (yum install rpm-build)
* git 1.7.12 or higher

#### traffic_ops:
* go 1.7 or higher

#### traffic_stats:
* go 1.7 or higher

#### traffic_monitor and traffic_router:
* java-1.8.0-openjdk and java-1.8.0-openjdk-devel
* apache-maven 3.3.1 or higher

#### traffic_monitor_golang:
* go 1.7 or higher

#### traffic_portal
* npm (yum install npm)
  * grunt (npm install -g grunt)

# Docker build instructions

__Building using `docker` is experimental at this time and has not been fully vetted.__

Dockerfiles for each sub-project are located in the build directory (e.g. `traffic_ops/build/Dockerfile`)

## Optionally set these environment variables to control the source to start with:
* `GITREPO` (default is `https://github.com/Comcast/traffic_control`) and `BRANCH` (default is master).

> export GITHUB_REPO=https://github.com/myuser/traffic_control
> export BRANCH=feature/my-new-feature
> ./build/docker-build.sh
