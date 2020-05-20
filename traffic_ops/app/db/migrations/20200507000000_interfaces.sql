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
CREATE TABLE interface (
	max_bandwidth bigint DEFAULT NULL CHECK (max_bandwidth IS NULL OR max_bandwidth >= 0),
	monitor boolean NOT NULL,
	mtu bigint DEFAULT 1500 CHECK (mtu IS NULL OR mtu > 1280),
	name text NOT NULL CHECK (name != ''),
	server bigint NOT NULL,
	PRIMARY KEY (name, server),
	FOREIGN KEY (server) REFERENCES server(id) ON DELETE RESTRICT ON UPDATE CASCADE
);
CREATE TABLE ip_address (
	address inet NOT NULL,
	gateway inet CHECK (
		gateway IS NULL OR (
			family(gateway) = 4 AND
			masklen(gateway) = 32
		) OR (
			family(gateway) = 6 AND
			masklen(gateway) = 128
		)
	),
	interface text NOT NULL,
	server bigint NOT NULL,
	service_address boolean NOT NULL DEFAULT FALSE,
	PRIMARY KEY (address, interface, server),
	FOREIGN KEY (server) REFERENCES server(id) ON DELETE RESTRICT ON UPDATE CASCADE,
	FOREIGN KEY (interface, server) REFERENCES interface(name, server) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TYPE server_ip_address AS (address inet, gateway inet, service_address boolean);
CREATE TYPE server_interface AS (ip_addresses server_ip_address ARRAY, max_bandwidth bigint, monitor boolean, mtu bigint, name text);

INSERT INTO interface(
	max_bandwidth,
	monitor,
	mtu,
	name,
	server
)
SELECT NULL, TRUE, server.interface_mtu::bigint, server.interface_name, id
FROM server;

INSERT INTO ip_address(
	address,
	gateway,
	interface,
	server,
	service_address
)
SELECT
	set_masklen(
		server.ip_address::inet,
		CASE
			WHEN server.ip_netmask IS NULL THEN 32
			WHEN server.ip_netmask = '' THEN 32
			ELSE
				(SELECT SUM(
					get_bit(digit, 0) +
					get_bit(digit, 1) +
					get_bit(digit, 2) +
					get_bit(digit, 3) +
					get_bit(digit, 4) +
					get_bit(digit, 5) +
					get_bit(digit, 6) +
					get_bit(digit, 7)
				)::int FROM (
					SELECT regexp_split_to_table::int::bit(8) AS digit
					FROM regexp_split_to_table(server.ip_netmask, '\.')
				) AS digits)
		END
	),
	NULLIF(regexp_replace(server.ip_gateway, '/\d+$', ''), '')::inet,
	server.interface_name,
	server.id,
	server.ip_address_is_service
FROM server
WHERE server.ip_address IS NOT NULL
AND server.ip_address != '';

INSERT INTO ip_address(
	address,
	gateway,
	interface,
	server,
	service_address
)
SELECT
	trim(BOTH '[]' FROM server.ip6_address)::inet,
	NULLIF(regexp_replace(server.ip6_gateway, '/\d+$', ''), '')::inet,
	server.interface_name,
	server.id,
	server.ip6_address_is_service
FROM server
WHERE server.ip6_address IS NOT NULL
AND server.ip6_address != '';

-- +goose Down
DROP TABLE IF EXISTS ip_address;
DROP TABLE IF EXISTS interface;

DROP TYPE IF EXISTS server_interface;
DROP TYPE IF EXISTS server_ip_address;
