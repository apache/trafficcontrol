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
