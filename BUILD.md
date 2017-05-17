# Building Traffic Control Components

## Build using docker-compose

This is the easiest way to build all the components of Traffic Control; all requirements
are automatically loaded into the image used to build each component.

### Requirements
- `docker` (https://docs.docker.com/engine/installation/)
- `docker-compose` (https://docs.docker.com/compose/install/) (optional, but recommended)

If `docker-compose` is not available, the `pkg` script will automatically download
and run it in a container. This is noticeably slower than running it natively.

### Steps

From the top level of the incubator-trafficcontrol directory.  The source in
the current directory is used for the process.   One or more components (with
\_build suffix added) can be added on the command line.

This is all run automatically by the `pkg` script at the root of the repository.

    $ ./pkg -?
    Usage: ./pkg [options] [projects]
        -q      Quiet mode. Supresses output.
        -v      Verbose mode. Lists all build output.
        -l      List available projects.

        If no projects are listed, all projects will be packaged.
        Valid projects:
                - traffic_portal_build
                - traffic_router_build
                - traffic_monitor_build
                - source
                - traffic_ops_build
                - traffic_stats_build

If no component names are provided on the command line, all components will be built.

All rpms are copied to `dist` at the top level of the `incubator-trafficcontrol` directory.
