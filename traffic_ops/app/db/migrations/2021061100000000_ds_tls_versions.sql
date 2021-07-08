-- syntax:postgresql
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
CREATE TABLE IF NOT EXISTS public.deliveryservice_tls_version (
	deliveryservice bigint NOT NULL REFERENCES public.deliveryservice(id) ON DELETE CASCADE ON UPDATE CASCADE,
	tls_version text NOT NULL CHECK (tls_version <> ''),
	PRIMARY KEY (deliveryservice, tls_version)
);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_ds_timestamp_on_insert()
	RETURNS trigger
	AS $$
BEGIN
	UPDATE public.deliveryservice
	SET last_updated=now()
	WHERE id IN (
		SELECT deliveryservice
		FROM new_table
	);
	RETURN NULL;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_ds_timestamp_on_delete()
	RETURNS trigger
	AS $$
BEGIN
	UPDATE public.deliveryservice
	SET last_updated=now()
	WHERE id IN (
		SELECT deliveryservice
		FROM old_table
	);
	RETURN NULL;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER update_ds_timestamp_on_tls_version_insertion
	AFTER INSERT ON public.deliveryservice_tls_version
	REFERENCING NEW TABLE AS new_table
	FOR EACH STATEMENT EXECUTE PROCEDURE update_ds_timestamp_on_insert();

CREATE TRIGGER update_ds_timestamp_on_tls_version_delete
	AFTER DELETE ON public.deliveryservice_tls_version
	REFERENCING OLD TABLE AS old_table
	FOR EACH STATEMENT EXECUTE PROCEDURE update_ds_timestamp_on_delete();

UPDATE public.deliveryservice_request
SET
	deliveryservice = jsonb_set(deliveryservice, '{tlsVersions}', 'null')
WHERE
	deliveryservice IS NOT NULL;
UPDATE public.deliveryservice_request
SET
	original = jsonb_set(original, '{tlsVersions}', 'null')
WHERE
	original IS NOT NULL;

-- +goose Down
UPDATE public.deliveryservice_request
SET
	deliveryservice = deliveryservice - 'tlsVersions'
WHERE
	deliveryservice IS NOT NULL;

UPDATE public.deliveryservice_request
SET
	original = original - 'tlsVersions'
WHERE
	original IS NOT NULL;

DROP TRIGGER IF EXISTS update_ds_timestamp_on_tls_version_insertion_or_update ON public.deliveryservice_tls_version;
DROP TRIGGER IF EXISTS update_ds_timestamp_on_tls_version_delete ON public.deliveryservice_tls_version;
DROP TABLE IF EXISTS public.deliveryservice_tls_version;
DROP FUNCTION IF EXISTS update_ds_timestamp_on_insert;
DROP FUNCTION IF EXISTS update_ds_timestamp_on_delete;
