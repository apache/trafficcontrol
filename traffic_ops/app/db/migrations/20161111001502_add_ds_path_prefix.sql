/*
	Copyright 2016 Cisco Systems, Inc.

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
CREATE TABLE `deliveryservice_path_prefix` (
  `deliveryservice` int(11) NOT NULL,
  `path_prefix` varchar(255) NOT NULL,
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`deliveryservice`,`path_prefix`),
  CONSTRAINT `fk_ds_to_path_prefix_deliveryservice1` FOREIGN KEY (`deliveryservice`) REFERENCES `deliveryservice` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB DEFAULT CHARACTER SET = latin1;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE deliveryservice_path_prefix;
