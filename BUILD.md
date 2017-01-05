# Building Traffic Control Components

## Build using docker-compose

This is the easiest way to build all the components of Traffic Control; all requirements
are automatically loaded into the image used to build each component.

### Requirements
- `docker` (https://docs.docker.com/engine/installation/)
- `docker-compose` (https://docs.docker.com/compose/install/)

### Steps

From the top level of the incubator-trafficcontrol directory.  Use the BRANCH
environment variable to specify the version of Traffic Control to build.   One
or more components (with \_build suffix added) can be added on the command
line:

> $ BRANCH=1.8.x docker-compose -f infrastructure/docker/build/docker-compose.yml up traffic\_monitor\_build

If no component names are provided here, all components will be built.


