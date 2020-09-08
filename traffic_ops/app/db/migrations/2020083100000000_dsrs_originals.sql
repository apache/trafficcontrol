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
ALTER TABLE deliveryservice_request
ADD COLUMN original jsonb DEFAULT NULL;

UPDATE deliveryservice_request
SET original=deliveryservice
WHERE status = 'complete' OR status = 'rejected' OR change_type = 'delete';

ALTER TABLE deliveryservice_request
ADD CONSTRAINT closed_has_original
CHECK (
	(
		(status = 'complete' OR status = 'rejected' OR change_type = 'delete') AND
		original IS NOT NULL
	)
	OR
	(
		(status = 'submitted' OR status = 'pending' OR status='draft') AND
		change_type <> 'delete' AND
		original IS NULL
	)

);

-- +goose Down
ALTER TABLE deliveryservice_request
DROP COLUMN original;
