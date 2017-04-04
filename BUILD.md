# Building Traffic Control Components

## Build using docker-compose

This is the easiest way to build all the components of Traffic Control; all requirements
are automatically loaded into the image used to build each component.

### Requirements
- `docker` (https://docs.docker.com/engine/installation/)
- `docker-compose` (https://docs.docker.com/compose/install/)

### Steps

From the top level of the incubator-trafficcontrol directory.  The source in
the current directory is used for the process.   One or more components (with
\_build suffix added) can be added on the command line.

Clean up any previously-built docker containers:
> $ docker-compose -f infrastructure/docker/build/docker-compose.yml down -v

And images:
> $ docker images | awk '/traffic\_.*\_build/ { print $3 }' | xargs docker rmi -f

Create and run new build containers:
> $ docker-compose -f infrastructure/docker/build/docker-compose.yml up [ container name ...] 

Container names can be one or more of these:
* `source`  (builds only the source tarball)
* `traffic_monitor_build`
* `traffic_ops_build`
* `traffic_portal_build`
* `traffic_router_build`
* `traffic_stats_build`

If no component names are provided on the command line, all components will be built.

All rpms are copied to `dist` at the top level of the `incubator-trafficcontrol` directory.
