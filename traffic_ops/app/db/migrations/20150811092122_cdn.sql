/*
    Copyright 2015 Comcast Cable Communications Management, LLC

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
	config_file 	varchar(45) 	NOT NULL,
	last_updated 	timestamp 		NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
);

ALTER TABLE `deliveryservice` ADD `cdn_id` int(11) DEFAULT NULL AFTER `profile`;
CREATE INDEX `fk_cdn1` ON `deliveryservice`(`cdn_id`);
ALTER TABLE `deliveryservice` ADD CONSTRAINT `fk_cdn1` FOREIGN KEY (`cdn_id`) REFERENCES `cdn` (`id`) ON DELETE SET NULL;

ALTER TABLE `server` ADD `cdn_id` int(11) DEFAULT NULL AFTER `profile`;
CREATE INDEX `fk_cdn2` ON `server`(`cdn_id`);
ALTER TABLE `server` ADD CONSTRAINT `fk_cdn2` FOREIGN KEY (`cdn_id`) REFERENCES `cdn` (`id`) ON DELETE SET NULL;

ALTER TABLE `cachegroup` ADD `cdn_id` int(11) DEFAULT NULL AFTER `type`;
CREATE INDEX `fk_cdn3` ON `cachegroup`(`cdn_id`);
ALTER TABLE `cachegroup` ADD CONSTRAINT `fk_cdn3` FOREIGN KEY (`cdn_id`) REFERENCES `cdn` (`id`) ON DELETE SET NULL;

ALTER TABLE `profile` ADD `cdn_id` int(11) DEFAULT NULL AFTER `description`;
CREATE INDEX `fk_cdn4` ON `profile`(`cdn_id`);
ALTER TABLE `profile` ADD CONSTRAINT `fk_cdn4` FOREIGN KEY (`cdn_id`) REFERENCES `cdn` (`id`) ON DELETE SET NULL;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `deliveryservice` DROP FOREIGN KEY `fk_cdn1`;
DROP INDEX `fk_cdn1` ON `deliveryservice`;
ALTER TABLE `deliveryservice` DROP `cdn_id`;

ALTER TABLE `server` DROP FOREIGN KEY `fk_cdn2`;
DROP INDEX `fk_cdn2` ON `server`;
ALTER TABLE `server` DROP `cdn_id`;

ALTER TABLE `cachegroup` DROP FOREIGN KEY `fk_cdn3`;
DROP INDEX `fk_cdn3` ON `cachegroup`;
ALTER TABLE `cachegroup` DROP `cdn_id`;

ALTER TABLE `profile` DROP FOREIGN KEY `fk_cdn4`;
DROP INDEX `fk_cdn4` ON `profile`;
ALTER TABLE `profile` DROP `cdn_id`;

DROP TABLE cdn;
