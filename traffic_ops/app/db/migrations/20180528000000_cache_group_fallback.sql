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


-- cachegroup_fallbacks
CREATE TABLE cachegroup_fallbacks (
    primary_cg bigint NOT NULL,
    backup_cg bigint CHECK (primary_cg != backup_cg) NOT NULL,
    set_order bigint NOT NULL,
    CONSTRAINT fk_primary_cg FOREIGN KEY (primary_cg) REFERENCES cachegroup(id) ON DELETE CASCADE,   
    CONSTRAINT fk_backup_cg FOREIGN KEY (backup_cg) REFERENCES cachegroup(id) ON DELETE CASCADE,
    UNIQUE (primary_cg, backup_cg),
    UNIQUE (primary_cg, set_order)
); 

ALTER TABLE cachegroup ADD COLUMN fallback_to_closest BOOLEAN DEFAULT TRUE;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE cachegroup DROP COLUMN fallback_to_closest;

DROP TABLE cachegroup_fallbacks;
