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
# Server Capabilities

## Problem Description

Suppose a Traffic Control operator has servers of a particular type. For example, servers with only Memory and no Disk cache. It's possible today to only route to those Edges, via manual Delivery Service Server assignments. But suppose you have a Mid server with only Memory and no Disk. Then suppose you have a certain class of traffic you need to route to this Mid, but not other traffic. For example, Delivery Services with small images; but not DSes with large binary files, which would destroy the cache. Right now, this isn't possible in Traffic Control.

We propose a feature, called "Server Capabilitites" to solve this problem. Servers will have "capabilitites" and Delivery Services will have "Required Capabilitites," and if a server does not have a required capability, then it should not be manually assignable as an Edge, and a Mid must not be added as a parent for that DS in the Edge ATS config.

Initially, this is completely backwards-compatible on upgrade, because initially no Delivery Service will have any Required Capabilities.

This feature will be completely optional. TC operators who don't need Server Capabilities will simply not create them.

Server Capabilities will be arbitrary strings. The ATS project will not "seed" any, and impose as little direction as possible. For example, Server Capabilitites could be Cache types, Server types, Hardware types, ATS versions, Lua support, other ATS plugin support, or any other features an operator needs to route (or not route) on.


## Proposed Change

- Servers have Capabilities (i.e.: CACHE_MEMORY, CACHE_DISK )

- Delivery Services have Required Capabilities (i.e.: CACHE_MEMORY, CACHE_DISK)

- When generating Configuration (e.g. ATS parent.config):
  - If a Mid Server which is otherwise parented to an Edge does not have all Required Capabilities of a Delivery Service, that Mid will not be inserted as a parent for that Delivery Service’s remap and parent rules.
  - Backward Compatibility is automatic:
    - Initially, no Delivery Services have Required Capabilities. Hence, everything behaves like it does today.
    - If an EDGE server does not have all capabilities required by a DS, that server SHOULD NOT be assignable to that DS. 
    - If an EDGE server is assigned to a DS which requires capabilities that server has, it SHOULD NOT be possible to remove those capabilities from that server.

- Initially, Server Capability names are limited to `[a-Z]`, `[0-9]`, `_`, and `-`. This allows them to be put in most parts of a URI, as well as being the characters of Base64 URL Encoding. We may decide to allow high unicode later. But for now, since we're not sure yet what usage will look like, it's much easier to add characters later without breaking people, than it is to remove them.

### Traffic Portal Impact

- Server Page
  - New dropdown to add to a text list: Capabilities
  - List elements have a button to remove the capability
  - Dropdown/List is just one option - other GUI components may be used
  - Component should consider that there could be many Capabilities

- Delivery Service Page
  - New dropdown to add to a text list: Required Capabilities
  - List elements have a button to remove a capability
  - Dropdown/List is one option - other GUI components may be used
  - Component should consider that there could be many Capabilities

- New page for Server Capability Types
  - Ability to add, delete Server Capability Types
  - May be a text input, which adds to a list, with list element buttons to remove from the list


### Traffic Ops Impact

`/server_server_capabilities`
  - List all server capabilities assigned to a server
  - GET+POST+DELETE
  - PUT doesn't make sense. To remove one and add another, DELETE and POST.

`/deliveryservices_required_capabilities`
  - List all DS Required Capabilities
  - GET+POST+DELETE
  - PUT doesn't make sense. To remove one and add another, DELETE and POST.

`/server_capabilities`
  - List, create, and delete Server Capabilities.
  - POST+GET+DELETE
  - PUT doesn’t make sense, Capabilities are strings, if one is misspelled it should be removed and re-added. Simpler, less error-prone than cascading changes.

Additionally, ATS config generation (which is currently in Traffic Ops, but in the process of being moved to its own library/app) will require changes. Primarily `parent.config`, as in the description, to prevent Mids being assigned as Parents for DSes for which they lack the capabilitites.

#### REST API Impact

See Traffic Ops Impact

#### Client Impact

See Traffic Ops Impact - client functions for each endpoint.

#### Data Model / Database Impact

New tables for Capabilitites: `server_capability`, `server_server_capability`, `delivery_service_required_server_capability`.

### ORT Impact

n/a

### Traffic Monitor Impact

n/a

### Traffic Router Impact

n/a

### Traffic Stats Impact

n/a

### Traffic Vault Impact

n/a

### Documentation Impact

The new UI and API endpoints will be documented. See Traffic Ops Impact.

### Testing Impact

The new UI, API endpoints, and ATS config changes, will have TO API tests and unit tests. See Traffic Ops Impact.


### Performance Impact

No performance impact expected.

### Security Impact

n/a

### Upgrade Impact

None. Upgrade will have no Required Capabilities, and config and server assignments will remain unchanged until an operator creates Server Capabilities.

### Operations Impact

None, until Server Capabilities are added. If a TC operator decides to use Server Capabilitites, they will they have to learn how to use them, and create and assign the Server Capabilitites they desire.

### Developer Impact

Code is added to TO endpoints, and config generation. Config generation complexity increases very slightly. TO endpoints are stand-alone and don't affect other endpoints.

## Alternatives

There are several other ways to do this.

Within 'Server Capabilities' the API endpoints could be added to existing Server and Delivery Service endpoints. We are proposing they be standalone for now, because it's unclear how long-lived this feature will be, and how it will interact with 'Flexible Cachegroups' and other planned projects. Standalone endpoints give us more flexibility, especially around API versioning, to change or remove the feature in the future.

Instead of 'Capabilitites,' the Mid assignment problem could be solved in other ways, by making server and/or assignments and parentage more powerful. There are myriad ways to do that, and we need to be very thoughtful, to make sure the added complexity is as simple and manageable as possible, as it would be very easy to make TC even more difficult to operate.

Server Capabilities seems the simplest way to solve the immediate problem, the most flexible to also solve other similar problems, and also the easiest to adapt or remove with any large future change to server assignment and parentage.

## Dependencies

n/a

## References

n/a
