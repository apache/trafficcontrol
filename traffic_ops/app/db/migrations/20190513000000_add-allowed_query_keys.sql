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

-- deliveryservice_consistent_hash_query_param
CREATE TABLE IF NOT EXISTS deliveryservice_consistent_hash_query_param (
    name TEXT NOT NULL,
    deliveryservice_id bigint NOT NULL,

    CONSTRAINT name_empty CHECK (length(name) > 0),
    CONSTRAINT name_reserved CHECK (name NOT IN ('format','trred')),
    CONSTRAINT fk_deliveryservice FOREIGN KEY (deliveryservice_id) REFERENCES deliveryservice(id) ON DELETE CASCADE,
    PRIMARY KEY (name, deliveryservice_id)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS deliveryservice_consistent_hash_query_param;
