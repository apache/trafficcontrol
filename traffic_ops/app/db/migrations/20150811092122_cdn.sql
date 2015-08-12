
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE cdn (
	id 				int(11) 		NOT NULL AUTO_INCREMENT,
	name			varchar(1024) 	NOT NULL,
	last_updated 	timestamp 		NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
);

ALTER TABLE `server` ADD `cdn_id` int(11) DEFAULT NULL AFTER `profile`;
CREATE INDEX `fk_cdn2` ON `server`(`cdn_id`);
ALTER TABLE `server` ADD CONSTRAINT `fk_cdn2` FOREIGN KEY (`cdn_id`) REFERENCES `cdn` (`id`) ON DELETE SET NULL;

ALTER TABLE `deliveryservice` ADD `cdn_id` int(11) DEFAULT NULL AFTER `profile`;
CREATE INDEX `fk_cdn1` ON `deliveryservice`(`cdn_id`);
ALTER TABLE `deliveryservice` ADD CONSTRAINT `fk_cdn1` FOREIGN KEY (`cdn_id`) REFERENCES `cdn` (`id`) ON DELETE SET NULL;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `deliveryservice` DROP FOREIGN KEY `fk_cdn1`;
DROP INDEX `fk_cdn1` ON `deliveryservice`;
ALTER TABLE `deliveryservice` DROP `cdn_id`;

ALTER TABLE `server` DROP FOREIGN KEY `fk_cdn2`;
DROP INDEX `fk_cdn2` ON `server`;
ALTER TABLE `server` DROP `cdn_id`;

DROP TABLE cdn;
