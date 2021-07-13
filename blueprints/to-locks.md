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

1. Currently, there is no way to guarantee that only your changes will make their way into a snapshot of a CDN. 
If somebody else makes changes while you're in the middle of modifying your CDN components or properties, 
their changes will "dirty" your snapshot and change the integrity of the data that you'd expect in your 
snapshot.

2. There is no way to prevent another user from performing a snapshot and prematurely propagating the (incomplete) changes 
   to Traffic Router and Traffic Monitor. 

3. There is no way to prevent other users from queueing updates on servers, and prematurely propagating the (incomplete) changes 
   to the Traffic Servers.

All the above three cases could result in data mismatches and inconsistencies.
CDN locks will serve as a way to avoid data corruption in snapshots and queue update activities by ensuring
that only the intended user(s) is able to pass their changes to the snapshot/ queue updates. By placing locks 
on a CDN, a user can essentially block out other users from modifying data that would dirty their view of the 
CDN.

This will be an optional feature in Traffic Ops/ Traffic Portal, which can be made use of by a user for their 
own guarantee of a "clean slate".

## Proposed Change

We propose a locking mechanism by which you can guarantee the exclusivity of your changes.
With CDN locks, you can ensure one or both of the following conditions:
- Only you can make changes and snap/ queue your CDN (hard locks)
- Anybody can make changes but only you can snap/ queue your CDN (soft locks)

`Hard` locks will help in cases where you are not collaborating with other users on the changes that you are
making to a CDN. We anticipate all automated scripts and cron jobs to make use of these hard locks.

The second category of locks is `soft` locks. This will be useful in a scenario when a group of users is working
on different parts of a CDN, but their changes are not related to each other. In such a case, say user `A` has acquired the `soft` lock
on cdn `foo`. Now, any number of users (including `A`) can make changes on the CDN. However, only the first user who grabbed the soft lock,
in this case, user `A`, can actually snap/ queue updates on the CDN. This way, the first user to grab the lock still has the chance to review 
the changes, and also ensure that nobody else snaps/ queues the CDN before they are done with their intended modifications.

There is also the third option of not using locks at all, in which case, the software will behave exactly the way it does
today, that is, no safety that your changes will not be corrupted by somebody else before you snap/ queue.

For unlocking a CDN under normal circumstances, only the user who has acquired the lock can unlock it (both soft 
and hard). However, in the rare case that a user `A` forgets to unlock a CDN, an `admin` role user can unlock 
the CDN on behalf of user `A`.

### Traffic Portal Impact

- Landing Page
    - New "lock" icon to give the user the ability to lock a particular CDN
    - A dropdown list that appears when you click the above mentioned lock icon, which lists the list of CDNs that you can lock
    - A `Message` field that appears under the dropdown list that the user can populate with a custom message stating the reason behind locking the CDN
    - A set of radio buttons to show the type of lock that the user wants, that is, `Soft Lock` and `Hard Lock`. A tooltip explaining what each lock denotes will also be provided.
    - A CDN notification displaying which CDN is locked by which user
    - An "unlock" option that appears next to the notification, only for the user who has locked the CDN
    - A way for the `admin` user to be able to unlock CDNs on other users' behalf
    
- Snapshot/ Queue Updates Page
    - Disable the snap/ queue button if another user has the lock on that CDN.

### Traffic Ops Impact

`/cdn_locks`
- List all CDN Locks
- GET+POST+DELETE

Traffic Ops will need to add the logic to check for locks before snapping/ queueing a CDN. It'll also need to account for
`soft` vs `hard` locks, and forbid a user from snapping/ queueing a CDN if another user possesses a lock on that CDN.
The following endpoints will need to handle the locks logic:
- `/cachegroups/{{ID}}/queue_update`
- `/cdns/{{ID}}/queue_update`
- `/servers/{{hostname}}/queue_update`
- `/snapshot`
- `/topologies/{{name}}/queue_update`

In addition to these, all `PUT`, `POST`, `DELETE` endpoints that scope to a CDN will have to add in the logic to check if there 
is a "hard" lock by some other user, before a user can modify data.
Basically, every endpoint that affects data related to a CDN directly (for example, `profiles`), and every endpoint that affects data 
related to a CDN indirectly (for example, `parameters`, `cachegroups`) will need to be updated to perform the above mentioned check.

#### REST API Impact

The following is the JSON representation of a `CDN_Lock` object:
```JSON
{
  "username": "foo",
  "cdn": "cdn1",
  "message": "snapping cdn",
  "soft": false,
  "lastUpdated": "2021-05-10T16:03:34-06:00"
}
```

The following table describes the top level `CDN_Lock` object:

