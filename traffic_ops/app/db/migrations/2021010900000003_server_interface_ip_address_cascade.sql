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
ALTER TABLE interface
  DROP CONSTRAINT interface_server_fkey,
  ADD CONSTRAINT interface_server_fkey FOREIGN KEY (server) REFERENCES server(id) ON DELETE CASCADE ON UPDATE CASCADE;

-- +goose StatementBegin
DO $$
BEGIN
  IF EXISTS(
    SELECT conname
    FROM pg_constraint
    WHERE conrelid = (
      SELECT oid
      FROM pg_class
      WHERE relname LIKE 'ip_address'
    ) AND conname = 'ip_address_interface_server_fkey'
  ) THEN
    ALTER TABLE ip_address
      RENAME CONSTRAINT ip_address_interface_server_fkey TO ip_address_interface_fkey;
  END IF;
END
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

ALTER TABLE ip_address
  DROP CONSTRAINT ip_address_interface_fkey,
  ADD CONSTRAINT ip_address_interface_fkey FOREIGN KEY (server) REFERENCES server(id) ON DELETE CASCADE ON UPDATE CASCADE,
  DROP CONSTRAINT ip_address_server_fkey,
  ADD CONSTRAINT ip_address_server_fkey FOREIGN KEY (interface, server) REFERENCES interface(name, server) ON DELETE CASCADE ON UPDATE CASCADE;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE interface
  DROP CONSTRAINT interface_server_fkey,
  ADD CONSTRAINT interface_server_fkey FOREIGN KEY (server) REFERENCES server(id) ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE ip_address
  DROP CONSTRAINT ip_address_interface_fkey,
  ADD CONSTRAINT ip_address_interface_fkey FOREIGN KEY (server) REFERENCES server(id) ON DELETE RESTRICT ON UPDATE CASCADE,
  DROP CONSTRAINT ip_address_server_fkey,
  ADD CONSTRAINT ip_address_server_fkey FOREIGN KEY (interface, server) REFERENCES interface(name, server) ON DELETE RESTRICT ON UPDATE CASCADE;
