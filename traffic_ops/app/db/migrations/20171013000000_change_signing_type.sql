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
DROP DOMAIN IF EXISTS deliveryservice_signature_type RESTRICT;

CREATE DOMAIN deliveryservice_signature_type AS TEXT
CHECK(
   VALUE IN ('url_sig', 'uri_signing')
);

ALTER TABLE IF EXISTS deliveryservice
ALTER COLUMN signed
SET DEFAULT NULL;

ALTER TABLE IF EXISTS deliveryservice
ALTER COLUMN signed
SET DATA TYPE deliveryservice_signature_type
USING CASE WHEN signed THEN 'url_sig'::text::deliveryservice_signature_type ELSE NULL END;

ALTER TABLE IF EXISTS deliveryservice
RENAME COLUMN signed TO signing_algorithm;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE IF EXISTS deliveryservice
RENAME COLUMN signing_algorithm TO signed;

ALTER TABLE deliveryservice
ALTER COLUMN signed
SET DATA TYPE boolean
USING CASE WHEN signed='url_sig' THEN true ELSE false END;

ALTER TABLE IF EXISTS deliveryservice
ALTER COLUMN signed
SET DEFAULT false;

DROP DOMAIN IF EXISTS deliveryservice_signature_type RESTRICT;
