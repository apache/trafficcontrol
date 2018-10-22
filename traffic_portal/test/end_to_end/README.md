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

# UI Tests  (Using Protractor)

- Install and run protractor
  http://www.protractortest.org/#/tutorial
    
- Start up Selenium Server
    webdriver-manager start

- Make sure Traffic Portal is running

- Edit conf.js if necessary to match your environment.

- Run your tests
  protractor conf.js --params.adminUser 'user' --params.adminPassword 'password'

- NOTE: Errors with webdriver
  Most errors with webdriver can be remedied by running the following:
  webdriver-manager clean
  webdriver-manager update

