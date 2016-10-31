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
insert into job_status (name, description) values ('PENDING', 'Job is queued, but has not been picked up by any agents yet');
insert into job_status (name, description) values ('IN_PROGRESS', 'Job is being processed by agents');
insert into job_status (name, description) values ('COMPLETED', 'Job has finished');
insert into job_status (name, description) values ('CANCELLED', 'Job was cancelled');
insert into job_status (name, description) values ('PURGE', 'Initial Purge state');
insert into job_agent (name, description, active) values ('dummy','Description of Purge Agent','1');

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
