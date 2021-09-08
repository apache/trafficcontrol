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
# Independent Delivery Service Changes

## Problem Description
In order to propagate changes made to a Delivery Service, updates must be
queued - often for an entire CDN - and a CDN Snapshot must be taken. Generally
speaking, both of these actions will propagate the entirety of the changes
made to any Delivery Services since the last time they were taken. This means
that in order to make changes for one Tenant's Delivery Service, operators must
ensure that no other unfinished changes are pending in the ATC system. This
presents a barrier to Tenants being able to make modifications to their own
Delivery Services, and specifically to their ability to make those changes live.

## Proposed Change
For "Self-Service", it must be possible to change a single Delivery Service
without applying data changes by other Tenants to other Delivery Services. To
accomplish this, the following changes are suggested:

- Change CDN Snapshots to generate when requested, rather than stored

	This requires all tables used to build Snapshots to use a boolean to track
	whether a row has been deleted and use that instead of actual deletions as
	well as using this fake deletion and an insertion instead of normal
	updates.

- Delivery Service Snapshots should be introduced which, similar to CDN
	Snapshots, will be a frozen set of Delivery Service data that can be made
	as a copy of the current data set at will through the Traffic Ops API.

- Add the notion of "default" Delivery Service server assignments.

	A CDN should have a default Topology that is assigned to new Delivery
	Services upon creation if one isn't specified. Alternatively/Additionally,
	servers should be able to be marked as "default assigned", which will cause
	them to be assigned directly to Delivery Services created within the same
	CDN when no Topology is specified (and the CDN does not have a default
	Topology).

### Traffic Portal Impact
Traffic Portal will need UI controls for Delivery Service Snapshots that can be
more or less exactly the same as its current controls for CDN Snapshots. It will
also need a toggle-able switch for the new server field, as well as
controls/display for a CDN's default Topology.

### Traffic Ops Impact
This proposed change would incur massive changes to Traffic Ops.

#### REST API Impact
From the API's perspective, the only changes would be an addition of a new
endpoint for taking/retrieving Delivery Service Snapshots, some changes to
headers and caching for CDN Snapshots, and adding a new field to Server and CDN
objects.

##### `/cdns` and `/cdns/{{ID}}`
These endpoints have identical payload requirements (POST for the former, PUT
for the latter) and response representations for CDN objects. This should not
change, and the response structure for both should be modified to include a new
property:

```typescript
interface NewCDN extends CDN {
	defaultTopology: string | null;
}
```

This new property should appear on response representations as well as being a
required field in request bodies (where request bodies are required).

##### `/deliveryservice_snapshots`
This endpoint will support two methods, POST and GET, requiring the
`ds-snapshot-write` and `ds-snapshot-read` Permissions, respectively (to be
given to the users with `cdn-snapshot` and `ds-read` on migration,
respectively). A POST request body should consist of an array of Delivery
Service XMLIDs like so:

```json
["demo1", "demo2"]
```

The Delivery Services with these given XMLIDs will then be "snapshotted". If any
of the specified Delivery Services falls outside the requesting user's Tenancy,
then no Snapshots are taken and a `403 Forbidden` response is returned.

A GET request will support all of the same filtering currently supported by the
`/deliveryservices` API endpoint through query string parameters. The response
will be an array of Delivery Service Snapshot objects, which have the following
structure (given as a TypeScript interface):

