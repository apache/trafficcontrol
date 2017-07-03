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

# Create a TO RPM with Dependencies

1. Download repo Traffic Control (or from your favorite repo)
```
$ git clone http://github.com/Comcast/traffic_control.git
```
2. Bring up Vagrant environment (http://www.vagrantup.com)
```
$ cd <repo dir>
$ cp traffic_control/traffic_ops/build/Vagrantfile ./
$ vagrant up
```
3. ssh into vagrant environment
```
$ vagrant ssh
```
4. **OPTIONAL** Set environment variables to control build.  All are
   automatically set to reasonable defaults (in parentheses) and it is
   recommended to leave them unset.   They can be overridden if necessary:
   - *BRANCH* (master)
   - *HOTFIX\_BRANCH* (none)
   - *WORKSPACE* (top level of local repo tree -- workspace must be a clone of the repository)
   - *BUILD\_NUMBER* (# of commits in branch + last commit identifier)
5. Build the RPM
```
$ cd /vagrant/traffic_control/traffic_ops/build
$ ./build_rpm.sh
```
Notes:  
This is known to work with CentOS 6.7 as the Vagrant environment.
