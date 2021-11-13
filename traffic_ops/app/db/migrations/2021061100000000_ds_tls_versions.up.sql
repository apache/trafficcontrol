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

CREATE TABLE IF NOT EXISTS public.deliveryservice_tls_version (
	deliveryservice bigint NOT NULL REFERENCES public.deliveryservice(id) ON DELETE CASCADE ON UPDATE CASCADE,
	tls_version text NOT NULL CHECK (tls_version <> ''),
	PRIMARY KEY (deliveryservice, tls_version)
);

CREATE OR REPLACE FUNCTION update_ds_timestamp_on_insert()
	RETURNS trigger
	AS $$
BEGIN
	UPDATE public.deliveryservice
	SET last_updated=now()
	WHERE id IN (
		SELECT deliveryservice
		FROM CAST(NEW AS deliveryservice_tls_version)
	);
	RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_ds_timestamp_on_delete()
	RETURNS trigger
	AS $$
BEGIN
	UPDATE public.deliveryservice
	SET last_updated=now()
	WHERE id IN (
		SELECT deliveryservice
		FROM CAST(OLD AS deliveryservice_tls_version)
	);
	RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_ds_timestamp_on_tls_version_insertion
	AFTER INSERT ON public.deliveryservice_tls_version
	FOR EACH ROW EXECUTE PROCEDURE update_ds_timestamp_on_insert();

CREATE TRIGGER update_ds_timestamp_on_tls_version_delete
	AFTER DELETE ON public.deliveryservice_tls_version
	FOR EACH ROW EXECUTE PROCEDURE update_ds_timestamp_on_delete();

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
