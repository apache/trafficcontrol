<!--
    Licensed to the Apache Software Foundation (ASF) under one
    or more contributor license agreements.  See the NOTICE file
    distributed with this work for additional information
    regarding copyright ownership.  The ASF licenses this file
    to you under the Apache License, Version 2.0 (the
    "License"); you may not use this file except in compliance
    with the License.  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing,
    software distributed under the License is distributed on an
    "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
    KIND, either express or implied.  See the License for the
    specific language governing permissions and limitations
    under the License.
-->
# Configuration
The `config.json` file has basically 2 sections
1. Server Info
2. Endpoints

**For Developers**

The underlying data structures for these can be found in [endpoint.go](../endpoint/endpoint.go) and enums can be found in [endpoint_enums.go](../endpoint/endpoint_enums.go).  We may periodically inject new configs into your config.json, but that should be only if they are required with a sane default or inert.

## Server Info
```json
"server": {
  "http_port": 8080,
  "https_port": 8443,
  "ssl_cert": "server.crt",
  "ssl_key": "server.key",
  "binding_address": "127.0.0.1",
  "crossdomain_xml_file": "./example/crossdomain.xml",
  "read_timeout": 0,
  "write_timeout": 0,
}
```
This allows you to set what ports and IP addresses you'd like to listen on.  You can also supply where your SSL certificates and key should be located or created.  They should be in PEM format. The read and write timeouts are optional with a default of 0 that means no timeout, values should be in seconds and will impose a time limit on the read or write side of transactions.

## Endpoints
This is where the meat of your config will be.
```json
"endpoints": [
  {
    "id": "EXAMPLE_ENDPOINT",
    "source": "./example/video/kelloggs.mp4",
    "outputdir": "./out/EXAMPLE_ENDPOINT",
    "type": "live",
    "default_headers": {
      "another-custom-header": [
        "\"foo\", \"bar\""
      ],
      "my-custom-header": [
        "foo",
        "bar"
      ]
    },
    "manual_command": [
      "transcode.cli.executable",
      "-arguments",
      "-anotherArg",
      "%SOURCE%",
      "%OUTPUTDIRECTORY%/%DISKID%.out"
    ]
  }
]
```
A few things you should know:
1. `id` must be unique and will be the first segment of your endpoint URL path
2. `override_disk_id` will default to be id and can usually be ommitted.  This is only for when you need different endpoints that share the same transcoder output
3. Be sure `outputdir` has enough disk space to handle whatever you're asking fakeOrigin to generate.  This should also be unique to each `override_disk_id`.
4. `manual_command` includes some string required token replacements to enable fakeOrigin to parse the output

To help make things a bit easier, you can take a look at [some samples](./Endpoint.Examples.md) on how to build endpoints.
The transcoder phase will be skipped for static and dir `type` as well as other `type` where the a hash of the input file and the `manual_command` are unchanged.

## Manual Command Tokens
You cannot use unescaped spaces in commands.  Each segment of the command is a separate entry in the array.

Since fakeOrigin still needs to know about where certain things are, such as the metadata generation, file organization, and the live manifest interceptor, a set of tokens are defined for use in a manual command.  These are mostly taken from the endpoint properties themselves.
```
%DISKID% :          Identifier for use in file paths, defaults to same as ID
%ENDPOINTTYPE% :    Endpoint Type [live, vod, event, static, dir]
%ID% :              Identifier for use in url paths
%MASTERMANIFEST% :  A primary video manifest path on disk [OutputDirectory + "/" + DiskID + ".m3u8"]
%OUTPUTDIRECTORY% : The directory on disk for output files
%SOURCE% :          The original source file used in transcoding
```
When writing manual commands, you really want to aim for your transcoder to output all files into a single directory (%OUTPUTDIRECTORY%) after consuming your single source asset (%SOURCE%) with a primary manifest (%MASTERMANIFEST%).  In fact, these 3 parameters are required.

## Adaptive Bit Rate
When using HLS with adaptive bitrate layers, the master manifest generator within fakeOrigin expects all layer manifests in the `%OUTPUTDIRECTORY%` with the format `%DISKID%.*(?<width>\d+)x(?<height>\d+)-(?<bandwidth>\d+).m3u8`.  For example, `%DISKID%_1920x1080-40000.m3u8`.

# Client-side Header Overrides
Setting `default_headers` in the config file makes sense for cases like video players who can't supply custom HTTP headers.  However for testing specific use cases, you can also supply custom headers as a client to override what gets returned in the fakeOrigin response.  If you give a header with the Prefix `Fakeorigin-`, that prefix will be stripped and the rest of the client header string will be returned.

Example (Normal):
```bash session
curl -vs4 -o /dev/null http://localhost:8080/SampleVideo/kelloggs.mp4
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /SampleVideo/kelloggs.mp4 HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.59.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Another-Custom-Header: foo, bar
< Content-Type: video/mp4
< My-Custom-Header: foo
< My-Custom-Header: bar
< Date: Thu, 19 Jul 2018 16:08:35 GMT
< Transfer-Encoding: chunked
<
```
Example (With Client Header Override):
```bash session
curl -vs4 -o /dev/null http://localhost:8080/SampleVideo/kelloggs.mp4 -H 'Fakeorigin-A-Custom_header: foo' -H 'fakeOrigin-Another-Custom-Header: "baz"'
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /SampleVideo/kelloggs.mp4 HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.59.0
> Accept: */*
> Fakeorigin-A-Custom_header: foo
> Fakeorigin-Another-Custom-Header: "baz"
>
< HTTP/1.1 200 OK
< A-Custom_header: foo
< Another-Custom-Header: "baz"
< Content-Type: video/mp4
< My-Custom-Header: foo
< My-Custom-Header: bar
< Date: Thu, 19 Jul 2018 16:10:43 GMT
< Transfer-Encoding: chunked
<
```
## Deterministic Testing Protocol
This is a special type of endpoint when using type `testing`. The user can give various arguments in the path of the request on this endpoint in order to generate expected responses from the origin. Headers can be generated or manipulated and various conditions set like an initial stall or delay in response and various other things. See the [dtp](./docs/Testing.md) documentation