# Traffic Ops Config File / Snapshot Compare

This test allows you to compare all generated config files and CDN snapshots (CRConfig.json) from 2 instances of Traffic Ops. For example, you could compare config files / snapshots of a MySQL vs Postgres Traffic Ops. You could even compare across releases (1.7.0 vs 1.8.0).

*Prerequisites*

1. Make sure the data in your databases are synced to avoid getting false positives.
2. Queue updates for all servers in both instances.
3. Modify test.config with proper settings. Set perform_snapshot=1 if you want to force a snapshot in both instances.

*Running the Test*

1. `./cfg_test.pl getref test.config` your ref files go into `/tmp/files/ref`
2. `./cfg_test.pl getnew test.config` your new files go into `/tmp/files/new`
3. `./cfg_test.pl compare test.config` - all `not ok` lines should be explained.

It will compare all files for all profiles, _including_ the CRConfig.json. 

