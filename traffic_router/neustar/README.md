## Neustar Geolocation Provider

This module provides a bean "neustarGeolocationService" that implements the GeolocationService interface defined
in the geolocation maven submodule. You can configure delivery services in traffic ops to use this module.

The default build does not include this into the Traffic Router war. You must add the 'neustar' profile id to your
maven build command. This module depends on the Neustar bff-reader library. See below for more details.

### Dependencies

This module depends on
 * Neustar Database Reader library bff-reader
 * Neustar Database files

## Getting the Neustar Database Reader library

Contact http://www.neustar.biz

### Installing the Neustar Database Reader Library

Run the following maven command:

`mvn install:install-file -Dfile=bff-reader-1.1.0.jar -DgroupId=com.quova.bff -DartifactId=bff-reader -Dversion=1.1.0 -Dpackaging=jar`


