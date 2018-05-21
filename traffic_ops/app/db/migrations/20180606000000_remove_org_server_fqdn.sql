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

INSERT INTO origin (name, protocol, fqdn, port, deliveryservice, tenant, is_primary)
SELECT
d.xml_id,
lower(split_part(d.org_server_fqdn, '://', 1))::origin_protocol,
regexp_replace(d.org_server_fqdn, '(^https?://)|(:\d+$)', '', 'gi'),
(SELECT (regexp_matches(d.org_server_fqdn, '(?<=:)\d+$'))[1]::bigint),
d.id,
d.tenant_id,
TRUE
FROM deliveryservice d
WHERE d.org_server_fqdn IS NOT NULL;

ALTER TABLE deliveryservice DROP COLUMN org_server_fqdn;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE deliveryservice ADD COLUMN org_server_fqdn text;

UPDATE deliveryservice d
SET org_server_fqdn = (
    SELECT o.protocol::text || '://' || o.fqdn || rtrim(concat(':', o.port::text), ':')
    FROM origin o
    WHERE o.deliveryservice = d.id AND o.is_primary
);

DELETE FROM origin o
WHERE o.is_primary;
