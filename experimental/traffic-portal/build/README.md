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

# Traffic Portal Installation / Configuration

### 1. Build Traffic Portal RPM w/ Docker

* Build instructions: https://github.com/apache/trafficcontrol/blob/master/BUILD.md

### 2. Install

* Install the Node.js JavaScript runtime

    ```
    $ curl --silent --location https://rpm.nodesource.com/setup_16.x | sudo bash -
    $ sudo yum install -y nodejs
    ```

* Install the Traffic Portal RPM

    ```
    $ sudo yum install -y traffic_portal-[version]-[commits].[sha].x86_64.rpm
    ```

### 3. Configure

* Configure Traffic Portal

    ```
    $ sudo vim /etc/traffic_portal/conf/config.js (read the inline comments)
    ```

### 4. Run

* Start Traffic Portal

    ```
    $ sudo service traffic_portal start
    ```

* Navigate to Traffic Portal

    ```
    $ http(s)://ip-address:port
    ```

### Notes
- Traffic Portal consumes the Traffic Ops API, therefore, an instance of Traffic Ops must be running.
- This is known to work with CentOS 7 and Centos 8 as the host environment.
