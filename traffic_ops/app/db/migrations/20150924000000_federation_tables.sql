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
-- federation_mapping
-- federation_resolver
CREATE TABLE `federation_resolver` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ip_address` VARCHAR(50) NOT NULL,
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `federation_resolver_ip_address` (`ip_address`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE IF NOT EXISTS `federation_mapping` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(45) NOT NULL,
  `description` VARCHAR(4096) NULL,
  `cname` VARCHAR(1024) NOT NULL,
  `ttl` INT(8) NOT NULL,
  `federation_resolver_id` INT(11) NOT NULL,
  `type` INT(11) NOT NULL,
  `last_updated` TIMESTAMP NOT NULL DEFAULT now(),
  PRIMARY KEY (`id`,`type`),
  KEY `fk_federation_resolver_id1` (`federation_resolver_id`),
  CONSTRAINT `fk_federation_resolver_id1` FOREIGN KEY (`federation_resolver_id`) REFERENCES `federation_resolver` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_federation_mapping_type` FOREIGN KEY (`type`) REFERENCES `type` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- federation_deliveryservice
CREATE TABLE `federation_deliveryservice` (
  `federation_mapping` int(11) NOT NULL,
  `deliveryservice` int(11) NOT NULL,
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`federation_mapping`,`deliveryservice`),
  KEY `fk_fed_to_ds1` (`deliveryservice`),
  CONSTRAINT `fk_fed_mapping_to_ds1` FOREIGN KEY (`deliveryservice`) REFERENCES `deliveryservice` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fed_mapping_to_fed1` FOREIGN KEY (`federation_mapping`) REFERENCES `federation_mapping` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
