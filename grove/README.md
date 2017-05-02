# grove

HTTP caching proxy, implementing RFC 7234

# Installing

The reverse proxy and caching rules are implemented as a library. This project includes a sample single-file application which can be run via `go run`.

1. Install and set up a Golang development environment.
    * See https://golang.org/doc/install
2. Clone this repository into your GOPATH.
```bash
mkdir -p $GOPATH/src/github.com/apache/incubator-trafficcontrol
cd $GOPATH/src/github.com/apache/incubator-trafficcontrol
git clone https://github.com/apache/incubator-trafficcontrol/grove
```
3. Build the library
```bash
cd $GOPATH/src/github.com/apache/incubator-trafficcontrol/grove
go build
```
4. Build the example application, if desired (it may also be run directly via `go run`; see [Running](#Running)).
```bash
cd $GOPATH/src/github.com/apache/incubator-trafficcontrol/grove/service
go build
```

# Running

1. Add remap rules to `remap.config`. See the sample `remap.config`.
    * These rules are of the form `map source-url destination-url`. For example, `map http://localhost:8080 https://www.example.com`.
    * This must include the scheme (`http://`), and may include path parts (`/foo/bar`).
    * The source domain must be sent to the application. For example, `localhost` will automatically be sent to your application.
    * If the application is running on a nonstandard port (not `80` for HTTP or `443` for HTTPS), the `source-url` must include the port.
2. Configure the application, via the config file. See `/service/cfg.json`. This is a JSON file with the following keys:
    * rfc_compliant - whether to use strict RFC 7341 compliance, or to ignore things like `no-cache` from the client in order to protect the origin.
    * port - the port to serve on.
    * cache_size_bytes - size of the cache, in bytes. When this size is exceeded, cached objects will be evicted with a least-recently-used algorithm.
    * remap_rules_file - the file with the proxy remap rules (see step 1)
3. Run the application
```bash
cd $GOPATH/src/github.com/apache/incubator-trafficcontrol/grove
go run service/service.go -config ./service/cfg.json
```
4. Verify it's working by making a request to a remapped endpoint.
    * For example, with the sample config, `curl -vs http://localhost:8080/foo/bar` should return a response from `http://example.com/foo/bar`, and the server should log messages regarding the request and cacheability.
