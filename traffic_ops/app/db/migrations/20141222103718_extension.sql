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

update serverstatus set last_recycle_date = NULL;

-- migrate the current table
alter table serverstatus 
  CHANGE `ilo_pingable` `aa` INT(11)  NULL DEFAULT NULL,
  CHANGE `teng_pingable` `ab` INT(11)  NULL DEFAULT NULL,
  CHANGE `fqdn_pingable` `ac` INT(11) NULL DEFAULT NULL,
  CHANGE `dscp` `ad` INT(11) NULL DEFAULT NULL,
  CHANGE `firmware` `ae` INT(11) NULL DEFAULT NULL,
  CHANGE `marvin` `af` INT(11) NULL DEFAULT NULL,
  CHANGE `ping6` `ag` INT(11) NULL DEFAULT NULL,
  CHANGE `upd_pending` `ah` INT(11) NULL DEFAULT NULL,
  CHANGE `stats` `ai` INT(11) NULL DEFAULT NULL,
  CHANGE `prox` `aj` INT(11) NULL DEFAULT NULL,
  CHANGE `mtu` `ak` INT(11) NULL DEFAULT NULL,
  CHANGE `ccr_online` `al` INT(11) NULL DEFAULT NULL,
  CHANGE `rascal` `am` INT(11) NULL DEFAULT NULL,
  CHANGE `chr` `an` INT(11) NULL DEFAULT NULL,
  CHANGE `cdu` `ao` INT(11) NULL DEFAULT NULL,
  CHANGE `ort_errors` `ap` INT(11) NULL DEFAULT NULL,
  CHANGE `mbps_out` `aq` INT(11) NULL DEFAULT NULL,
  CHANGE `clients_connected` `ar` INT(11) NULL DEFAULT NULL,
  CHANGE `last_recycle_date` `as` INT(11) NULL DEFAULT NULL,
  CHANGE `last_recycle_duration_hrs` `at` INT(11) NULL DEFAULT NULL;

alter table serverstatus modify `server` INT(11) NOT NULL AFTER `id`;
alter table serverstatus add column `au` INT(11) NULL DEFAULT NULL after `at`;
alter table serverstatus add column `av` INT(11) NULL DEFAULT NULL after `au`;
alter table serverstatus add column `aw` INT(11) NULL DEFAULT NULL after `av`;
alter table serverstatus add column `ax` INT(11) NULL DEFAULT NULL after `aw`;
alter table serverstatus add column `ay` INT(11) NULL DEFAULT NULL after `ax`;
alter table serverstatus add column `az` INT(11) NULL DEFAULT NULL after `ay`;
alter table serverstatus add column `ba` INT(11) NULL DEFAULT NULL after `az`;
alter table serverstatus add column `bb` INT(11) NULL DEFAULT NULL after `ba`;
alter table serverstatus add column `bc` INT(11) NULL DEFAULT NULL after `bb`;
alter table serverstatus add column `bd` INT(11) NULL DEFAULT NULL after `bc`;
alter table serverstatus add column `be` INT(11) NULL DEFAULT NULL after `bd`;

-- there shouldn't be any updates pending while doing the TM update, and upd_pending gets moved to the server table.
alter table server add column `upd_pending` TINYINT(1) NOT NULL DEFAULT 0 after status; 

rename table serverstatus to servercheck;


CREATE TABLE IF NOT EXISTS `to_extension` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(45) NOT NULL,
  `version` VARCHAR(45) NOT NULL,
  `info_url` VARCHAR(45) NOT NULL,
  `script_file` VARCHAR(45) NOT NULL,
  `isactive` TINYINT(1) NOT NULL,
  `additional_config_json` VARCHAR(4096) NULL,
  `description` VARCHAR(4096) NULL,
  `servercheck_short_name` VARCHAR(8) NULL,
  `servercheck_column_name` VARCHAR(10) NULL,
  `type` INT(11) NOT NULL,
  `last_updated` TIMESTAMP NOT NULL DEFAULT now(),
  PRIMARY KEY (`id`),
  UNIQUE INDEX `id_UNIQUE` (`id` ASC),
  INDEX `fk_ext_type_idx` (`type` ASC),
  CONSTRAINT `fk_ext_type`
    FOREIGN KEY (`type`)
    REFERENCES `type` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB 
DEFAULT CHARACTER SET = latin1;


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

-- drop table to_extension;
-- delete from type where name like ('%XTENSION');
-- delete from tm_user where username='extension';
-- drop table dynserverstatus;
-- drop table dynserverstatusentry;
