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
CREATE TABLE IF NOT EXISTS service_category (
	name TEXT PRIMARY KEY CHECK (name <> ''),
	tenant_id BIGINT NOT NULL REFERENCES tenant(id),
	last_updated TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);
ALTER TABLE deliveryservice ADD COLUMN service_category TEXT REFERENCES service_category(name) ON UPDATE CASCADE;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE deliveryservice DROP COLUMN service_category;
DROP TABLE service_category;
