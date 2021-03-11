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
ALTER TABLE server DROP CONSTRAINT need_at_least_one_ip;
ALTER TABLE server DROP CONSTRAINT need_gateway_if_ip;
ALTER TABLE server DROP CONSTRAINT need_netmask_if_ip;

ALTER TABLE server DROP COLUMN interface_name;
ALTER TABLE server DROP COLUMN ip_address;
ALTER TABLE server DROP COLUMN ip_netmask;
ALTER TABLE server DROP COLUMN ip_gateway;
ALTER TABLE server DROP COLUMN ip6_address;
ALTER TABLE server DROP COLUMN ip6_gateway;
ALTER TABLE server DROP COLUMN interface_mtu;
ALTER TABLE server DROP COLUMN ip_address_is_service;
ALTER TABLE server DROP COLUMN ip6_address_is_service;


-- +goose Down
ALTER TABLE server ADD COLUMN interface_name text DEFAULT '' NOT NULL;
ALTER TABLE server ADD COLUMN ip_address text DEFAULT '';
ALTER TABLE server ADD COLUMN ip_netmask text DEFAULT '';
ALTER TABLE server ADD COLUMN ip_gateway text DEFAULT '';
ALTER TABLE server ADD COLUMN ip6_address text DEFAULT '';
ALTER TABLE server ADD COLUMN ip6_gateway text DEFAULT '';
ALTER TABLE server ADD COLUMN interface_mtu bigint DEFAULT '9000'::bigint NOT NULL;
ALTER TABLE server ADD COLUMN ip_address_is_service boolean DEFAULT true;
ALTER TABLE server ADD COLUMN ip6_address_is_service boolean DEFAULT true;

ALTER TABLE server ADD CONSTRAINT need_at_least_one_ip CHECK (ip_address IS NOT NULL OR ip6_address IS NOT NULL OR ip_address = '' OR ip6_address = '');
ALTER TABLE server ADD CONSTRAINT need_gateway_if_ip CHECK (ip_address IS NULL OR ip_address = '' OR ip_gateway IS NOT NULL);
ALTER TABLE server ADD CONSTRAINT need_netmask_if_ip CHECK (ip_address IS NULL OR ip_address = '' OR ip_netmask IS NOT NULL);

UPDATE server SET ip_address = host(ip_address.address),
  ip_netmask = COALESCE(host(netmask(ip_address.address)), ''),
  ip_gateway = COALESCE(host(ip_address.gateway), ''),
  ip_address_is_service = ip_address.service_address,
  interface_name = ip_address.interface,
  interface_mtu = COALESCE(interface.mtu, 0)
  FROM ip_address
  JOIN interface ON ip_address.interface = interface.name
  WHERE server.id = ip_address.server
  AND family(ip_address.address) = 4
  AND ip_address.service_address;

UPDATE server SET ip6_address = ip_address.address,
  ip6_gateway = ip_address.gateway,
  ip6_address_is_service = ip_address.service_address
  FROM ip_address
  WHERE server.id = ip_address.server
  AND family(ip_address.address) = 6
  AND ip_address.service_address;
