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
ADD COLUMN original jsonb DEFAULT NULL;

UPDATE public.deliveryservice_request
SET original=deliveryservice
WHERE (status = 'complete' OR status = 'rejected' OR status = 'pending' OR change_type = 'delete') AND change_type != 'create';

ALTER TABLE public.deliveryservice_request
ALTER COLUMN deliveryservice
DROP NOT NULL;

ALTER TABLE public.deliveryservice_request
ALTER COLUMN deliveryservice
SET DEFAULT NULL;

UPDATE public.deliveryservice_request
SET deliveryservice=NULL
WHERE change_type='delete';

/* 'create' dsrs have no original, 'delete' dsrs have no requested, and only*/
/* closed 'update' dsrs have originals. */
ALTER TABLE public.deliveryservice_request
ADD CONSTRAINT appropriate_requested_and_original_for_change_type
CHECK (
	(change_type = 'delete' AND original IS NOT NULL AND deliveryservice IS NULL)
	OR
	(change_type = 'create' AND original IS NULL AND deliveryservice IS NOT NULL)
	OR (
		change_type = 'update' AND
		deliveryservice IS NOT NULL AND
		(
			(
				(status = 'complete' OR status = 'rejected' OR status = 'pending')
				AND
				original IS NOT NULL
			)
			OR
			(
				(status = 'draft' OR status = 'submitted')
				AND
				original IS NULL
			)
		)
	)
);
