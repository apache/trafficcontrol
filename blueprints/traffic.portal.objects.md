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

# Traffic Portals as API objects

## Problem Description
Traffic Portals are, today, configured in Traffic Ops as regular servers, with
all the details that entails. A Traffic Portal can have updates and
revalidations pending, has a physical location, is associated with a CDN
despite being primarily for operating arbitrary numbers of CDNs, belongs to a
Cache Group despite not being a Cache Server, has a Type with no semantic
meaning beyond identifying it as a Traffic Portal instance, has a Profile with
no meaningful Parameters, has a HashID despite not being hashed for routing by
Traffic Router, and can have multiple network interfaces with their own
bandwidth limits - despite that Traffic Monitor does not monitor Traffic Portals
at all.

## Proposed Change
All of these properties are superfluous with no real meaning for Traffic
Portals, and for that reason alone they'd be better off as their own object.
The new Traffic Portal object would be compact, containing only the properties
that are needed by and make sense on Traffic Portal instances. Doing this would
also be a step toward avoiding complex filtering on servers for processes that
are looking for specific kinds of servers (though as far as this writer knows,
there aren't any that look specifically for Traffic Portal servers).

## Data Model Impact
The proposed new object type is expressed below as a TypeScript interface.

```typescript
interface TrafficPortal {
	ipv4Address: string | null;
	ipv6Address: string | null;
	notes: string;
	online: boolean;
	tags: Set<string>;
	url: URL;
}
```

Each property is further described in the following subsections.

### IPv4 Address
This string contains a valid IPv4 Address at which the UI is served. If it is
`null`, then it is assumed that there is no static IPv4 Address at which
the Traffic Portal instance may be reached - or it does not communicate over
IPv4 at all.

### IPv6 Address
This string contains a valid IPv6 Address at which the UI is served. If it is
`null`, then it is assumed that there is no static IPv6 Address at which
the Traffic Portal instance may be reached - or it does not communicate over
IPv6 at all.

### Notes
This string contains arbitrary text for containing miscellaneous information.

### Online
This boolean value indicates whether or not this Traffic Portal is online and
serving its user interface.

### Tags
The Tags associated with a Traffic Portal are represented by a set of strings
that are Tag Names.

(Note that this property should only be added if the implementation is done
after #4819 is merged - if ever.)

### URL
This string is a full URL at which the UI is served, including scheme (e.g.
`http://`) and optionally port (e.g. `80`), e.g.
`https://trafficportal.infra.ciab.test:443`.

If the port is omitted, it will be guessed based on the protocol indicated by
the schema, e.g. `80` for `http://`, which indicates the HTTP protocol. The only
valid schemes are `https://` and `http://`.

This property uniquely identifies a Traffic Portal, and must therefore
obviously be unique (this allows multiple Traffic Portals - or other
services - to be running at the same static IP address(es) as long as they
are not literally the same network location; port, host, or protocol could
differ).

## Impacted Components
Since Traffic Portals are intended only as a UI, and need not even exist for an
ATC system to function normally, the impact should be quite limited.

### Traffic Ops Impact

#### Database Impact
A new table must be created to store Traffic Portal information. It can be
populated trivially from the information of servers with Type IDs that identify
Types which have the exact Name "TRAFFIC_PORTAL". Removal of existing data in
the generic `servers` table should be left up to operators' discretion. The
table should look similar to the one shown below:

```text
                       Table "public.traffic_portal"
    Column    |           Type           | Collation | Nullable | Default
--------------+--------------------------+-----------+----------+----------
 ipv4_address | inet                     |           |          |
 ipv6_address | inet                     |           |          |
 last_updated | timestamp with time zone |           | not null | now()
 notes        | text                     |           | not null | ''::text
 online       | boolean                  |           | not null | false
 url          | text                     |           | not null |
Indexes:
    "traffic_portal_pkey" PRIMARY KEY, btree (url)
Check constraints:
    "traffic_portal_url_check" CHECK (lower(url) ~~ 'http://%'::text OR lower(url) ~~ 'https://%'::text)
```

Furthermore, if Tags are implemented as a property of Traffic Portals (subject
to the aforementioned conditions), a "join table" should exist between the two,
which should be extremely similar to the one shown below:

```text
       Table "public.traffic_portal_tags"
 Column | Type | Collation | Nullable | Default
--------+------+-----------+----------+---------
 tp     | text |           | not null |
 tag    | text |           | not null |
Indexes:
    "traffic_portal_tags_pkey" PRIMARY KEY, btree (tp, tag)
Foreign-key constraints:
    "traffic_portal_tags_tag_fkey" FOREIGN KEY (tag) REFERENCES tag(name)
    "traffic_portal_tags_tp_fkey" FOREIGN KEY (tp) REFERENCES traffic_portal(url)
```

#### API Impact
The following new endpoints will need to be created (examples and specs assume
that Tag support is added).

##### `/traffic_portals`
This endpoint deals with manipulations and representations of Traffic Portal
objects configured in Traffic Ops.

###### GET
Retrieves Traffic Portal representations.

- *Required Permissions* `traffic-portals-read`
- *Response Type* Array

**Request Structure**
This method of this endpoint implements the Age Filtering and Sorting and
Pagination query parameters as outlined in [the relevant official documentation
sections](https://traffic-control-cdn.readthedocs.io/en/latest/development/api_guidelines.html#age-filtering).

It further provides the following query parameters:

*GET `/traffic_portals` Query Parameters*

Parameter | Description
==========+============
ipv4Address | Filters results to contain only the Traffic Portals with the given IPv4 Address.
ipv6Address | Filters results to contain only the Traffic Portals with the given IPv6 Address.
online | Filters results to contain only the Traffic Portals that have the provided "Online" value - valid values are "true" and "false".
tag | Filters results to contain only the Traffic Portals with the given Tag. Multiple tags may be given, either in regular, `application/x-www-form-urlencoded` format like e.g. `tag=Foo&tag=Bar` *or* as a comma-separated list, e.g. `tag=Foo,Bar`. In either case - or even in the case of mixed conventions - the filtered results will contain only those Traffic Portals that have *all* specified Tags.
url | Filters results to contain only the Traffic Portals with the given URL. This URL must be valid, and may not contain a non-root path - e.g. `https://example.com:443/` is acceptable, whereas `https://example.com:443/foo` is not.

*Request Example*
```http
GET /api/4.0/traffic_portals?online=true HTTP/1.1
Host: trafficops.infra.ciab.test
Accept: application/json, */*;q=0.9
Cookie: mojolicious=...

```

**Response Structure**
The response is a set of representations of Traffic Portal objects, each
representation extended with the `lastUpdated` property containing the
Date/Time at which the Traffic Portal object was last modified.\\
This method of this endpoint also implements the `count` property of the
top-level `summary` object as described in [the section on `summary` in the
official documentation](https://traffic-control-cdn.readthedocs.io/en/latest/api/index.html#summary).

*Response Example*
```http
HTTP/1.1 200 OK
Content-Type: application/json
Server: Traffic Ops/5.0
Date: Wed, 18 Nov 2020 20:46:57 GMT
Content-Length: 237
ETag: Sample Text

{ "response": [
	{
		"ipv4Address": "192.168.240.12",
		"ipv6Address": null,
		"lastUpdated": "2009-11-10T23:00:00Z"
		"notes": "The default Traffic Portal instance for CDN-in-a-Box.",
		"online": true,
		"tags": [],
		"url": "https://trafficportal.infra.ciab.test:443/"
	}
],
"summary": {
	"count": 1
}}
```

###### POST
Creates new Traffic Portal objects.

- *Required Permissions* `traffic-portals-write`
- *Response Type* Object

**Request Structure**
The body of a POST request to `/traffic_portals` is a representation of a
Traffic Portal object to be created. Only one query parameter is supported, as
shown below.

*POST `/traffic_portals` Query Parameters*

Parameter | Description
==========+=============
lookupAddress | When this query string parameter is given, and is "true", then when the request body contains a `null` `ipv4Address` or `ipv6Address`, the value will be "filled" in by looking up the host name portion of the request body's `url` property. This may also be one of the strings "ipv4Address" and "ipv6Address", to limit the "filling in" behavior to the named IP address property.

*Request Example*
```http
POST /api/4.0/traffic_portals?lookupAddress=ipv4Address HTTP/1.1
Host: trafficops.infra.ciab.test
Accept: application/json, */*;q=0.9
Cookie: mojolicious=...
Content-Type: application/json
Content-Length: 197

{
	"ipv4Address": null,
	"ipv6Address": null,
	"notes": "The default Traffic Portal instance for CDN-in-a-Box.",
	"online": true,
	"tags": [],
	"url": "https://trafficportal.infra.ciab.test:443/"
}
```

**Response Structure**
The response is a representation of the created Traffic Portal object, augmented
with the `lastUpdated` property containing the Date/Time at which the
Traffic Portal object was last modified (which should be approximately equal to
the current time).

*Response Example*
```http
HTTP/1.1 201 Created
Content-Type: application/json
Server: Traffic Ops/5.0
Date: Wed, 18 Nov 2020 20:46:57 GMT
Content-Length: 237
ETag: Sample Text
Location: /traffic_portals/https%3A%2F%2Ftrafficportal.infra.ciab.test%3A443

{ "response": {
	"ipv4Address": "192.168.240.12",
	"ipv6Address": null,
	"lastUpdated": "2009-11-10T23:00:00Z",
	"notes": "The default Traffic Portal instance for CDN-in-a-Box.",
	"online": true,
	"tags": [],
	"url": "https://trafficportal.infra.ciab.test:443/"
},
"alerts": [
	{
		"level": "success",
		"text": "Traffic Portal 'https://trafficportal.infra.ciab.test:443' was created."
	}
]}
```

##### `/traffic_portals/{{URL}}`
This endpoint deals with manipulations and representations of a single Traffic
Portal object identified by the `URL` in the request path.

###### GET
Retrieves a Traffic Portal representation.
- *Required Permissions* `traffic-portals-read`
- *Response Type* Object

**Request Structure**
This method of this endpoint provides no query parameters.

*Request Example*
```http
GET /api/4.0/traffic_portals/https%3A%2F%2Ftrafficportal.infra.ciab.test%3A443 HTTP/1.1
Host: trafficops.infra.ciab.test
Accept: application/json, */*;q=0.9
Cookie: mojolicious=...

```

**Response Structure**
The response is a representation of the requested Traffic Portal object,
extended with the `lastUpdated` property containing the Date/Time at which
the Traffic Portal object was last modified.

*Response Example*
```http
HTTP/1.1 200 OK
Content-Type: application/json
Server: Traffic Ops/5.0
Date: Wed, 18 Nov 2020 20:46:57 GMT
Content-Length: 237
ETag: Sample Text

{ "response": {
	"ipv4Address": "192.168.240.12",
	"ipv6Address": null,
	"lastUpdated": "2009-11-10T23:00:00Z"
	"notes": "The default Traffic Portal instance for CDN-in-a-Box.",
	"online": true,
	"tags": [],
	"url": "https://trafficportal.infra.ciab.test:443/"
}}
```

###### PUT
Replaces a Traffic Portal object with one provided.

- *Required Permissions* `traffic-portals-write`
- *Response Type* Object

**Request Structure**
The body of a PUT request to `/traffic_portals/{{URL}}` is a
representation of a Traffic Portal object to replace the one identified in the
request path. Only one query parameter is supported, as shown below.

*POST `/traffic_portals} Query Parameters`*

Parameter | Description
==========+=============
lookupAddress | When this query string parameter is given, and is "true", then when the request body contains a `null` `ipv4Address` or `ipv6Address`, the value will be "filled" in by looking up the host name portion of the request body's `url` property. This may also be one of the strings "ipv4Address" and "ipv6Address", to limit the "filling in" behavior to the named IP address property.

*Request Example*
```http
PUT /api/4.0/traffic_portals/https%3A%2F%2Ftrafficportal.infra.ciab.test%3A443?lookupAddress=ipv4Address HTTP/1.1
Host: trafficops.infra.ciab.test
Accept: application/json, */*;q=0.9
Cookie: mojolicious=...
Content-Type: application/json
Content-Length: 197
If-Unmodified-Since: Wed, 18 Nov 2020 20:46:57 GMT

{
	"ipv4Address": null,
	"ipv6Address": "::1",
	"notes": "The default Traffic Portal instance for CDN-in-a-Box.",
	"online": true,
	"tags": [],
	"url": "https://trafficportal.infra.ciab.test:443/"
}
```

**Response Structure**
The response is a representation of the requested Traffic Portal object, after
modifications and augmented with the `lastUpdated` property containing the
Date/Time at which the Traffic Portal object was last modified (which should be
approximately equal to the current time).

*Response Example*
```http
HTTP/1.1 200 OK
Content-Type: application/json
Server: Traffic Ops/5.0
Date: Wed, 18 Nov 2020 20:46:57 GMT
Content-Length: 237
ETag: Sample Text

{ "response": {
	"ipv4Address": "192.168.240.12",
	"ipv6Address": "::1",
	"lastUpdated": "2020-18-11T20:46:57.1Z",
	"notes": "The default Traffic Portal instance for CDN-in-a-Box.",
	"online": true,
	"tags": [],
	"url": "https://trafficportal.infra.ciab.test:443/"
},
"alerts": [
	{
		"level": "success",
		"text": "Traffic Portal https://trafficportal.infra.ciab.test:443 was updated."
	}
]}
```

###### PATCH
Modifies the identified Traffic Portal with the partial representation provided.

- *Required Permissions* `traffic-portals-write`
- *Response Type* Object

**Request Structure**
The body of a PATCH request to `/traffic\_portals/{{URL}}` is a partial
representation of a Traffic Portal object to overwrite the properties of the one
identified in the request path with those provided in the body.

*Request Example*
```http
PATCH /api/4.0/traffic_portals/https%3A%2F%2Ftrafficportal.infra.ciab.test%3A443 HTTP/1.1
Host: trafficops.infra.ciab.test
Accept: application/json, */*;q=0.9
Cookie: mojolicious=...
Content-Type: application/json
Content-Length: 24
If-None-Match: Different Sample Text

{
	"ipv6Address": null
}
```

**Response Structure**
The response is a representation of the requested Traffic Portal object, after
modifications and augmented with the `lastUpdated` property containing the
Date/Time at which the Traffic Portal object was last modified (which should be
approximately equal to the current time).

*Response Example*
```http
HTTP/1.1 200 OK
Content-Type: application/json
Server: Traffic Ops/5.0
Date: Wed, 18 Nov 2020 20:46:57 GMT
Content-Length: 237
ETag: Sample Text

{ "response": {
	"ipv4Address": "192.168.240.12",
	"ipv6Address": null,
	"lastUpdated": "2020-18-11T20:46:57.2Z",
	"notes": "The default Traffic Portal instance for CDN-in-a-Box.",
	"online": true,
	"tags": [],
	"url": "https://trafficportal.infra.ciab.test:443/"
},
"alerts": [
	{
		"level": "success",
		"text": "Traffic Portal https://trafficportal.infra.ciab.test:443 was updated."
	}
]}
```

###### DELETE
Deletes the identified Traffic Portal.

- *Required Permissions* `traffic-portals-write`
- *Response Type* Object

**Request Structure**
DELETE requests to `/traffic_portals/{{URL}}` may have no body, nor
does this method of this endpoint provide any query parameters.

*Request Example*
```http
DELETE /api/4.0/traffic_portals/https%3A%2F%2Ftrafficportal.infra.ciab.test%3A443 HTTP/1.1
Host: trafficops.infra.ciab.test
Accept: application/json, */*;q=0.9
Cookie: mojolicious=...
If-Unmodified-Since: Wed, 18 Nov 2020 20:46:57 GMT

```

**Response Structure**
The response is a representation of the deleted Traffic Portal object, augmented
with the `lastUpdated` property containing the Date/Time at which the
Traffic Portal object was last modified (which should be approximately equal to
the current time).

*Response Example*
```http
HTTP/1.1 200 OK
Content-Type: application/json
Server: Traffic Ops/5.0
Date: Wed, 18 Nov 2020 20:46:57 GMT
Content-Length: 237
ETag: Sample Text

{ "response": {
	"ipv4Address": "192.168.240.12",
	"ipv6Address": null,
	"lastUpdated": "2020-18-11T20:46:57.3Z",
	"notes": "The default Traffic Portal instance for CDN-in-a-Box.",
	"online": true,
	"tags": [],
	"url": "https://trafficportal.infra.ciab.test:443/"
},
"alerts": [
	{
		"level": "success",
		"text": "Traffic Portal https://trafficportal.infra.ciab.test:443 was deleted."
	}
]}
```
