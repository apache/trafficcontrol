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

# Traffic Ops ORT Tests

The ORT tests are used to validate the ORT tools used with a
release of Traffic Ops.  The tests ensure that the ORT tools
are able to communicate with Traffic Ops, install required
packages and generate the correct ATS configuration files on
a Trafficserver cache.

The first thing you need do is to provide the Traffic Ops and
Traffic Ops ORT RPM files for the build you wish to test using 
this framework.  This test environment provides the necessary
Docker containers used to support and execute the tests when
given the necessary RPM's.  If you choose to not use the provided
Docker containers you will need to provide the following resources:

  - A running Traffic Ops with the installed release to be tested
  - A running Postgres SQL database which is loaded with the proper
    test data for the Traffic Ops release.
  - An Apache Traffic Server host that has the installed release of
    ORT to be tested against the release of Traffic Ops.
  - A yum server configured to provide the test rpm's herein.
  - A Traffic vault server.

# Directory layout

  - trafficcontrol/traffic_ops_ort/testing/docker:  has all the 
    necessary files for running the test Docker containers.
  - trafficcontrol/traffic_ops_ort/testing/ort-tests:  this directory.
    contains all the go files used to run the ORT tests.

# Setup.

  1.  Build the Traffic Ops and Traffic Ops ORT RPM's that you wish
      to test.  See the top level 'build' directory for building 
      instructions. 
  2.  Copy the Traffic Ops RPM to docker/traffic_ops/traffic_ops.rpm
      (NOTE:  Use the file name 'traffic_ops.rpm')
  3.  Copy the Traffic Ops ORT rpm to docker/ort_test/traffic_ops_ort.rpm
      (NOTE:  Use the file name 'traffic_ops_ort.rpm'
  4.  Copy an Apache Trafficserver rpm to 
      docker/yumserver/test-rpms/trafficserver-$VERSION.$COMMIT-HASH.el7.x86_64.rpm
      You will need to edit and adjust the trafficserver package value in
      ort-tests/tc-fixtures.json to match the $VERSION.$COMMIT-HASH used in the name
      of your RPM.  Search for 'CHANGEME' in the ort-tests/tc-fixtures.json 
      and change the value to match the RPM version you choose to use.
      For example, '8.0.8-19.77cb23a.el7.x86_64' and therefore the rpm file expected in 
      'docker/yumserver/test-rpms' is 'trafficserver-8.0.8-19.77cb23a.el7.x86_64.rpm'.
  6.  The container Docker files have the usernames and passwords used in the various
      containers ie, postgresql db, traffic_ops, and traffic_ops_ort.  The usernames
      and passwords passed to the 't3c' executable in in the 
      ort-tests/conf/docker-edge-cache.conf file.  Make sure that the usernames/passwords
      in the Docker files match those in the t3c configuration file.
      An example ort-tests/conf/edge-cache.conf file is provided should you choose to
      use your own Traffic Ops and Postgresql environment.
  7.  Build the Docker images and run the ort test:
      ``` 
      cd trafficcontrol/traffic_ops_ort/testing/docker
      docker-compose build
      docker-compose run ort_test
      ```
      After some time, test results should be available at
      'ort-tests/test.log'
  
  If you wish to run the tests manually use 'docker ps' to obtain the container id for
  the ort_test host and then:

  ```
     docker exec -it $ort_test /bin/bash -l
     cd /ort-tests
     go test -cfg=conf/docker-edge-cache.conf
  ```

  If you wish to run the tests manually using your own environment, create a config
  file with the necessary login information in the 'conf' directory and then rerun
  the tests using your config file.  WARNING: the traffic ops database will be dropped
  and initialized using the data in tc-fixtures.json and then the tests are run.  
  DO NOT USE a production Traffic Ops database with these test scripts.


 
