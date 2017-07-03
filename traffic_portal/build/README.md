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

# Traffic Portal Installation

### 1. Build RPM w/ Docker

* See https://github.com/apache/incubator-trafficcontrol/blob/master/BUILD.md

### 2. Install

* Install the Node.js JavaScript runtime

    ```
    $ curl --silent --location https://rpm.nodesource.com/setup_6.x | sudo bash -
    $ sudo yum install -y nodejs
    ```

* Install the Traffic Portal RPM

    ```
    $ sudo yum install -y traffic_portal-[version]-[commits].[sha].x86_64.rpm
    ```

### 3. Configure

* Configure Traffic Portal

    ```
    $ sudo vi /etc/traffic_portal/conf/config.js (read the inline comments)
    ```

### 4. Run

* Start Traffic Portal

    ```
    $ sudo service traffic_portal start
    ```

* Navigate to Traffic Portal

    ```
    $ http://localhost:8080
    ```

#### Notes

    - Traffic Portal consumes the Traffic Ops API, therefore, an instance of Traffic Ops must be running.
    - This is known to work with CentOS 6.7 and Centos 7 as the host environment.
