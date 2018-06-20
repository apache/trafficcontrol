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

# Range Request Handler plugin

The `range_req_handler` plugin handles range requests in Grove. It has 3 different modes, that are described below. To enable the plugin, add `range_req_handler` to the `plugins` array in the `grove.cfg` file, and add

```
"plugins": {
    "range_req_handler": {
        "mode": "get_full_serve_ranges",
        "cache-key-string":"range_rq_handler_ck_string"
    }
}
```

to your remap rule. The `cache-key-string` configuration is optional, when omitted, this defaults to `grove_range_req_handler_plugin_data`. Mode can be one of the three options described below.

#  1 get_full_serve_ranges
In this mode, the plugin will remove the `Range` header from the request to the origin on a MISS, causing a "normal" GET to the origin and for Grove to store the whole object in cache. When responding to the client, the plugin will build the appropriate `206 Partial Content` from the cached full object. This mode supports multi-part and will create a multi-part content type response if needed.

This mode is appropriate for MPEG-DASH video where the full objects are not too larget to GET without using range requests, and the range requests the clients use are not all on the same boundaries.

# 2 store_ranges
In this mode the plugin will add the range to the cache-key, and store them as is. When all clients are guaranteed to use the same boundaries for their range requests, and the full object will never be requested this mode can be used.

# 3 slice

In slice mode, the plugin will break the requested object up in to slices of a configurable size, and slices are requested from the origin and stored in the cache as they are needed to serve client requests. To configure this mode, add this to the plugins of your remap:

```
"plugins": {
    "range_req_handler": {
        "mode": "slice",
        "slice-size": 1048576,
        "wg-size": 32,
        "cache-key-string":"range_rq_handler_ck_string"
    }
}
```

`slice-size` is the sice of the slices being used in bytes, and `wg-size` is the number of concurrent requests the plugin will do to the origin.

Note that in the current implementation, the slice requests are being requested "through" the grove cache, meaning they will show up in the grove access logs as client requests from localhost.
