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

INSERT INTO capability (name, description) values ('parameters-read-secure', 'Ability to view secure parameter values') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) values ('servers-read-secure', 'Ability to view secure server values') ON CONFLICT (name) DO NOTHING;

INSERT INTO role_capability (role_id, cap_name) VALUES ((SELECT id FROM role WHERE name = 'admin'), 'parameters-read-secure') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) VALUES ((SELECT id FROM role WHERE name = 'admin'), 'servers-read-secure') ON CONFLICT DO NOTHING;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DELETE FROM role_capability WHERE role_id = (SELECT id FROM role WHERE name = 'admin') AND cap_name = 'servers-read-secure';
DELETE FROM role_capability WHERE role_id = (SELECT id FROM role WHERE name = 'admin') AND cap_name = 'parameters-read-secure';

DELETE FROM capability WHERE name = 'servers-read-secure'
DELETE FROM capability WHERE name = 'parameters-read-secure'
