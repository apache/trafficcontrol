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
# fakeOrigin

fakeOrigin is a simple HLS video server, capable of simulating live video traffic. It can:

* Serve HLS Live video, by transcoding a static video file as VOD, then manipulating the manifest on the fly to serve an infinitely looping live manifest
* Serve HLS VOD video, from a static video file
* Serve static video and other files

# How to build an rpm
From the root of the trafficcontrol repository use the `pkg` script such as:
```
./pkg -o fakeorigin_build
```
Since this is considered an extra non-required component of the overall functionality of Apache Traffic Control, it's placed in the optional grouping not built by default by `pkg`.

# How to install locally
Local build pre-requesites:
* Go 1.9+
  * OSX: ```brew install go```
  * CentOS: [Instructions]( https://www.itzgeek.com/how-tos/linux/centos-how-tos/install-go-1-7-ubuntu-16-04-14-04-centos-7-fedora-24.html)
* FFMPEG 3.4+ (Optional)
  * OSX: ```brew install ffmpeg --with-rtmp-dump```
  * CentOS: [Instructions](https://linuxadmin.io/install-ffmpeg-on-centos-7/)

and/or just a modern version of Docker & docker compose

If you're building locally, just run ```go install github.com/apache/trafficcontrol/v8/test/fakeOrigin@latest```

If you're just using docker, clone this repository.

# How to use
Running locally:
```
Usage:
fakeOrigin (generates a minimal config.json next to binary)
fakeOrigin -cfg config.json (same as above, but specify the location)
```
Running in docker:
```
docker compose build --no-cache
docker compose up --force-recreate
... customize the config.json created in ./docker_host (maps to /host inside the container, it's really important to customize this appropriately)
docker compose up --force-recreate
```

On startup it will print any routes that are available after transcoding.  You should just be able to plug those m3u8 url into VLC to start streaming.

I'd *highly* recommend going and reading about the [Configuration](./docs/Configuration.md) to learn about what fakeOrigin can do.

There is also another [set of instructions](build/README.md) if you're interested in building your own RPMs and binaries.

# Features
* Transcoding on startup only if the source file or transcoder configuration changes
* Single & Multiple Static file serving support
* HTTP & HTTPS support
* RFC7232 support for CDN caching
* RFC7233 support for Range requests
  * Both single and multi-part ranges are supported
* Arbitrary header response controls
  * Controlled via config file or by client request headers
* Optional in-memory caching
* Support for arbitrary external commands to perform transcoding
  * Supports vod, live, and event m3u8 HLS manifest types
* Support for testing various types of transaction elements using generated output when setting a type of "testing". See the [dtp](./docs/Testing.md) documentation
