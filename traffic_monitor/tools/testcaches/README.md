# testcaches

The `testcaches` tool simulates multiple ATS caches' `_astats` endpoints.

Its primary goal is for testing the Monitor under load, but it may be useful for testing other components.

A list of parameters can be seen by running `./testcaches -h`. There are only three: the first port to use, the number of ports to use, and the number of remaps (delivery services) to serve in each fake server.

Each port is a unique fake server, with distinct incrementing stats.

When run with no parameters, it defaults to ports 40000-40999 and 1000 remaps.

Stats are served at the regular ATS `stats_over_http` endpoint, `_astats`. For example, if it's serving on port 40000, it can be reached via `curl http://localhost:40000/_astats`. It also respects the `?application=system` query parameter, and will serve only system stats (the Monitor "health check" [as opposed to the "stat check"]). For example, `curl http://localhost:40000/_astats?application=system`.