```typescript
interface DeliveryServiceMatch {
	pattern: string;
	setNumber: number;
	type: string;
}

interface Parameter {
	configFile: string;
	id: number;
	name: string;
	secure: boolean;
	value: string;
}

interface URISigningKey {
	alg: string;
	kid: string;
	kty: string;
	k: string;
}

interface URISigningKeys {
	[issuer: string]: {
		renewal_kid: string;
		keys: Array<URISigningKey>;
	};
}

interface URLKeys {
	[keyNum: `key${number}`]: string;
}

interface DeliveryServiceSnapshot {
	active:                     boolean;
	anonymousBlockingEnabled:   boolean;

	// This field will be set at the moment the Snapshot is taken and will not
	// change even if the underlying Capabilities are deleted - the behavior
	// of ATC components under those conditions is undefined.
	capabilities:               Array<string>

	ccrDnsTtl:                  number | null;

	// These fields will be set at the moment the Snapshot is taken and will not
	// change even if the underlying CDN is renamed or deleted - the
	// behavior of ATC components under those conditions is undefined.
	cdnId:                      number;
	cdnName:                    string;

	checkPath:                  string | null;
	consistentHashRegex:        string | null;
	consistentHashQueryParams:  Array<string>;
	deepCachingType:            "ALWAYS" | "NEVER";
	displayName:                string;
	dnsBypassCname:             string | null;
	dnsBypassIp:                string | null;
	dnsBypassIp6:               string | null;
	dnsBypassTtl:               number | null;
	dscp:                       number;
	edgeHeaderRewrite:          string | null;
	firstHeaderRewrite:         string | null;
	geoLimit:                   0 | 1 | 2;
	geoLimitCountries:          string | null;
	geoLimitRedirectURL:        string | null;
	geoProvider:                0 | 1;
	globalMaxMbps:              number | null;
	globalMaxTps:               number | null;
	httpBypassFqdn:             string | null;
	id:                         number;
	infoUrl:                    string | null;
	initialDispersion:          number | null;
	innerHeaderRewrite:         string | null;
	ipv6RoutingEnabled:         boolean;
	lastHeaderRewrite:          string | null;

	// This will be the time at which the Snapshot was taken, not the actual
	// time at which the Delivery Service herein described was last updated.
	lastUpdated:                Date;

	logsEnabled:                boolean;
	longDesc:                   string;
	matchList:                  DeliveryServiceMatch[];
	maxDnsAnswers:              number | null;
	maxOriginConnections:       number;
	maxRequestHeaderBytes:      number;
	midHeaderRewrite:           string | null;
	missLat:                    number;
	missLong:                   number;
	multiSiteOrigin:            boolean;
	originShield:               string | null;

	// These fields will be set at the moment the Snapshot is taken, and will
	// not change even if the underlying Profile changes in any way, including
	// deletion or renaming. The behavior of ATC components under those
	// conditions is undefined.
	profileDescription:         string;
	profileId:                  number;
	profileName:                string;
	parameters:                 Array<Parameter>;

	protocol:                   0 | 1 | 2 | 3;
	qstringIgnore:              0 | 1 | 2;
	rangeRequestHandling:       0 | 1 | 2;
	regexRemap:                 string | null;
	regionalGeoBlocking:        boolean;
	remapText:                  string | null;
	routingName:                string;

	// This field will be set at the moment the Snapshot is taken and will not
	// change even if one of the identified servers is deleted - the behavior of
	// ATC components under those conditions is undefined.
	servers:                    Array<number> | null;

	// This field will be set at the moment the Snapshot is taken and will not
	// change even if the underlying Service Category is renamed or deleted -
	// the behavior of ATC components under those conditions is undefined.
	serviceCategory:            string | null;

	signed:                     boolean;
	signingAlgorithm:           string;
	sslKeyVersion:              number;

	// These fields will be set at the moment the Snapshot is taken, and will
	// not change even if the underlying Tenant changes - and in particular not
	// even if the original Delivery Service's Tenancy changes. Access based on
	// Tenancy to a Delivery Service Snapshot is granted based on this value,
	// and Tenants that no longer exist are determined to fall within all
	// Tenancies. Note that this means that if the Tenancy tree is reordered
	// after the Snapshot was taken, access may be granted to or lost from
	// unexpected answers.
	tenant:                     string;
	tenantId:                   number;

	// This field will be set at the moment the Snapshot is taken and will not
	// change even if the underlying Topology is renamed or deleted - the
	// behavior of ATC components under those conditions is undefined.
	topology:                   string | null;

	trRequestHeaders:           string | null;
	trResponseHeaders:          string | null;

	// These fields will be set at the moment the Snapshot is taken, and will
	// not change even if the underlying Type changes.
	type:                       string;
	typeId:                     number;

	uriSigningKeys:             URISigningKeys | null;
	urlKeys:                    URLKeys | null;
	xmlId:                      string;
}
```
Note that this structure differs from a Delivery Service in many ways.