| field       | type                        | optionality | description                                                                                                                  |
| ----------- | --------------------------- | ----------- | -----------------------------------------------------------------------------------------------------------------------------|
| username    | string                      | required    | the user name of the user that wants to acquire a lock on the CDN                                                            |
| cdn         | string                      | required    | the name of the CDN on which the lcok needs to be acquired                                                                   |
| message     | string                      | optional    | the message stating a reason behind locking the CDN                                                                          |
| soft        | boolean                     | required    | whether or not this is a shared lock, meaning if a user has the lock, whether or not other users can make changes to the CDN |
| lastUpdated | time                        | required    | the last time this lock was updated                                                                                          |

**API constraints:**
- a CDN can have only one `hard` or `soft` lock at a time
- a user can snap/ queue a CDN only if they are the one holding the lock (of either kind) on the `cdn`. Alternatively, if no one has the lock on the `cdn`, anyone can snap/ queue/ make changes to the CDN
- a user can delete their `hard` or `soft` lock whenever they want
- a user with at least `all-read` and `all-write` capabilities can delete the locks of other users on any CDN. The `all-write` capability will be modified
to include a new capability `delete-all-locks`. The way it works today, this translates to a user with an `admin` role. Such a user can delete other users' locks. 

Three new endpoints will be added for `GET`, `POST` and `DELETE` functionality with respect to CDN locks.
##### GET `cdn_locks`

response JSON:
```JSON
{
  "response": [
    {
      "username": "foo",
      "cdn": "cdn1",
      "message": "snapping cdn",
      "soft": true,
      "lastUpdated": "2021-05-10T16:03:34-06:00"
    },
    {
      "username": "bar",
      "cdn": "cdn2",
      "message": "queue cdn",
      "soft": false,
      "lastUpdated": "2021-05-10T17:04:34-06:00"
    }
  ]
}
```

##### POST `cdn_locks`

request JSON:
```JSON
{
  "cdn": "bar",
  "message": "snapping cdn",
  "soft": false
}
```

response JSON:
```JSON
{
  "alerts": [
    {
      "text": "hard CDN lock acquired!",
      "level": "success"
    }
  ],
  "response": {
    "username": "foo",
    "cdn": "bar",
    "message": "snapping cdn",
    "soft": false,
    "lastUpdated": "2021-05-10T17:05:30-06:00"
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
  ],
  "response": {
    "username": "foo",
    "cdn": "bar",
    "message": "snapping cdn",
    "soft": false,
    "lastUpdated": "2021-05-10T17:05:30-06:00"
  }
}
```

#### Client Impact

New Go client methods will be added for the `/cdn_locks` endpoints in order to write TO API tests for the new endpoints.

#### Data Model / Database Impact
A new database table for `cdn_lock`, as described below, will be created.
```text
            Table "traffic_ops.cdn_lock"
     Column    |  Type                    | Collation | Nullable | Default
---------------+--------------------------+-----------+----------+--------
 id            | bigint                   |           | not null |
 username      | text                     |           | not null |
 cdn_name      | text                     |           | not null |
 message       | test                     |           |          |
 soft          | boolean                  |           | not null | true
 last_updated  | timestamp with time zone |           | not null | now() 
Indexes:
    "pk_cdn_lock" PRIMARY KEY(cdn)
Foreign-key constraints:
    "fk_lock_cdn" FOREIGN KEY (cdn) REFERENCES cdn(name)
    "fk_lock_username" FOREIGN KEY (username) REFERENCES tm_user(username)
```

The `capability` table will need to add another capability by the name of `delete-all-locks` and the `all-write` capability 
will include this new capability. 
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
All new endpoints will need to be documented, along with the documentation explaining `cdn_locks`. 

### Testing Impact
Client/API integration tests should be written to verify the
functionality as described in [API](#sec:api).

### Performance Impact
We do not anticipate significant performance impact, as this process will just entail one additional check
in the database before performing a snap or queue operation.

### Security Impact
We do not anticipate any impact on security due to the introduction of `cdn_locks`.

### Upgrade Impact
There will be one database migrations, to create the `cdn_locks` table. 
The new capability can just be added to the `seeds.sql` file. 
This does not depend on any existing data, so nothing should ideally cause this migration to fail.

### Operations Impact
Operations will have to learn how to create, delete and work with CDN locks, if they wish to make use of this feature.
If not, there should be no impact to their use of the software.

### Developer Impact
Developers will most likely need to use CDN locks, so they'll need to be familiar with the process
of creating, deleting, debugging and working with locks.

## Alternatives
No other alternatives for this safety measure currently exist, other than posting a message on a
group chat or sending an email to ask everyone to refrain from making changes/ snapping/ queueing.

## Dependencies
None

## References
None
