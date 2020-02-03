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
# TITLE <!-- a concise title for the new feature this blueprint will describe -->

## Problem Description
<!--
*What* is being asked for?
*Why* is this necessary?
*How* will this be used?
-->

## Proposed Change
<!--
*How* will this be implemented (at a high level)?
-->

### Traffic Portal Impact
<!--
*How* will this impact Traffic Portal?
What new UI changes will be required?
Will entirely new pages/views be necessary?
Will a new field be added to an existing form?
How will the user interact with the new UI changes?
-->

### Traffic Ops Impact
<!--
*How* will this impact Traffic Ops (at a high level)?
-->

#### REST API Impact
<!--
*How* will this impact the Traffic Ops REST API?

What new endpoints will be required?
How will existing endpoints be changed?
What will the requests and responses look like?
What fields are required or optional?
What are the defaults for optional fields?
What are the validation constraints?
-->

#### Client Impact
<!--
*How* will this impact Traffic Ops REST API clients (Go, Python, Java)?

If new endpoints are required, will corresponding client methods be added?
-->

#### Data Model / Database Impact
<!--
*How* will this impact the Traffic Ops data model?
*How* will this impact the Traffic Ops database schema?

What changes to the lib/go-tc structs will be required?
What new tables and columns will be required?
How will existing tables and columns be changed?
What are the column data types and modifiers?
What are the FK references and constraints?
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

### Traffic Stats Impact
<!--
*How* will this impact Traffic Stats?
-->

### Traffic Vault Impact
<!--
*How* will this impact Traffic Vault?

Will there be any new data stored in or removed from Riak?
Will there be any changes to the Riak requests and responses?
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
