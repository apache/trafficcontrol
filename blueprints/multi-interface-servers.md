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
# Additional Server Interfaces

## Problem Description
Multiple network interfaces may be used by cache servers to service requests.
However, these cannot be expressed as parts of the current "server" object.
Specifically, assuming each interface has its own IP address, then each would
need to be treated as a separate cache server. This means that getting clients
to use one when the other is unavailable requires manual intervention.

## Proposed Change
Instead, cache servers should be configurable with numerous network interfaces
monitored by Traffic Monitor that are all associated with the server, each
having its own associated network information e.g. IP Address, but Traffic
Router should only be aware of the "main" interface on which the cache server
listens for connections.

The specific use-case outlined above, obviously, entails some management be done
on the cache servers to actually support multiple interfaces answering to the
same "main" IP address and this should be considered outside of
<abbr title="Apache Traffic Control">ATC</abbr>'s scope of operation.

## Data Model Impact
The proposed structure of each of these interfaces is given below as a pair of
TypeScript<sup>[1](#typescript)</sup> interfaces:
```typescript
interface IPAddress {
	/**
	 * The actual IP (v4 or v6) address which is used by an interface.
	 * Any "mask" should be present here as a CIDR-notation suffix.
	 */
	address: string;

	/**
	 * The IP (v4 or v6) address of the gateway used by this IP address.
	 */
	gateway: string;

	/**
	 * Tells whether or not this address of this interface is the server's
	 * "service" address.
	 * At least one address of EXACTLY ONE interface MUST have this set to
	 * 'true' for a server.
	 */
	serviceAddress: boolean;
}

interface Interface {
	/**
	 * The name of the interface device on the server e.g. eth0.
	 */
	name: string;
	/**
	 * These will be all of the IPv4/IPv6 addresses assigned to the interface,
	 * including gateways and "masks".
	 * It is illegal for an interface to not have at least one associated IP
	 * address.
	 */
	ipAddresses: Array<IPAddress> & {0: IPAddress};
	/**
	 * Whether or not Traffic Monitor should monitor this particular interface.
	 */
	monitor: boolean;
	/**
	 * The interface's Maximum Transmission Unit.
	 * If this is 'null' it is assumed that the interface's MTU is not known/is
	 * irrelevant.
	 */
	mtu: 1500 | 9000 | null;
}
```
We don't have a real data model for ATC, so what follows is an approximately
defined representation of a "server" that is valid within certain contexts. The
modifications made to it for the purposes of the changes herein proposed should
be considered valid in those same contexts - whatever they may be.

```typescript
interface Server {

	cachegroup?:  string | null;
	cachegroupId: number;
	cdnId:        number;
	cdnName?:     string | null;
	domainName:   string;
	guid?:        string | null;
	hostName:     string;
	httpsPort?:   number | null;
	id?:          number;

	// ILO things aren't being moved in because they also require a
	// username/password that aren't used by any other interface type.
	iloIpAddress?: string | null;
	iloIpGateway?: string | null;
	iloIpNetmask?: string | null;
	iloPassword?:  string | null;
	iloUsername?:  string | null;

	/**
	 * New field containing additional interfaces.
	 * It is illegal for a server to not have at least one interface.
	 */
	interfaces: Array<Interface> & {0: Interface};
	lastUpdated?:  string;

	// notice not here: mgmtIpAddress, mgmtIpGateway, mgmtIpNetmask
	// these will be moved into the interfaces property

	offlineReason:   string;
	physLocation?:   string | null;
	physLocationId:  number;
	profile?:        string | null;
	profileDesc?:    string | null;
	profileId:       number;
	revalPending?:   boolean;
	rack?:           string | null;
	routerHostName?: string | null;
	routerPortName?: string | null;
	status?:         string | null;
	statusId:        number;
	tcpPort:         number;
	type?:           string | null;
	typeId:          number;
	updPending?:     boolean;
	xmppId?:         string | null;
	xmppPassword?:   string | null;
}
```

`/monitoring` payloads and CDN Snapshots are purposely un-versioned and
un-modeled, so this section will make no attempts to define the model changes to
those objects.

## Impacted Components
The only components that should need to be modified to accommodate this feature
are Traffic Ops, Traffic Portal, and Traffic Monitor.

### Traffic Ops Impact

#### Database
A new table will need to be introduced to contain interface information:
`interface`, which is described below.
```text
              Table "public.interface"
 Column  |  Type   | Collation | Nullable | Default
---------+---------+-----------+----------+---------
 name    | text    |           | not null |
 server  | bigint  |           | not null |
 monitor | boolean |           | not null |
 mtu     | integer |           |          | 1500
Indexes:
    "interface_pkey" PRIMARY KEY, btree (name, server)
Foreign-key constraints:
    "interface_server_fkey" FOREIGN KEY (server) REFERENCES server(id)
```
Another new table will need to be created to contain address information, which
is described below.
```text
                 Table "public.ip_address"
     Column     |  Type   | Collation | Nullable | Default
----------------+---------+-----------+----------+---------
 address        | inet    |           | not null |
 gateway        | inet    |           | not null |
 serviceaddress | boolean |           | not null |
 server         | bigint  |           | not null |
 interface      | text    |           | not null |
Indexes:
    "ip_address_pkey" PRIMARY KEY, btree (address, interface, server)
Check constraints:
    "ip_address_gateway_check" CHECK (masklen(gateway) = 0)
Foreign-key constraints:
    "ip_address_interface_fkey" FOREIGN KEY (interface, server) REFERENCES interface(name, server)
    "ip_address_server_fkey" FOREIGN KEY (server) REFERENCES server(id)
```

This should be sufficient to capture the model changes. In addition to creating
these tables, a migration will need to be written to convert currently stored
server information to utilize these new tables.

#### API
The affected endpoints will be:

##### `/cdns/{{name}}/configs/monitoring`
The current payload includes a `trafficServers` property that contains an array
of cache server representations to be monitored by a Traffic Monitor within the
path-specified CDN. An example of one such array member is given below.
```json
{
	"cacheGroup": "Ciab_Edge",
	"fqdn": "edge.infra.ciab.test",
	"hashId": "edge",
	"hostName": "edge",
	"interfacename": "eth0",
	"ip6": "dead::babe",
	"ip": "192.168.0.0",
	"port": 80,
	"profile": "ATS_EDGE",
	"status": "REPORTED",
	"type": "EDGE"
}
```
This will need to change to reflect the new set of interfaces, which for the
above example with no additional interfaces would look like:
```json
{
	"cacheGroup": "Ciab_Edge",
	"fqdn": "edge.infra.ciab.test",
	"hashId": "edge",
	"hostName": "edge",
	"interfaces": [
		{
			"name": "eth0",
			"ipAddresses": [
				{
					"address": "192.168.0.0/32",
					"serviceAddress": true
				},
				{
					"address": "dead::babe/128",
					"serviceAddress": true
				}
			]
		}
	],
	"port": 80,
	"profile": "ATS_EDGE",
	"status": "REPORTED",
	"type": "EDGE"
}
```
Note that only interfaces that have "monitor" set to `true` should appear at all
in the payloads of this endpoint. Further note that neither gateways nor
<abbr title="Maximum Transmission Unit">MTU</abbr> are reported.

##### `/cdns/{{name}}/snapshot`
This is a semantic change only - payloads remain structurally unchanged. The
reported `ip`, `ip6` and `interfaceName` fields for elements of the
`contentServers` array property should be an arbitrarily chosen IPv4 service
address, an arbitrarily chosen IPv6 service address, and the interface name of
that/those chosen address(es).

##### `/cdns/{{name}}/snapshot/new`
This is a semantic change only - payloads remain structurally unchanged. The
changes to this endpoint are the same as those described for
`/cdns/{{name}}/snapshot` above.

##### `/servers`/`/servers/{{ID}}`
An example of one element of an output array for this endpoint at the time of
this writing is shown below.
```json
{
	"cachegroup": "CDN_in_a_Box_Mid",
	"cachegroupId": 6,
	"cdnId": 2,
	"cdnName": "CDN-in-a-Box",
	"domainName": "infra.ciab.test",
	"guid": null,
	"hostName": "mid",
	"httpsPort": 443,
	"id": 10,
	"iloIpAddress": "",
	"iloIpGateway": "",
	"iloIpNetmask": "",
	"iloPassword": "",
	"iloUsername": "",
	"interfaceMtu": 1500,
	"interfaceName": "eth0",
	"ip6Address": "fc01:9400:1000:8::120",
	"ip6Gateway": "fc01:9400:1000:8::1",
	"ipAddress": "172.16.239.120",
	"ipGateway": "172.16.239.1",
	"ipNetmask": "255.255.255.0",
	"lastUpdated": "2018-12-05 18:45:05+00",
	"mgmtIpAddress": "",
	"mgmtIpGateway": "",
	"mgmtIpNetmask": "",
	"offlineReason": "",
	"physLocation": "Apachecon North America 2018",
	"physLocationId": 1,
	"profile": "ATS_MID_TIER_CACHE",
	"profileDesc": "Mid Cache - Apache Traffic Server",
	"profileId": 10,
	"rack": "",
	"revalPending": false,
	"routerHostName": "",
	"routerPortName": "",
	"status": "REPORTED",
	"statusId": 3,
	"tcpPort": 80,
	"type": "MID",
	"typeId": 12,
	"updPending": false,
	"xmppId": "mid",
	"xmppPasswd": "",
	"ipIsService": true,
	"ip6IsService": true
}
```
After this feature's inclusion, the same element would look like this:
```json
{
	"cachegroup": "CDN_in_a_Box_Mid",
	"cachegroupId": 6,
	"cdnId": 2,
	"cdnName": "CDN-in-a-Box",
	"domainName": "infra.ciab.test",
	"guid": null,
	"hostName": "mid",
	"httpsPort": 443,
	"id": 10,
	"iloIpAddress": "",
	"iloIpGateway": "",
	"iloIpNetmask": "",
	"iloPassword": "",
	"iloUsername": "",
	"interfaces": [
		{
			"ipAddresses": [
				{
					"address": "172.16.239.120/24",
					"gateway": "172.16.239.1",
					"serviceAddress": true
				},
				{
					"address": "fc01:9400:1000:8::120/128",
					"gateway": "fc01:9400:1000:8::1",
					"serviceAddress": true
				}
			],
			"monitor": true,
			"mtu": 1500,
			"name": "eth0"
		}
	],
	"lastUpdated": "2018-12-05 18:45:05+00",
	"offlineReason": "",
	"physLocation": "Apachecon North America 2018",
	"physLocationId": 1,
	"profile": "ATS_MID_TIER_CACHE",
	"profileDesc": "Mid Cache - Apache Traffic Server",
	"profileId": 10,
	"rack": "",
	"revalPending": false,
	"routerHostName": "",
	"routerPortName": "",
	"status": "REPORTED",
	"statusId": 3,
	"tcpPort": 80,
	"type": "MID",
	"typeId": 12,
	"updPending": false,
	"xmppId": "mid",
	"xmppPasswd": "",
}
```
Request payloads for `PUT` and `POST` methods will need to change similarly. A
server object should be accepted (based only on its `interfaces` property) if
it has at least one interface, every interface has at least one IP address,
exactly one interface has at least one IP address that is a service address, and
no two IP addresses of the same family/version are simultaneously marked as
'service addresses'. And, of course, the property's values must be syntactically
and semantically valid.

Below are some examples of invalid `interfaces` property arrays, to which
Traffic Ops MUST respond with a `400 Bad Request` and an appropriate Alert:

Empty array
```json
{"interfaces": []}
```
Null property
```json
{"interfaces": null}
```
Missing property
```json
{}
```
Empty `ipAddresses` sub-property
```json
{
	"interfaces": [
		{
			"ipAddresses": [],
			"monitor": true,
			"mtu": 1500,
			"name": "eth0"
		}
	]
}
```
Null `ipAddresses` sub-property
```json
{
	"interfaces": [
		{
			"ipAddresses": null,
			"monitor": true,
			"mtu": 1500,
			"name": "eth0"
		}
	]
}
```
Missing `ipAddresses` sub-property
```json
{
	"interfaces": [
		{
			"monitor": true,
			"mtu": 1500,
			"name": "eth0"
		}
	]
}
```
No service addresses
```json
{
	"interfaces": [
		{
			"ipAddresses": [
				{
					"address": "0.0.0.1/32",
					"gateway": "10.0.0.1",
					"serviceAddress": false
				}
			],
			"monitor": true,
			"mtu": 1500,
			"name": "eth0"
		}
	]
}
```
Too many interfaces with service addresses
```json
{
	"interfaces": [
		{
			"ipAddresses": [
				{
					"address": "0.0.0.1/32",
					"gateway": "10.0.0.1",
					"serviceAddress": true
				}
			],
			"monitor": true,
			"mtu": 1500,
			"name": "eth0"
		},
		{
			"ipAddresses": [
				{
					"address": "f1d0::f00d/128",
					"gateway": "dead::beef",
					"serviceAddress": true
				}
			],
			"monitor": true,
			"mtu": 1500,
			"name": "wnlp0"
		}
	]
}
```
Too many service addresses of the same family/version
```json
{
	"interfaces": [
		{
			"ipAddresses": [
				{
					"address": "0.0.0.1/32",
					"gateway": "10.0.0.1",
					"serviceAddress": true
				},
				{
					"address": "0.0.0.2/32",
					"gateway": "10.0.0.1",
					"serviceAddress": true
				}
			],
			"monitor": true,
			"mtu": 1500,
			"name": "eth0"
		}
	]
}
```
##### `/servers/details`
This endpoint is subject to the same changes described above for
`/servers`/`/servers/{{ID}}` - though obviously the full output objects will
differ.

#### Client Impact
The structures output by various clients will change, but no client code changes
should be necessary.

### Traffic Portal Impact
<!--
*How* will this impact Traffic Portal?
What new UI changes will be required?
Will entirely new pages/views be necessary?
Will a new field be added to an existing form?
How will the user interact with the new UI changes?
-->

### ORT Impact
<!--
*How* will this impact ORT?
-->

### Traffic Monitor Impact
<!--
*How* will this impact Traffic Monitor?

Will new profile parameters be required?
-->

### Traffic Router Impact
<!--
*How* will this impact Traffic Router?

Will new profile parameters be required?
How will the CRConfig be changed?
How will changes in Traffic Ops data be reflected in the CRConfig?
Will Traffic Router remain backwards-compatible with old CRConfigs?
Will old Traffic Routers remain forwards-compatible with new CRConfigs?
-->

### Documentation Impact
<!--
*How* will this impact the documentation?

What new documentation will be required?
What existing documentation will need to be updated?
-->

### Testing Impact
<!--
*How* will this impact testing?

What is the high-level test plan?
How should this be tested?
Can this be tested within the existing test frameworks?
How should the existing frameworks be enhanced in order to test this properly?
-->

### Performance Impact
<!--
*How* will this impact performance?

Are the changes expected to improve performance in any way?
Is there anything particularly CPU, network, or storage-intensive to be aware of?
What are the known bottlenecks to be aware of that may need to be addressed?
-->

### Security Impact
<!--
*How* will this impact overall security?

Are there any security risks to be aware of?
What privilege level is required for these changes?
Do these changes increase the attack surface (e.g. new untrusted input)?
How will untrusted input be validated?
If these changes are used maliciously or improperly, what could go wrong?
Will these changes adhere to multi-tenancy?
Will data be protected in transit (e.g. via HTTPS or TLS)?
Will these changes require sensitive data that should be encrypted at rest?
Will these changes require handling of any secrets?
Will new SQL queries properly use parameter binding?
-->

### Upgrade Impact
<!--
*How* will this impact the upgrade of an existing system?

Will a database migration be required?
Do the various components need to be upgraded in a specific order?
Will this affect the ability to rollback an upgrade?
Are there any special steps to be followed before an upgrade can be done?
Are there any special steps to be followed during the upgrade?
Are there any special steps to be followed after the upgrade is complete?
-->

### Operations Impact
<!--
*How* will this impact overall operation of the system?

Will the changes make it harder to operate the system?
Will the changes introduce new configuration that will need to be managed?
Can the changes be easily automated?
Do the changes have known limitations or risks that operators should be made aware of?
Will the changes introduce new steps to be followed for existing operations?
-->

### Developer Impact
<!--
*How* will this impact other developers?

Will it make it easier to set up a development environment?
Will it make the code easier to maintain?
What do other developers need to know about these changes?
Are the changes straightforward, or will new developer instructions be necessary?
-->

## Alternatives
<!--
What are some of the alternative solutions for this problem?
What are the pros and cons of each approach?
What design trade-offs were made and why?
-->

## Dependencies
<!--
Are there any significant new dependencies that will be required?
How were the dependencies assessed and chosen?
How will the new dependencies be managed?
Are the dependencies required at build-time, run-time, or both?
-->

## References
<!--
Include any references to external links here.
-->

<a name="typescript">1:</a> The syntax should be mostly self-explanatory, but
see their [official documentation for a five-minute introduction](https://www.typescriptlang.org/docs/handbook/typescript-in-5-minutes.html)
