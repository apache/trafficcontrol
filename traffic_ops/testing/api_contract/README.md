# Traffic Ops API Contract Tests

The Traffic Ops API Contract tests are used to validate the Traffic Ops API's.

## Setup

In order to run the tests you will need a running instance of Traffic Ops and Traffic Ops DB:

1. **Traffic Ops Database** configured port access
    - _Usually 5432 - should match the value set in database.conf and the **trafficOpsDB port** in traffic-ops-test.conf_
2. **Traffic Ops** configured port access
    - _Usually 443 or 60443 - should match the value set in cdn.conf and the **URL** in traffic-ops-test.conf_
3. Running Postgres instance with a database that matches the **trafficOpsDB dbname** value set in traffic-ops-test.conf
    - For example to set up the `to_test` database do the following:

         ```console
         $ cd trafficcontrol/traffic_ops/app
         $ db/admin --env=test reset
         ```

      The Traffic Ops users will be created by the tool for accessing the API once the database is accessible.

      To test if `db/admin` ran all migrations successfully, you can run the following command from the `traffic_ops/app` directory:

        ```console
        db/admin -env=test dbversion
        ```
      The result should be something similar to:
        ```
        dbversion 2021070800000000
        ```
      If migrations did not run successfully, you may see:
        ```
        dbversion 20181206000000 (dirty)
        ```
      Make sure **trafficOpsDB dbname** in traffic-ops-test.conf is set to: `to_test`

      For more info see: http://trafficcontrol.apache.org/docs/latest/development/traffic_ops.html?highlight=reset

4. A running Traffic Ops Golang instance pointing to the `to_test` database.

    ```console
	$ cd trafficcontrol/traffic_ops/traffic_ops_golang
    $ cp ../app/conf/cdn.conf $HOME/cdn.conf 
    $ go build && ./traffic_ops_golang -cfg $HOME/cdn.conf -dbcfg ../app/conf/test/database.conf
    ```
   Verify that the passwords defined for your `to_test` database match:
    - `trafficcontrol/traffic_ops/app/conf/test/database.conf`
    - `traffic-ops-test.conf`

5. Install the requirements for testing API contract tests

    ```console
    pip install -r /path/to/requirements.txt
    ``` 

## Running the API Contract Tests

The API Contract tests are run using `pytest` from the **traffic_ops/testing/api_contract/v4** directory

Example commands to run the tests:

Only Test a specific endpoint with Arguments
> Note: For staging and nightly environments (Steps 1-4 Not mandatory)

```console
pytest --to_user Username--to_password Password --to_url URL test_cdns.py
```

Only Test a specific endpoint with Local Traffic Ops Instance
> Note: For local testing environment (Steps 1-4 mandatory)
```console
pytest -rA test_cdns.py
```
