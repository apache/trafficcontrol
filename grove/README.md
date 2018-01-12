# grove

HTTP caching proxy, implementing RFC 7234

# Building

1. Install and set up a Golang development environment.
    * See https://golang.org/doc/install
2. Clone this repository into your GOPATH.
```bash
mkdir -p $GOPATH/src/github.com/apache/incubator-trafficcontrol
cd $GOPATH/src/github.com/apache/incubator-trafficcontrol
git clone https://github.com/apache/incubator-trafficcontrol/grove
```
3. Build the application
```bash
cd $GOPATH/src/github.com/apache/incubator-trafficcontrol/grove
go build
```
5. Install and configure an RPM development environment
   * See https://wiki.centos.org/HowTos/SetupRpmBuildEnvironment
4. Build the RPM
```bash
./build/build_rpm.sh
```

# Configuration

A config file must be passed with the `-cfg` flag on startup. The RPM uses a config file at `/etc/grove/grove.cfg`.

The config file is JSON of the following format:

```json
{ "rfc_compliant": false,
  "port": 8080,
  "cache_size_bytes": 50000,
  "remap_rules_file": "./remap.json",
```

The config file has the following fields:

| Field | Description |
| --- | --- |
| `rfc_compliant` | Whether to strictly adhere to RFC 7234. If false, client requests which can harm a parent, such as `no-cache` are ignored. |
| `port` | The HTTP port to serve on. |
| `https_port` | The HTTPS port to serve on. |
| `cache_size_bytes` | The maximum size of the memory cache, in bytes. This is a soft maximum, and the cache may temporarily exceed this size until older values can be purged. The cache uses a Least Recently Used algorithm, purging the oldest requested object when a request for an uncached object is received with a full cache. Also note the cache size calculation does not currently count headers. |
| `remap_rules_file` | The file with remap rules. See [Remap Rules](#remap-rules). |
| `concurrent_rule_requests` | The maximum number of simultaneous requests which will be issued to a parent for any rule. |
| `cert_file` | The global HTTPS certificate file to use, for HTTPS remap rules without certificates specified. |
| `key_file` | The global HTTPS certificate key file to use, for HTTPS remap rules without certificates specified. |
| `interface_name` | The name of the network interface to gather statistics for. This does _not_ affect which addresses are bound for listening, currently the app listens on the given port for all addresses, irrespective of interface. |
| `connection_close` | Whether to send a `Connection: Close` header with responses. This is primarily designed for debugging and operations use, for example, to help remove clients from a cache in order to take it out of service. |
| `log_location_error` | The location to log error messages to. May be any file, `stdout`, `stderr`, or `null`. |
| `log_location_warning` | The location to log warning messages to. May be any file, `stdout`, `stderr`, or `null`. |
| `log_location_info` | The location to log informational messages to. May be any file, `stdout`, `stderr`, or `null`. |
| `log_location_debug` | The location to log debug messages to. May be any file, `stdout`, `stderr`, or `null`. |
| `log_location_event` | The location to log access events to. May be any file, `stdout`, `stderr`, or `null`. |
| `parent_request_timeout_ms` | The timeout in milliseconds for requests to parents. |
| `parent_request_keep_alive_ms` | The length of time in milliseconds to keep connections to parents alive, for multiple requests. |
| `parent_request_max_idle_connections` | The maximum number of idle kept-alive connections to retain per parent. |
| `parent_request_idle_connection_timeout_ms` | The length of time in milliseconds to keep an idle parent connection alive, before terminating it. |
| `server_idle_timeout_ms` | The length of time in milliseconds to allow a kept-alive client connection to remain idle, before terminating it. |
| `server_read_timeout_ms` | The length of time in milliseconds to allow a client to read data, before the connection is terminated. This value should be carefully considered, as too short a timeout will result in terminating legitimate clients with slow connections, while too long a timeout will make the server vulnerable to SlowLoris attacks.  |
| `server_write_timeout_ms` | The length of time in milliseconds to allow a client to write data, before the connection is terminated. This value should be carefully considered, as too short a timeout will result in terminating legitimate clients with slow connections, while too long a timeout will make the server vulnerable to SlowLoris attacks.|


# Remap Rules

The remap rules file is specified in the [config file](#configuration).

Note there exists a tool for generating remap rules from [Traffic Control](https://github.com/apache/incubator-trafficcontrol), available [here](https://github.com/apache/incubator-trafficcontrol/grove/tree/master/grovetccfg).

The remap rules file is JSON of the following form:

```json
{
    "parent_selection": "consistent-hash",
    "retry_codes": [ 501, 404 ],
    "retry_num": null,
    "rules": [
        {
            "allow": [ "::1/128", "0.0.0.0/0" ],
            "certificate-file": "",
            "certificate-key-file": "",
            "concurrent_rule_requests": 0,
            "connection-close": false,
            "deny": [ "::1/128", "0.0.0.0/0" ],
            "from": "http://foo.example.net",
            "name": "foo.example.com.http.http",
            "parent_selection": "consistent-hash",
            "query-string": { "cache": true, "remap": true },
            "retry_codes": [ 404, 500 ],
            "retry_num": 5,
            "timeout_ms": 5000,
            "to": [
                {
                    "parent_selection": "consistent-hash",
                    "proxy_url": "http://proxy.example.net:80",
                    "retry_codes": [ 500, 404 ],
                    "retry_num": 5,
                    "timeout_ms": 5000,
                    "url": "http://bar.example.net",
                    "weight": 1
                }
            ],
            "to_client_headers": {
                "drop": [],
                "set": [ { "name": "Accept-Ranges", "value": "None" } ]
            },
            "to_origin_headers": {
                "drop": [ "Range" ],
                "set": []
            }
        }
    ],
    "timeout_ms": 5000
}
```

Rule configuration may be specified at the global, rule, or `to` level, and the most specific field applies. Remap rules have the following configuration fields:

| Field | Description |
| --- | --- |
| `retry_num` | The number of times to retry a parent request. |
| `retry_codes` | The HTTP codes which will be considered failures and cause a failure and cause a retry on the next parent. If `retry_num` tries are exceeded, the final failure response will be cached and returned to the client. |
| `timeout_ms` | The request timeout in milliseconds for the given parent. |
| `parent_selection` | The parent selection algorithm. Currently, only `consistent-hash` is supported. |
| `concurrent_rule_requests` | The maximum number of concurrent requests to make to the parent, for this rule. |
| `allow` | An array of CIDR networks to allow access. This may include both IPv4 and IPv6 networks. Note single IPs must be in CIDR format, e.g. `192.0.2.1/32`. |
| `deny` | An array of CIDR networks to deny access to. This may include both IPv4 and IPv6 networks. Note single IPs must be in CIDR format, e.g. `192.0.2.1/32`. |

The global object must also include a `rules` key, with an array of rule objects. Each remap rule has the following fields:

| Field | Description |
| --- | --- |
| `name` | The internal name for the given rule. This is not used in request mapping, and may be any unique string. |
| `from` | The request to remap, including the scheme and fully qualified domain name. This may also optionally include URL path parts. |
| `certificate-file` | The file path for the certificate for this HTTPS request. This field is not used for HTTP requests. |
| `certificate-key-file` | The file path for the certificate key for this HTTPS request. This field is not used for HTTP requests. |
| `connection-close` | Whether to add a `Connection: Close` header to client responses for this rule. This is designed for maintenance, operations, or debugging. |
| `query-string` | A JSON object with the boolean keys `remap` and `cache`. The `remap` key indicates whether to append request query strings to the parent request. The `cache` key incidates whether to cache requests with different query strings separately. |
| `to` | The array of parents for the given rule. |
|`to_client_headers`| JSON object of header manipulation rules for the response to the client. `set` is an array of headers to set, and `drop` is an array of headers to remove. |
|`to_origin_headers`| JSON object of header manipulation rules for the request to the origin server. `set` is an array of headers to set, and `drop` is an array of headers to remove. |

The objects in the `to` array of parents have the following fields:

| Field | Description |
| --- | --- |
| `url` | The parent URL to remap to, including the scheme and fully qualified domain name. This may also optionally include URL path parts. |
| `weight` | The weight of this parent in the parent selection algorithm. |
| `proxy_url` | The proxy URL, if this parent is being used as a forward proxy. Must include the scheme, fully qualified domain name, and port. If this rule is omitted, the parent will be requested directly with the `url` as a reverse proxy. |

# Running

The application may be run manually via `./grove -cfg grove.cfg`, or if installed via the RPM, as a service via `service grove start` or `systemctl start grove`.

If there are errors, they will be logged to the error location in the config file (`/etc/grove/grove.cfg` for the service), or if the errors are with the config file itself, to stdout.
