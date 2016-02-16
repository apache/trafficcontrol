# Traffic Monitor

### Why Tests are not in exactly matching packages

The "com.comcast.cdn.traffic_control.traffic_monitor" portion of the package name was omitted from unit
tests to prevent improper referencing of package private fields and methods of the code under test.

### Running Traffic Monitor locally

The "com.comcast.cdn.traffic_control.traffic_monitor.Start" class allows one to run Traffic Monitor
locally provided that necessary configuration is in place. By default, the files are specified
with paths relevant to certain IDEs, but these paths can be changed by specifying different
properties via System.properties. These properties are:

* traffic_monitor.path.config
* traffic_monitor.path.db

The first property refers to the location of traffic_monitor_config.js. The second property
refers to the directory that will be used for certain data files that are downloaded at runtime.
If you need to specify a different path, use the -D option to the Java command, or modify the
paths in the Start class directly.
