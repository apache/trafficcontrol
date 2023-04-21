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

# Traffic Ops API contract tests

The Traffic Ops API contract tests are used to validate the Traffic Ops API's.

## Setup

In order to run the tests you will need a running instance of Traffic Ops and Traffic Ops DB:

1. Follow the instructions for building and running Traffic Ops from docs.

2. Make sure you don't have the published Apache TrafficControl

   ```console
   pip uninstall Apache-TrafficControl
   ```

3. Install packages via setup.py under traffic_control/clients/python/setup.py
   
   ```console
   sudo python3 setup.py install
   ```

4. Install local Apache Traffic Control under trafficcontrol root directory

   ```console
   pip install traffic_control/clients/python
   ```

5. Make sure a build is generated under traffic_control/clients/python

6. set the PYTHONPATH environment variable to that directory (replacing the absolute path with your own):

   ```console
   export PYTHONPATH=/absolute/path/to/your/repo/trafficcontrol/traffic_control/clients/python/build/lib
   ```

7. Install the requirements under traffic_ops/testing/api_contract/v4/requirements.txt.

    ```console
    pip install -r /path/to/requirements.txt
    ``` 

## Running the API contract tests

The API contract tests are run using `pytest` from the ``traffic_ops/testing/api_contract/v4`` directory.

Example commands to run the tests:

Only test a specific endpoint with arguments.
> Note: For particular environments (Step 1 is not mandatory)
```console
pytest --to-user Username --to-password Password --to-url URL test_cdns.py
```

Only test a specific endpoint with local Traffic Ops instance.
> Note: For local environment (Step 1 is mandatory)
```console
pytest -rA test_cdns.py
```
