# Traffic Ops Config Test

This test allows you to check all generated config files from Traffic Ops. 

*How to Test*

# Make sure no changes are made in the reference system, and all changes are "snapped" and "queued".
# Get a copy of the reference DB (using tools->DB Dump), and save it. 
# Get the files from the reference system generates by running `./cfg_test.pl getref test.config` - make sure `test.config` is right. It will prompt you for the user's passwd, and get all files for all profiles into `/tmp/files/ref`.
# Load the database into your new system. run migrations (includeing mysql -> postgres) , and move the riak data to your test system from the ref system. 
# `./cfg_test.pl getnew test.config` your new files go into `/tmp/files/new`
# `./cfg_test.pl compare test.config` - all `not ok` lines should be explained.

It will compare all files for all profiles, _including_ the CRConfig.json. 

