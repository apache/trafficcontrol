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

ALTER TABLE public.deliveryservice_request
DROP CONSTRAINT appropriate_requested_and_original_for_change_type;

UPDATE public.deliveryservice_request
SET deliveryservice = original
WHERE deliveryservice IS NULL;

ALTER TABLE public.deliveryservice_request
ALTER COLUMN deliveryservice
DROP DEFAULT;

ALTER TABLE public.deliveryservice_request
ALTER COLUMN deliveryservice
SET NOT NULL;

ALTER TABLE public.deliveryservice_request
DROP COLUMN original;
