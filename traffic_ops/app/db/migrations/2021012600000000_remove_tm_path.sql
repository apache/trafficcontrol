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
UPDATE snapshot
SET crconfig = crconfig::jsonb #- '{stats,tm_path}'
WHERE crconfig::jsonb ? 'stats' AND (crconfig::jsonb -> 'stats') ? 'tm_path';

-- +goose Down
UPDATE snapshot SET crconfig = jsonb_set(crconfig::jsonb, '{stats,tm_path}', ('"/api/4.0/cdns/' || (crconfig::jsonb -> 'stats' ->> 'CDN_name') || '/snapshot"')::jsonb)
WHERE crconfig::jsonb ? 'stats' AND (crconfig::jsonb -> 'stats') ? 'CDN_name';
