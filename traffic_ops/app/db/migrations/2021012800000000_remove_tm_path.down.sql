/*

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

        http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

/*
This migration removes the 'tm_path' property from the 'stats' property of the
'crconfig' column of stored CDN Snapshots, if it exists.

When reverted, it will insert (if possible) a 'tm_path' value of
/api/4.0/cdns/{{CDN Name}}/snapshot
where 'CDN Name' is the name of the CDN snapshotted as determined by the
Snapshot data - NOT the linked CDN object. This is so it does not self-confilct
afterward regarding which CDN is named, even if the one to which it is linked
is wrong, somehow.
*/

UPDATE snapshot SET crconfig = jsonb_set(crconfig::jsonb, '{stats,tm_path}', ('"/api/4.0/cdns/' || (crconfig::jsonb -> 'stats' ->> 'CDN_name') || '/snapshot"')::jsonb)
WHERE crconfig::jsonb ? 'stats' AND (crconfig::jsonb -> 'stats') ? 'CDN_name';
