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

-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE public.cdn ADD COLUMN domain_name text;

UPDATE cdn SET domain_name=domainlist.value
  FROM (SELECT distinct cdn_id,value FROM server,parameter WHERE type=(SELECT id FROM type WHERE name='EDGE') 
    AND parameter.id in (select parameter from profile_parameter WHERE profile_parameter.profile=server.profile) 
    AND parameter.name='domain_name' 
    AND config_file='CRConfig.json') AS domainlist
WHERE id = domainlist.cdn_id;

UPDATE public.cdn SET domain_name='-' WHERE name='ALL';

ALTER TABLE public.cdn ALTER COLUMN domain_name SET NOT NULL;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE public.cdn DROP COLUMN domain_name;