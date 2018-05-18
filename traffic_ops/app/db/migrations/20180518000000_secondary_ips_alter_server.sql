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

-- alter server
ALTER TABLE server
    DROP COLUMN IF EXISTS interface_name,
    DROP COLUMN IF EXISTS interface_mtu,
    DROP COLUMN IF EXISTS ip_address,
    DROP COLUMN IF EXISTS ip_netmask,
    DROP COLUMN IF EXISTS ip_gateway,
    DROP COLUMN IF EXISTS ip6_address,
    DROP COLUMN IF EXISTS ip6_gateway,
    DROP COLUMN IF EXISTS mgmt_ip_address,
    DROP COLUMN IF EXISTS mgmt_ip_netmask,
    DROP COLUMN IF EXISTS mgmt_ip_gateway,
    DROP COLUMN IF EXISTS ilo_ip_address,
    DROP COLUMN IF EXISTS ilo_ip_netmask,
    DROP COLUMN IF EXISTS ilo_ip_gateway;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE server
    ADD COLUMN interface_name text NOT NULL,
    ADD COLUMN interface_mtu bigint DEFAULT '9000'::bigint NOT NULL,
    ADD COLUMN ip_address text NOT NULL,
    ADD COLUMN ip_netmask text NOT NULL,
    ADD COLUMN ip_gateway text NOT NULL,
    ADD COLUMN ip6_address text,
    ADD COLUMN ip6_gateway text,
    ADD COLUMN mgmt_ip_address text,
    ADD COLUMN mgmt_ip_netmask text,
    ADD COLUMN mgmt_ip_gateway text,
    ADD COLUMN ilo_ip_address text,
    ADD COLUMN ilo_ip_netmask text,
    ADD COLUMN ilo_ip_gateway text;
