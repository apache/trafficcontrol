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

-- we are going to clean up capabilities for a reseed
-- it is really just a few, but I'm deleting everything just so I don't miss anything..
DELETE FROM api_capability;
DELETE FROM role_capability;
DELETE FROM capability;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
