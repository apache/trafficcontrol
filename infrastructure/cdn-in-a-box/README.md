# CDN-In-A-bOx (CIAO) - Notes

## Build step
So far, the build is done from this directory (/infrastructure/cdn-in-a-box/), using the following command:

```bash
docker-compose -f ./docker-compose.yml -f traffic_ops/docker-compose.yml -f traffic_vault/docker-compose.yml up --build
```

This will build and start all of the currently-implemented pieces of CIAO. By default, it leaves your terminal's stdout open to the logs produced by various parts of the CDN.


## <a name="ports"></a> Ports and Interfaces
CIAO provides API gateways and user interfaces over HTTP(S) for the pieces of the CDN that support it. By default, all of these services bind to the local address `0.0.0.0`, and are spread out over various ports. Here's a list of the ports and the services they provide:

* 443 - Exposes the Traffic Ops landing page, from which you may log in and view/manage/manipulate your CIAO at will (See [the Traffic Ops documentation](http://traffic-control-cdn.readthedocs.io/en/latest/admin/traffic_ops/using.html) for details). This is also the endpoint for the Go-based Traffic Ops API, and will handle all REST API requests (including those intended to be processed by Perl).

* 5432 - Exposes a connection to the Traffic Ops's Postgres database. Note that this is *not* meant to be used through a browser, and instead is an endpoint for direct connection to the database (this is used by Traffic Ops to store information).

* 8080 - The "adminer" portal which allows access to the Traffic Ops Database through a browser.

* 60443 - This is the endpoint for the old Perl-based API for Traffic Ops. Rather than send requests here, you should send them to port 443, as the Go-based API there will act as a reverse proxy to pass off API requests it doesn't directly handle back to this port.

For the credentials used by default to access these services, see [Login Credentials](#creds)


## <a name="creds"></a>Login Credentials
For testing/educational purposes, you may want to login to the various services exposed on local ports. For ease of building and use, the CIAO build process has a set of default login credentials it uses (you may change these, but you risk causing some part of the CDN being unable to access resources it needs). Here are the various services and their login credentials:


Service                | Username         | Password
-----------------------|------------------|-----------
Traffic Ops (admin)    | `admin`          | `!!twelve`
RIAK                   | `riakuser`       | `riakpass`
Traffic Ops database | `traffic_ops`    | `password`


To know where these services are located on the local network, see [Ports and Interfaces](#ports)
