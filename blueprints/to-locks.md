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
# CDN Traffic Ops Locks

## Problem Description

Currently, there is no way to guarantee that only your changes will make their way into a snapshot or queue updates 
for the servers of a CDN. If somebody else makes changes while you're in the middle of modifying your CDN components
or properties, their changes will "dirty" your snapshot and change the integrity of the data that you'd expect in your 
snapshot. This could result in data mismatches and inconsistencies.

CDN locks will serve as a way to avoid data corruption in snapshots and queue update activities by ensuring
that only the intended user(s) is able to pass their changes to the snapshot/ queue updates. By placing locks 
on a CDN, a user can essentially block out other users from modifying data that would dirty their view of the 
CDN.

This will be an optional feature in Traffic Ops/ Traffic Portal, which can be made use of by a user for their 
own guarantee of a "clean slate".

## Proposed Change

We propose a locking mechanism by which you can guarantee the exclusivity of your changes.
With CDN locks, you can ensure one or both of the following conditions:
- Only you can make changes and snap/ queue your CDN (exclusive locks)
- Anybody can make changes but only you can snap/ queue your CDN (shared locks)

`Exclusive` locks will help in cases where you are not collaborating with other users on the changes that you are
making to a CDN. We anticipate all automated scripts and cron jobs to make use of these exclusive locks.

The second category of locks is `shared` locks. This will be useful in a scenario when a group of users is working
on different parts of a CDN, but their changes are not related to each other. In such a case, any number of users can 
grab the shared lock and make changes on the CDN. However, only the first user who grabbed the shared lock can actually
snap/ queue updates on the CDN. This way, the first user to grab the lock still has the chance to review the changes, and
also ensure that nobody else snaps/ queues the CDN before they are done with their intended modifications.

There is also the third option of not using locks at all, in which case, the software will behave exactly the way it does
today, that is, no safety that your changes will not be corrupted by somebody else before you snap/ queue.

For unlocking a CDN under normal circumstances, only the user who has acquired the lock can unlock it (both shared 
and exclusive). However, in the rare case that a user `A` forgets to unlock a CDN, an `admin` role user can unlock 
the CDN on behalf of user `A`.

### Traffic Portal Impact

- Landing Page
    - New "lock" icon to give the user the ability to lock a particular CDN
    - A dropdown list that appears when you click the above mentioned lock icon, which lists the list of CDNs that you can lock
    - A `Message` field that appears under the dropdown list that the user can populate with a custom message stating the reason behind locking the CDN
    - A `Shared` field that can be set to `true` or `false` based on user requirement
    - A CDN notification displaying which CDN is locked by which user
    - An "unlock" option that appears next to the notification, only for the user who has locked the CDN
    - A way for the `admin` user to be able to unlock CDNs on other users' behalf
    
- Snapshot/ Queue Updates Page
    - A way to display an error on a snap/ queue operation by one user, if another user has the lock on that CDN

### Traffic Ops Impact

`/cdn_locks`
- List all CDN Locks
- GET+POST+DELETE

Traffic Ops will need to add the logic to check for locks before snapping/ queueing a CDN. It'll also need to account for
`shared` vs `exclusive` locks, and forbid a user from snapping/ queueing a CDN if another user possesses a lock on that CDN.

#### REST API Impact

The following is the JSON representation of a `CDN_Lock` object:
```JSON
{
  "userName": "foo",
  "cdnName": "cdn1",
  "message": "snapping cdn",
  "shared": false,
  "creator": true,
  "lastUpdated": "2021-05-10 16:03:34-06"
}
```

The following table describes the top level `CDN_Lock` object:

| field       | type                        | optionality | description                                                                             |
| ----------- | --------------------------- | ----------- | ----------------------------------------------------------------------------------------|
| userName    | string                      | required    | the user name of the user that wants to acquire a lock on the CDN                       |
| cdnName     | string                      | required    | the name of the CDN on which the lcok needs to be acquired                              |
| message     | string                      | optional    | the message stating a reason behind locking the CDN                                     |
| shared      | boolean                     | required    | whether or not this is a shared lock                                                    |
| creator     | boolean                     | optional    | whether or not the requesting `userName` is the first one to acquire a lock on `cdnName`| 
| lastUpdated | time                        | optional    | the last time this lock was updated                                                     |

**API constraints:**
- a `userName` and `cdnName` combination must be unique id `shared` is set to `false`. In other words, there can be only one `exclusive`
lock for a CDN
- a CDN can have multiple `shared` locks
- the `creator` will be set to `true` if the user is the first one acquiring a lock for the specified `cdnName`
- a user can snap/ queue a CDN only if the `creator` field corresponding to the `userName` and `cdnName` combination is `true`
- a user can delete their `exclusive` lock whenever, since the user is the creator in this case
- a user who is the creator cannot delete their `shared` lock until all other non-creator users have releases their shared locks for that CDN
- an `admin` user can delete the locks of other users on any CDN

Three new endpoints will be added for `GET`, `POST` and `DELETE` functionality with respect to CDN locks.
##### GET `cdn_locks`

response JSON:
```JSON
{
  "response": [
    {
      "userName": "foo",
      "cdnName": "cdn1",
      "message": "snapping cdn",
      "shared": true,
      "creator": true,
      "lastUpdated": "2021-05-10 16:03:34-06"
    },
    {
      "userName": "bar",
      "cdnName": "cdn2",
      "message": "queue cdn",
      "shared": false,
      "creator": true,
      "lastUpdated": "2021-05-10 17:04:34-06"
    }
  ]
}
```

##### POST `cdn_locks`

request JSON:
```JSON
{
  "cdnName": "bar",
  "message": "snapping cdn",
  "shared": false
}
```

response JSON:
```JSON
{
  "alerts": [
    {
      "text": "CDN lock acquired!",
      "level": "success"
    }
  ],
  "response": {
    "userName": "foo",
    "cdnName": "bar",
    "message": "snapping cdn",
    "shared": false,
    "creator": true,
    "lastUpdated": "2021-05-10 17:05:30-06"
  }
}
```

##### `DELETE /cdn_locks?cdn=bar`

response JSON:
```JSON
{
  "alerts": [
    {
      "text": "Cdn lock deleted",
      "level": "success"
    }
  ]
}
```

#### Client Impact

New Go client methods will be added for the `/cdn_locks` endpoints in order to write TO API tests for the new endpoints.

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
No impact

### Traffic Monitor Impact
No impact

### Traffic Router Impact
No impact

### Traffic Stats Impact
No impact

### Traffic Vault Impact
No impact

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