##### `/servers`, `/servers/details`, `/servers/{{ID}}`, `/deliveryservices/{{ID}}/servers`, and `/deliveryservices/{{ID}}/servers/eligible`
These endpoints have _nearly_ identical payload requirements (for POST and PUT
requests) and response representations for Server objects. This should not
change, and the response structures for all should be modified to include a new
property:

```typescript
interface NewServer extends Server {
	defaultAssigned: boolean;
}

interface NewServerDetails extends ServerDetails {
	defaultAssigned: boolean;
}
```

This new property should appear on response representations as well as being a
required field in request bodies (where request bodies are required).

##### CDN Snapshot endpoints
Now that Traffic Ops generates Snapshots on request every time instead of
storing them, Traffic Ops may cache these Snapshots but only for very short
periods of time, e.g. 1 second, and should return an `Age` HTTP header. The
generation API endpoint(s) should respect `no-cache` directives from the client.

#### Client Impact
The client will need new methods added for the new endpoint, but should
otherwise be covered by shared structure updates.

#### Database Impact
New tables will need to be added to accommodate the Delivery Service Snapshots
above, and new columns will need to be added to the `cdn` and `server` tables to
accommodate the new data they need to store. `cdn.default_topology` should be a
foreign key into `topology`, and `server.default_assign` should not be allowed
to be `NULL` (should default to `FALSE`).

A migration should also be added to add the new Permissions required for taking
Delivery Service Snapshots to all users with the `cdn-snapshot` Permission.

Most notably, all tables that contain rows used to construct a CDN Snapshot
(which is far and away most of them) will need to be changed to have a boolean
field that tells whether or not the object in question has been deleted, and all
DELETE and UPDATE statements will need to be changed, respectively, to UPDATEs
that toggle that flag and INSERTs that insert new data instead of updating old
data.

### T3C Impact
`t3c` will need to be overhauled to construct configuration from Delivery
Service Snapshots in the new API version (only) instead of from raw Delivery
Service information

### Traffic Monitor Impact
Traffic Monitor's Monitoring Snapshots will need to be overhauled to pull from
Delivery Service Snapshots in the new API version (only) instead of from raw
Delivery Service information, but no actual TM changes should be necessary.

### Traffic Router Impact
Traffic Router will need to be overhauled to pull information from Delivery
Service Snapshots in the new API version (only) instead of from raw Delivery
Service information.

### Traffic Stats Impact
None.

### Traffic Vault Impact
None.

### Documentation Impact
The new endpoints, fields, and data model changes will need to be documented.

### Testing Impact
The API tests will probably be impacted in various unforeseen ways that will
need to be compensated for.

### Performance Impact
This will add massive amounts of data stored in the database, potentially
causing database sizes to multiply in a matter of months. As a result, reads
from the database will become gradually slower and more expensive.

### Security Impact
None.

### Upgrade Impact
None.

### Operations Impact
Operators will be freed to make changes to individual Delivery Services as they
see fit without needing to concern themselves with changes being made to other
Delivery Services, and as a result their operating procedures should likely
change.

### Developer Impact
Developers need to be aware of the changes.

## Alternatives
An alternative to timestamps, is doing away with the "snapshot" concept, and
changing Traffic Ops/Portal to do changes entirely in the user's browser local
storage. Then, when the user clicks an "apply" button, all changes would be sent
in a single HTTP request, to be atomically applied.

This has the advantage of simplifying the data model, as well as the application
to Routers and Caches. The idea of a CDN Snapshots and Update Queuing would go
away, but caches, routers, and monitors would simply get all the latest data.

This idea does not yet account for additions versus subtractions, that some data
must be applied to routers, then caches, and some the reverse.

This proposal's timestamp system has the following advantages over local storage
with bulk updates:

1. Timestamps thoroughly accommodate the snapshot-queue direction problem.
	Whereas the local storage does not, and arguably any system will be an
	afterthought, complex, and error-prone. Timestamps solve the problem clearly
	and simply, with extensive history and debugging support.

Timestamps provide extensive history, which make debugging, finding, and fixing
bad states significantly easier. Whereas the Local Storage solution does not
provide any history or debugging aid. A history could be logged, of bulk updates
sent by each user. However, it would be an afterthought, and only provide data
as to which user performed which change, it would not help to debug how Traffic
Ops, ORT, or the Monitor/Router data got into a bad state.

## Dependencies
None.
