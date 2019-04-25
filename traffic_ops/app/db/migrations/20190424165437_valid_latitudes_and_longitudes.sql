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

-- Database constraint for valid lat-long values

-- +goose Up
CREATE DOMAIN latitude AS numeric CHECK (VALUE >= -90 AND VALUE <= 90);
CREATE DOMAIN longitude AS numeric CHECK (VALUE >= -180 AND VALUE <= 180);

UPDATE coordinate SET latitude=0 WHERE latitude>90 OR latitude<-90;
UPDATE coordinate SET longitude=0 WHERE longitude>180 OR longitude<-180;

ALTER TABLE coordinate ALTER COLUMN latitude SET DATA TYPE latitude;
ALTER TABLE coordinate ALTER COLUMN longitude SET DATA TYPE longitude;

-- +goose Down
ALTER TABLE coordinate ALTER COLUMN latitude SET DATA TYPE numeric;
ALTER TABLE coordinate ALTER COLUMN longitude SET DATA TYPE numeric;
DROP DOMAIN IF EXISTS latitude;
DROP DOMAIN IF EXISTS longitude;
