#!/usr/bin/env bash

# Get rid of exploded war file directory, but leave everything else untouched because this runs after the install
# steps during an upgrade (http://www.ibm.com/developerworks/library/l-rpm2/)
rm -rf /opt/traffic_monitor/webapps/ROOT