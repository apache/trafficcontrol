# Traffic Monitor

### Why Tests are not in exactly matching packages

The "com.comcast.cdn.traffic_control.traffic_monitor" portion of the package name was omitted from unit
tests to prevent improper referencing of package private fields and methods of the code under test.