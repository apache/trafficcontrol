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
CREATE TYPE workflow_states AS ENUM ('draft', 'submitted', 'rejected', 'pending', 'complete');
CREATE TYPE change_types AS ENUM ('create', 'update', 'delete');

CREATE TABLE deliveryservice_request (
    assignee_id bigint,
    author_id bigint NOT NULL,
    change_type change_types NOT NULL,
    created_at timestamp WITH time zone NOT NULL DEFAULT now(),
    id bigserial primary key NOT NULL,
    last_edited_by_id bigint NOT NULL,
    last_updated timestamp WITH time zone NOT NULL DEFAULT now(),
    request jsonb NOT NULL,
    status workflow_states NOT NULL
);

ALTER TABLE deliveryservice_request
    ADD CONSTRAINT fk_author FOREIGN KEY (author_id) REFERENCES tm_user(id) ON DELETE CASCADE;

ALTER TABLE deliveryservice_request
    ADD CONSTRAINT fk_assignee FOREIGN KEY (assignee_id) REFERENCES tm_user(id) ON DELETE SET NULL;

ALTER TABLE deliveryservice_request
    ADD CONSTRAINT fk_last_edited_by FOREIGN KEY (last_edited_by_id) REFERENCES tm_user(id) ON DELETE CASCADE;

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON deliveryservice_request FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE deliveryservice_request;

DROP TYPE change_types;

DROP TYPE workflow_states;
