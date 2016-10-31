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

-- this migration only works if there is no "action" entry, which is true for us

alter table deliveryservice add column temp varchar(2048);
update deliveryservice set temp=(select action from header_rewrite where header_rewrite.id=deliveryservice.header_rewrite);
alter table deliveryservice drop foreign key fk_deliveryservice_header_rewrite1;
alter table deliveryservice drop key fk_deliveryservice_header_rewrite1_idx;
alter table deliveryservice drop column header_rewrite;
alter table deliveryservice change `temp` `header_rewrite` varchar(2048);
drop table header_rewrite;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
