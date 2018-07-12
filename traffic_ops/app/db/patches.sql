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


-- THIS FILE INCLUDES POST-MIGRATION DATA FIXES REQUIRED OF TRAFFIC OPS
UPDATE steering_target SET "type" = (SELECT id FROM type WHERE name = 'STEERING_WEIGHT') where "type" is NULL;;
ALTER TABLE steering_target ALTER COLUMN type SET NOT NULL;

UPDATE deliveryservice SET routing_name = 'edge' WHERE routing_name IS NULL AND type IN (SELECT id FROM type WHERE name like 'DNS%');
UPDATE deliveryservice ds
SET routing_name = (
  SELECT p.value
  FROM parameter p
    JOIN profile_parameter pp ON p.id = pp.parameter
    JOIN profile pro ON pp.profile = pro.id
    JOIN cdn ON pro.cdn = cdn.id
  WHERE p.name = 'upgrade_http_routing_name'
    AND cdn.id = ds.cdn_id)
WHERE routing_name IS NULL;
UPDATE deliveryservice SET routing_name = 'tr' WHERE routing_name IS NULL;
ALTER TABLE deliveryservice ALTER COLUMN routing_name SET NOT NULL;
ALTER TABLE deliveryservice ALTER COLUMN routing_name SET DEFAULT 'cdn';
ALTER TABLE deliveryservice DROP CONSTRAINT IF EXISTS routing_name_not_empty;
ALTER TABLE deliveryservice ADD CONSTRAINT routing_name_not_empty CHECK (length(routing_name) > 0);
-- This patch adds a new enum value to profile_type
-- We need to add the enum directly with sql as postgres will not allow
-- altering any enum values within a transaction block and by default
-- goose runs migrations within a transaction.
ALTER TYPE profile_type ADD VALUE IF NOT EXISTS 'GROVE_PROFILE';

-- If tenancy is not needed, simply default to the root tenant.  This will help
-- to eliminate the use_tenancy property

DO
$tenantnotnull$
BEGIN
    IF EXISTS (SELECT value FROM parameter WHERE name = 'use_tenancy' AND config_file = 'global' AND value = '1') THEN
        -- tenancy turned ON -- anything without a tenant gets the "unassigned" tenant
        UPDATE tm_user SET tenant_id = (SELECT id FROM tenant WHERE name = 'unassigned')
        WHERE tenant_id IS NULL;
        UPDATE deliveryservice SET tenant_id = (SELECT id FROM tenant WHERE name = 'unassigned')
        WHERE tenant_id IS NULL;
        UPDATE origin SET tenant = (SELECT id FROM tenant WHERE name = 'unassigned')
        WHERE tenant IS NULL;
    ELSE
        -- tenancy turned OFF -- anything without a tenant gets the "root" tenant
        UPDATE tm_user SET tenant_id = (SELECT id FROM tenant WHERE name = 'root')
        WHERE tenant_id IS NULL;
        UPDATE deliveryservice SET tenant_id = (SELECT id FROM tenant WHERE name = 'root')
        WHERE tenant_id IS NULL;
        UPDATE origin SET tenant = (SELECT id FROM tenant WHERE name = 'root')
        WHERE tenant IS NULL;
    END IF;
END
$tenantnotnull$;

-- add constraints so we can depend on everything being assigned a tenant
ALTER TABLE tm_user ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE deliveryservice ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE origin ALTER COLUMN tenant SET NOT NULL;
-- get rid of the use_tenancy flag
DELETE FROM parameter WHERE name = 'use_tenancy' AND config_file = 'global';
