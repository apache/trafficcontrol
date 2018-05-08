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

# Cache Inspector Plugin

The cache inspector plugin allows you to view the contents of the cache in a browser.  Point your browser at `http://<yourcacheiporhostname:yourcacheport>/_cacheinspect` to get started. Access to this endpoint is limited to the IP ranges defined in the `stats` object of the global configuration.

Example output:

```
Jump to:  disk  my-disk-cache-two  


*** Cache "" ***

  * Size of in use cache:      1.5M
  * Cache capacity:            9.5M
  * Number of elements in LRU: 54
  * Objects in cache sorted by Least Recently Used on top, showing only first 100 and last 100:

            #    Code      Size                   Age              FreshFor    HitCount      Key
            0     200       15K             600.922ms            50.541023s           1      GET:http://localhost/15k.bin?
            1     200       16K             627.929ms            50.540856s           1      GET:http://localhost/16k.bin?
            2     200       17K             650.813ms            50.540698s           1      GET:http://localhost/17k.bin?
            3     200       18K             671.732ms            50.540547s           1      GET:http://localhost/18k.bin?
            4     200       19K             691.981ms            50.540373s           1      GET:http://localhost/19k.bin?
            5     200       20K             715.566ms            50.540261s           1      GET:http://localhost/20k.bin?
```


Any of the keys can be clicked to peek at the details of this object in cache, and this will not update the LRU list:

```
Key: GET:http://localhost/34k.bin?hah cache: ""

  > User-Agent: curl/7.54.0
  > Accept: */*
  > Host: localhost

  < Last-Modified: Sun, 01 Apr 2018 19:42:43 GMT
  < Etag: "8800-568cead43b2c0"
  < Accept-Ranges: bytes
  < Content-Length: 34816
  < Content-Type: application/octet-stream
  < Date: Sat, 14 Apr 2018 22:18:35 GMT
  < Server: Apache/2.4.29 (Unix)

  Code:                         200
  OriginCode:                   200
  ProxyURL:                     
  ReqTime:                      2018-04-14 16:18:35.098419 -0600 MDT m=+33.256292902
  ReqRespTime:                  2018-04-14 16:18:35.098843 -0600 MDT m=+33.256716928
  RespRespTime:                 2018-04-14 22:18:35 +0000 GMT
  LastModified:                 2018-04-01 19:42:43 +0000 GMT
```

The following querystrings can be used: 

- `cache=<cachename>`
Only display the contents of the cache `<cachename>`. Use `cache=` (empty cachename) for the default memory cachename `""`.
- `head=<number>`
Number of items to list from the top of the top of the LRU. Default is 100.
- `key=<keystring>`
Show details page of the given key.
- `search=<searchstring>`
List only items that have `<searchstring>` in the key. This overrules `<head>` and `<tail>`, when search is used these are ignored. 
- `tail=<number>`
Number of items to list from the bottom of the top of the LRU. Default is 100.
