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
# Federations Redesign

## Problem Description
Currently, Federations are among the most convoluted and difficult to comprehend
concepts/objects in Traffic Control. The process to create a Federation - and
what one even does - is poorly understood by many, and that's reflected by the
fact that
[the latest documentation's quick-how-to guide on the subject](https://traffic-control-cdn.readthedocs.io/en/v6.0.0/admin/quick_howto/federations.html)
still has steps for and screenshots of the Traffic Ops UI that was removed years
ago.

Part of this is undoubtedly caused by the myriad of Traffic Ops API endpoints
used to manipulate Federations. Currently, these exist:

- `/cdns/{{name}}/federations`
- `/cdns/{{name}}/federations/{{ID}}`
- `/federations/all`
- `/federations`
- `/federations/{{ID}}/deliveryservices`
- `/federations/{{ID}}/deliveryservices/{{dsID}}`
- `/federation_resolvers`
- `/federations/{{ID}}/federation_resolvers`
- `/federations/{{ID}}/users`
- `/federations/{{ID}}/users/{{userID}}`

## Proposed Change
The notion of a "Federation" will be reworked into a first-class citizen of the
ATC data model. Access will be controlled by Tenancy (and Permissions), so the
only necessary Traffic Ops API endpoint will be `/federations`.

## Data Model Impact
The proposed model for Federations is as follows:

```typescript
interface Federation {
    cname: `${string}.`;
    description: string;
     // RFC3339 timestamp string - only present in responses
    lastUpdated?: Date;
    resolvers: Array<string>;
    ttl: number;
}
```


### Traffic Portal Impact
### Traffic Ops Impact

#### REST API Impact

#### Client Impact

#### Data Model / Database Impact

### ORT Impact
### Traffic Monitor Impact

### Traffic Router Impact

### Traffic Stats Impact

### Traffic Vault Impact

### Documentation Impact

### Testing Impact

### Performance Impact
### Security Impact

### Upgrade Impact

### Operations Impact

### Developer Impact

## Alternatives

## Dependencies

## References
