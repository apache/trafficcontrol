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

CREATE TABLE cdn (
	id 				int(11) 		NOT NULL AUTO_INCREMENT,
	name			varchar(1024) 	NOT NULL,
	last_updated 	timestamp 		NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
) ENGINE = InnoDB DEFAULT CHARACTER SET = latin1;

ALTER TABLE `deliveryservice` ADD `cdn_id` int(11) DEFAULT NULL AFTER `profile`;
CREATE INDEX `fk_cdn1` ON `deliveryservice`(`cdn_id`);
ALTER TABLE `deliveryservice` ADD CONSTRAINT `fk_cdn1` FOREIGN KEY (`cdn_id`) REFERENCES `cdn` (`id`) ON DELETE SET NULL;

ALTER TABLE `server` ADD `cdn_id` int(11) DEFAULT NULL AFTER `profile`;
CREATE INDEX `fk_cdn2` ON `server`(`cdn_id`);
ALTER TABLE `server` ADD CONSTRAINT `fk_cdn2` FOREIGN KEY (`cdn_id`) REFERENCES `cdn` (`id`) ON DELETE SET NULL;

INSERT INTO cdn(name) (
  SELECT parameter.value
  FROM parameter
  WHERE parameter.name = 'CDN_name'
);

update deliveryservice ds
set ds.cdn_id = ( 
  select cdn.id
  from profile p, profile_parameter pp, parameter param, cdn
  where ds.profile = p.id and pp.profile = p.id and pp.parameter = param.id
  and param.name = 'CDN_name'
  and cdn.name = param.value
);

update server s
set s.cdn_id = ( 
  select cdn.id
  from profile p, profile_parameter pp, parameter param, cdn
  where s.profile = p.id and pp.profile = p.id and pp.parameter = param.id
  and param.name = 'CDN_name'
  and cdn.name = param.value
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `deliveryservice` DROP FOREIGN KEY `fk_cdn1`;
DROP INDEX `fk_cdn1` ON `deliveryservice`;
ALTER TABLE `deliveryservice` DROP `cdn_id`;

ALTER TABLE `server` DROP FOREIGN KEY `fk_cdn2`;
DROP INDEX `fk_cdn2` ON `server`;
ALTER TABLE `server` DROP `cdn_id`;

DROP TABLE cdn;
