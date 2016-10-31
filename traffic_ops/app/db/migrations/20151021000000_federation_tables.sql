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
-- federation
-- federation_resolver
CREATE TABLE `federation_resolver` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ip_address` VARCHAR(50) NOT NULL,
  `type` INT(11) NOT NULL,
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_federation_mapping_type` FOREIGN KEY (`type`) REFERENCES `type` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  UNIQUE KEY `federation_resolver_ip_address` (`ip_address`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE IF NOT EXISTS `federation` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `cname` VARCHAR(1024) NOT NULL,
  `description` VARCHAR(1024) NULL,
  `ttl` INT(8) NOT NULL,
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- federation_deliveryservice
CREATE TABLE `federation_deliveryservice` (
  `federation` int(11) NOT NULL,
  `deliveryservice` int(11) NOT NULL,
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`federation`,`deliveryservice`),
  KEY `fk_fed_to_ds1` (`deliveryservice`),
  CONSTRAINT `fk_federation_to_ds1` FOREIGN KEY (`deliveryservice`) REFERENCES `deliveryservice` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_federation_to_fed1` FOREIGN KEY (`federation`) REFERENCES `federation` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- federation_federation_resolver
CREATE TABLE `federation_federation_resolver` (
  `federation` int(11) NOT NULL,
  `federation_resolver` int(11) NOT NULL,
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`federation`,`federation_resolver`),
  KEY `fk_federation_federation_resolver` (`federation`),
  CONSTRAINT `fk_federation_federation_resolver1` FOREIGN KEY (`federation`) REFERENCES `federation` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_federation_resolver_to_fed1` FOREIGN KEY (`federation_resolver`) REFERENCES `federation_resolver` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- federation_tm_user
CREATE TABLE `federation_tmuser` (
  `federation` int(11) NOT NULL,
  `tm_user` int(11) NOT NULL,
  `role` INT(11) NOT NULL,
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`federation`,`tm_user`),
  KEY `fk_federation_federation_resolver` (`federation`),
  CONSTRAINT `fk_federation_tmuser_tmuser` FOREIGN KEY (`tm_user`) REFERENCES `tm_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_federation_tmuser_federation` FOREIGN KEY (`federation`) REFERENCES `federation` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_federation_tmuser_role` FOREIGN KEY (`role`) REFERENCES `role` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
