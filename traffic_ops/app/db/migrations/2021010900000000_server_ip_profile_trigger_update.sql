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
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION before_server_table()
    RETURNS TRIGGER AS
$$
DECLARE
    server_count BIGINT;
BEGIN
    WITH server_ips AS (
        SELECT s.id, i.name, ip.address, s.profile
        FROM server s
                 JOIN interface i on i.server = s.ID
                 JOIN ip_address ip on ip.Server = s.ID and ip.interface = i.name
        WHERE i.monitor = true
    )
    SELECT count(*)
    INTO server_count
    FROM server_ips sip
             JOIN server_ips sip2 on sip.id <> sip2.id
    WHERE sip2.address = sip.address
      AND sip2.profile = sip.profile;

    IF server_count > 0 THEN
        RAISE EXCEPTION 'Server [id:%] does not have a unique ip_address over the profile [id:%], [%] conflicts',
            NEW.id,
            NEW.profile,
            server_count;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION before_ip_address_table()
    RETURNS TRIGGER
AS
$$
DECLARE
    server_count   BIGINT;
    server_id      BIGINT;
    server_profile BIGINT;
BEGIN
    WITH server_ips AS (
        SELECT s.id as sid, i.name, ip.address, s.profile
        FROM server s
                 JOIN interface i
                     on i.server = s.ID
                 JOIN ip_address ip
                     on ip.Server = s.ID and ip.interface = i.name
        WHERE i.monitor = true
    )
    SELECT count(distinct(sip.sid)), sip.sid, sip.profile
    INTO server_count, server_id, server_profile
    FROM server_ips sip
             JOIN server_ips sip2 on sip.sid <> sip2.sid
    WHERE (sip.sid <> NEW.server AND sip.address = NEW.address AND sip.name = NEW.interface)
    GROUP BY sip.sid, sip.profile;
    IF server_count > 0 THEN
        RAISE EXCEPTION 'ip_address is not unique accross the server [id:%] profile [id:%], [%] conflicts',
            server_id,
            server_profile,
            server_count;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE PLPGSQL;
-- +goose StatementEnd
