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

rename table location to cachegroup;
rename table location_parameter to cachegroup_parameter;

alter table cachegroup drop key loc_name_UNIQUE;
alter table cachegroup add key cg_name_UNIQUE (name);
alter table cachegroup drop key loc_short_UNIQUE;
alter table cachegroup add key cg_short_UNIQUE (short_name);
alter table cachegroup drop foreign key fk_location_1;
alter table cachegroup drop key fk_location_1;
alter table cachegroup change parent_location_id parent_cachegroup_id int(11);
alter table cachegroup add CONSTRAINT `fk_cg_1` FOREIGN KEY (`parent_cachegroup_id`) REFERENCES `cachegroup` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION;
alter table cachegroup drop foreign key fk_location_type1;
alter table cachegroup drop key fk_location_type1;
alter table cachegroup add CONSTRAINT `fk_cg_type1` FOREIGN KEY (`type`) REFERENCES `type` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION;

alter table cachegroup_parameter drop foreign key fk_location;
alter table cachegroup_parameter drop key fk_location;
alter table cachegroup_parameter change location cachegroup int(11);
alter table cachegroup_parameter add CONSTRAINT `fk_cg_param_cachegroup1` FOREIGN KEY (`cachegroup`) REFERENCES `cachegroup` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION;

alter table server drop foreign key fk_contentserver_location;
alter table server drop  key fk_contentserver_location;
alter table server change location cachegroup int(11);
alter table server add CONSTRAINT `fk_server_cachegroup1` FOREIGN KEY (`cachegroup`) REFERENCES `cachegroup` (`id`) ON DELETE CASCADE;

alter table cran drop foreign key fk_cran_location1;
alter table cran drop key fk_cran_location1;
alter table cran change location cachegroup int(11);
alter table cran add CONSTRAINT `fk_cran_cachegroup1` FOREIGN KEY (`cachegroup`) REFERENCES `cachegroup` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION;

alter table staticdnsentry drop foreign key fk_staticdnsentry_location;
alter table staticdnsentry drop key fk_staticdnsentry_location;
alter table staticdnsentry change location cachegroup int(11);
alter table staticdnsentry add CONSTRAINT `fk_staticdnsentry_cachegroup1` FOREIGN KEY (`cachegroup`) REFERENCES `cachegroup` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION;

update type set use_in_table='cachegroup' where use_in_table='location';

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back