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

INSERT INTO coordinate (name, latitude, longitude)
SELECT
'from_cachegroup_' || cg.name,
cg.latitude,
cg.longitude
FROM cachegroup cg
WHERE cg.latitude IS NOT NULL AND cg.longitude IS NOT NULL;

ALTER TABLE cachegroup ADD COLUMN coordinate bigint REFERENCES coordinate(id);

UPDATE cachegroup
SET coordinate = (
    SELECT co.id
    FROM coordinate co
    WHERE co.name = 'from_cachegroup_' || cachegroup.name
);

CREATE INDEX cachegroup_coordinate_fkey ON cachegroup USING btree (coordinate);

ALTER TABLE cachegroup DROP COLUMN latitude;
ALTER TABLE cachegroup DROP COLUMN longitude;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE cachegroup ADD COLUMN latitude numeric;
ALTER TABLE cachegroup ADD COLUMN longitude numeric;

UPDATE cachegroup
SET latitude = (
    SELECT co.latitude
    FROM coordinate co
    WHERE cachegroup.coordinate = co.id
),
longitude = (
    SELECT co.longitude
    FROM coordinate co
    WHERE cachegroup.coordinate = co.id
)
WHERE cachegroup.coordinate IS NOT NULL;

ALTER TABLE cachegroup DROP CONSTRAINT cachegroup_coordinate_fkey;

DELETE FROM coordinate
WHERE id IN (
    SELECT cg.coordinate FROM cachegroup cg
);

ALTER TABLE cachegroup DROP COLUMN coordinate;
