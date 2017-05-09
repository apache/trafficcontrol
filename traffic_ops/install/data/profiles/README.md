
# Profiles:
 
- This directory contains sample profiles that can be loaded in to a new
  instance of Traffic Ops.  
  
## Load Sample Profiles

### Prerequisites:

  * Postgres database is setup
  * Traffic Ops RPM is installed and configured to point to the Postgres Database
  * Traffic Ops `/opt/traffic_ops/install/bin/postinstall` has been run to create the tables, seed with data, and tables migrated

- To build the profiles to be loaded above one can `export` profiles from a Traffic Ops Instance 
 with production data

  After a successful `postinstall` execution.  Run the following commands from the new Traffic Ops instance
  ```
  $ cd /opt/traffic_ops/install/data/profiles
  $ ./load_profiles.sh
  ```
----
## Export Sample Profiles

### Prerequisites:

  * Postgres database with production data
  * Traffic Ops RPM is installed and configured to point to the appropriate Postgres database
  * Verify that the `$lookupTable` variable matches the `latest_` profiles in your production database with the `export_profiles.pl` script.


  After a successful `postinstall` execution.  Run the following command from the new Traffic Ops instance
  ```
  $ cd /opt/traffic_ops/install/data/profiles
  $ ./export_profiles.pl
  ```

