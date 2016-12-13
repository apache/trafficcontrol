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
-- MySQL dump 10.13  Distrib 5.6.19, for osx10.9 (x86_64)
--
-- Host: localhost    Database: twelve_monkeys
-- ------------------------------------------------------
-- Server version	5.6.19

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `cran`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cran` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `asn` int(11) NOT NULL,
  `location` int(11) NOT NULL,
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`,`location`),
  UNIQUE KEY `cr_id_UNIQUE` (`id`),
  KEY `fk_cran_location1` (`location`),
  CONSTRAINT `fk_cran_location1` FOREIGN KEY (`location`) REFERENCES `location` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=47 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `deliveryservice`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `deliveryservice` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `xml_id` varchar(48) NOT NULL,
  `active` tinyint(4) NOT NULL,
  `dscp` int(11) NOT NULL,
  `signed` tinyint(1) DEFAULT NULL,
  `qstring_ignore` tinyint(1) DEFAULT NULL,
  `geo_limit` tinyint(1) DEFAULT '0',
  `http_bypass_fqdn` varchar(255) DEFAULT NULL,
  `dns_bypass_ip` varchar(45) DEFAULT NULL,
  `dns_bypass_ip6` varchar(45) DEFAULT NULL,
  `dns_bypass_ttl` int(11) DEFAULT NULL,
  `org_server_fqdn` varchar(255) DEFAULT NULL,
  `type` int(11) NOT NULL,
  `profile` int(11) NOT NULL,
  `ccr_dns_ttl` int(11) DEFAULT NULL,
  `global_max_mbps` int(11) DEFAULT NULL,
  `global_max_tps` int(11) DEFAULT NULL,
  `long_desc` varchar(255) DEFAULT NULL,
  `long_desc_1` varchar(255) DEFAULT NULL,
  `long_desc_2` varchar(255) DEFAULT NULL,
  `max_dns_answers` int(11) DEFAULT '0',
  `info_url` varchar(255) DEFAULT NULL,
  `miss_lat` double DEFAULT NULL,
  `miss_long` double DEFAULT NULL,
  `check_path` varchar(255) DEFAULT NULL,
  `header_rewrite` int(11) DEFAULT NULL,
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `protocol` tinyint(4) DEFAULT '0',
  PRIMARY KEY (`id`,`type`),
  UNIQUE KEY `ds_name_UNIQUE` (`xml_id`),
  UNIQUE KEY `ds_id_UNIQUE` (`id`),
  KEY `fk_deliveryservice_type1` (`type`),
  KEY `fk_deliveryservice_profile1` (`profile`),
  KEY `fk_deliveryservice_header_rewrite1_idx` (`header_rewrite`),
  CONSTRAINT `fk_deliveryservice_header_rewrite1` FOREIGN KEY (`header_rewrite`) REFERENCES `header_rewrite` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `fk_deliveryservice_profile1` FOREIGN KEY (`profile`) REFERENCES `profile` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `fk_deliveryservice_type1` FOREIGN KEY (`type`) REFERENCES `type` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=311 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `deliveryservice_regex`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `deliveryservice_regex` (
  `deliveryservice` int(11) NOT NULL,
  `regex` int(11) NOT NULL,
  `set_number` int(11) DEFAULT '0',
  PRIMARY KEY (`deliveryservice`,`regex`),
  KEY `fk_ds_to_regex_regex1` (`regex`),
  CONSTRAINT `fk_ds_to_regex_deliveryservice1` FOREIGN KEY (`deliveryservice`) REFERENCES `deliveryservice` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_ds_to_regex_regex1` FOREIGN KEY (`regex`) REFERENCES `regex` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `deliveryservice_server`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `deliveryservice_server` (
  `deliveryservice` int(11) NOT NULL,
  `server` int(11) NOT NULL,
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`deliveryservice`,`server`),
  KEY `fk_ds_to_cs_contentserver1` (`server`),
  CONSTRAINT `fk_ds_to_cs_contentserver1` FOREIGN KEY (`server`) REFERENCES `server` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_ds_to_cs_deliveryservice1` FOREIGN KEY (`deliveryservice`) REFERENCES `deliveryservice` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `deliveryservice_tmuser`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `deliveryservice_tmuser` (
  `deliveryservice` int(11) NOT NULL,
  `tm_user_id` int(11) NOT NULL,
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`deliveryservice`,`tm_user_id`),
  KEY `fk_tm_userid` (`tm_user_id`),
  CONSTRAINT `fk_tm_user_ds` FOREIGN KEY (`deliveryservice`) REFERENCES `deliveryservice` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_tm_user_id` FOREIGN KEY (`tm_user_id`) REFERENCES `tm_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `division`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `division` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(45) NOT NULL,
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name_UNIQUE` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `header_rewrite`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `header_rewrite` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `hr_condition` varchar(1024) DEFAULT NULL,
  `action` varchar(1024) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `hwinfo`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `hwinfo` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `serverid` int(11) NOT NULL,
  `description` varchar(256) NOT NULL,
  `val` varchar(256) NOT NULL,
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `serverid` (`serverid`,`description`),
  KEY `fk_hwinfo1` (`serverid`),
  CONSTRAINT `fk_hwinfo1` FOREIGN KEY (`serverid`) REFERENCES `server` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=4021555 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `job`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `job` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `agent` int(11) DEFAULT NULL,
  `object_type` varchar(48) DEFAULT NULL,
  `object_name` varchar(256) DEFAULT NULL,
  `keyword` varchar(48) NOT NULL,
  `parameters` varchar(256) DEFAULT NULL,
  `asset_url` varchar(512) NOT NULL,
  `asset_type` varchar(48) NOT NULL,
  `status` int(11) NOT NULL,
  `start_time` datetime NOT NULL,
  `entered_time` datetime NOT NULL,
  `job_user` int(11) NOT NULL,
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `fk_job_agent_id1` (`agent`),
  KEY `fk_job_status_id1` (`status`),
  KEY `fk_job_user_id1` (`job_user`),
  CONSTRAINT `fk_job_user_id1` FOREIGN KEY (`job_user`) REFERENCES `tm_user` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `fk_job_agent_id1` FOREIGN KEY (`agent`) REFERENCES `job_agent` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_job_status_id1` FOREIGN KEY (`status`) REFERENCES `job_status` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `job_agent`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `job_agent` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(128) DEFAULT NULL,
  `description` varchar(512) DEFAULT NULL,
  `active` int(1) NOT NULL DEFAULT '0',
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `job_result`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `job_result` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `job` int(11) NOT NULL,
  `agent` int(11) NOT NULL,
  `result` varchar(48) NOT NULL,
  `description` varchar(512) DEFAULT NULL,
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `fk_job_id1` (`job`),
  KEY `fk_agent_id1` (`agent`),
  CONSTRAINT `fk_agent_id1` FOREIGN KEY (`agent`) REFERENCES `job_agent` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_job_id1` FOREIGN KEY (`job`) REFERENCES `job` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `job_status`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `job_status` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(48) DEFAULT NULL,
  `description` varchar(256) DEFAULT NULL,
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `location`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `location` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(45) NOT NULL,
  `short_name` varchar(255) NOT NULL,
  `latitude` double DEFAULT NULL,
  `longitude` double DEFAULT NULL,
  `parent_location_id` int(11) DEFAULT NULL,
  `type` int(11) NOT NULL,
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`,`type`),
  UNIQUE KEY `loc_name_UNIQUE` (`name`),
  UNIQUE KEY `loc_short_UNIQUE` (`short_name`),
  UNIQUE KEY `lo_id_UNIQUE` (`id`),
  KEY `fk_location_type1` (`type`),
  KEY `fk_location_1` (`parent_location_id`),
  CONSTRAINT `fk_location_1` FOREIGN KEY (`parent_location_id`) REFERENCES `location` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `fk_location_type1` FOREIGN KEY (`type`) REFERENCES `type` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=52 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `location_parameter`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `location_parameter` (
  `location` int(11) NOT NULL,
  `parameter` int(11) NOT NULL,
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`location`,`parameter`),
  KEY `fk_location` (`location`),
  KEY `fk_parameter` (`parameter`),
  CONSTRAINT `fk_location` FOREIGN KEY (`location`) REFERENCES `location` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_parameter` FOREIGN KEY (`parameter`) REFERENCES `parameter` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `log`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `level` varchar(45) DEFAULT NULL,
  `message` varchar(1024) NOT NULL,
  `tm_user` int(11) NOT NULL,
  `ticketnum` varchar(64) DEFAULT NULL,
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`,`tm_user`),
  KEY `fk_log_1` (`tm_user`),
  CONSTRAINT `fk_log_1` FOREIGN KEY (`tm_user`) REFERENCES `tm_user` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=21879 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `parameter`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `parameter` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(1024) NOT NULL,
  `config_file` varchar(45) NOT NULL,
  `value` varchar(1024) NOT NULL,
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=817 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `phys_location`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `phys_location` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(45) NOT NULL,
  `short_name` varchar(12) NOT NULL,
  `address` varchar(128) NOT NULL,
  `city` varchar(128) NOT NULL,
  `state` varchar(2) NOT NULL,
  `zip` varchar(5) NOT NULL,
  `poc` varchar(128) DEFAULT NULL,
  `phone` varchar(45) DEFAULT NULL,
  `email` varchar(128) DEFAULT NULL,
  `comments` varchar(256) DEFAULT NULL,
  `region` int(11) NOT NULL,
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name_UNIQUE` (`name`),
  UNIQUE KEY `short_name_UNIQUE` (`short_name`),
  KEY `fk_phys_location_region_idx` (`region`),
  CONSTRAINT `fk_phys_location_region` FOREIGN KEY (`region`) REFERENCES `region` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=200 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `profile`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `profile` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(45) NOT NULL,
  `description` varchar(256) DEFAULT NULL,
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name_UNIQUE` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=48 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `profile_parameter`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `profile_parameter` (
  `profile` int(11) NOT NULL,
  `parameter` int(11) NOT NULL,
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`profile`,`parameter`),
  KEY `fk_atsprofile_atsparameters_atsprofile1` (`profile`),
  KEY `fk_atsprofile_atsparameters_atsparameters1` (`parameter`),
  CONSTRAINT `fk_atsprofile_atsparameters_atsparameters1` FOREIGN KEY (`parameter`) REFERENCES `parameter` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_atsprofile_atsparameters_atsprofile1` FOREIGN KEY (`profile`) REFERENCES `profile` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `regex`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `regex` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pattern` varchar(255) NOT NULL DEFAULT '',
  `type` int(11) NOT NULL,
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`,`type`),
  UNIQUE KEY `re_id_UNIQUE` (`id`),
  KEY `fk_regex_type1` (`type`),
  CONSTRAINT `fk_regex_type1` FOREIGN KEY (`type`) REFERENCES `type` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=519 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `region`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `region` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(45) NOT NULL,
  `division` int(11) NOT NULL,
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name_UNIQUE` (`name`),
  KEY `fk_region_division1_idx` (`division`),
  CONSTRAINT `fk_region_division1` FOREIGN KEY (`division`) REFERENCES `division` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `role`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `role` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(45) NOT NULL,
  `description` varchar(128) DEFAULT NULL,
  `priv_level` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `server`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `server` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `host_name` varchar(45) NOT NULL,
  `domain_name` varchar(45) NOT NULL,
  `tcp_port` int(10) unsigned DEFAULT NULL,
  `xmpp_id` varchar(256) DEFAULT NULL,
  `xmpp_passwd` varchar(45) DEFAULT NULL,
  `interface_name` varchar(45) NOT NULL,
  `ip_address` varchar(45) NOT NULL,
  `ip_netmask` varchar(45) NOT NULL,
  `ip_gateway` varchar(45) NOT NULL,
  `ip6_address` varchar(50) DEFAULT NULL,
  `ip6_gateway` varchar(50) DEFAULT NULL,
  `interface_mtu` int(11) NOT NULL DEFAULT '9000',
  `phys_location` int(11) NOT NULL,
  `rack` varchar(64) DEFAULT NULL,
  `location` int(11) NOT NULL,
  `type` int(11) NOT NULL,
  `status` int(11) NOT NULL,
  `profile` int(11) NOT NULL,
  `mgmt_ip_address` varchar(45) DEFAULT NULL,
  `mgmt_ip_netmask` varchar(45) DEFAULT NULL,
  `mgmt_ip_gateway` varchar(45) DEFAULT NULL,
  `ilo_ip_address` varchar(45) DEFAULT NULL,
  `ilo_ip_netmask` varchar(45) DEFAULT NULL,
  `ilo_ip_gateway` varchar(45) DEFAULT NULL,
  `ilo_username` varchar(45) DEFAULT NULL,
  `ilo_password` varchar(45) DEFAULT NULL,
  `router_host_name` varchar(256) DEFAULT NULL,
  `router_port_name` varchar(256) DEFAULT NULL,
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`,`location`,`type`,`status`,`profile`),
  UNIQUE KEY `cs_ip_address_UNIQUE` (`ip_address`),
  UNIQUE KEY `se_id_UNIQUE` (`id`),
  UNIQUE KEY `host_name` (`host_name`),
  UNIQUE KEY `ip6_address` (`ip6_address`),
  KEY `fk_contentserver_location` (`location`),
  KEY `fk_contentserver_contentservertype1` (`type`),
  KEY `fk_contentserver_contentserverstatus1` (`status`),
  KEY `fk_contentserver_atsprofile1` (`profile`),
  KEY `fk_contentserver_phys_location1` (`phys_location`),
  CONSTRAINT `fk_contentserver_atsprofile1` FOREIGN KEY (`profile`) REFERENCES `profile` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `fk_contentserver_contentserverstatus1` FOREIGN KEY (`status`) REFERENCES `status` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `fk_contentserver_contentservertype1` FOREIGN KEY (`type`) REFERENCES `type` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `fk_contentserver_location` FOREIGN KEY (`location`) REFERENCES `location` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_contentserver_phys_location1` FOREIGN KEY (`phys_location`) REFERENCES `phys_location` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=580 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `serverstatus`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `serverstatus` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ilo_pingable` tinyint(1) NOT NULL DEFAULT '0',
  `teng_pingable` tinyint(1) NOT NULL DEFAULT '0',
  `fqdn_pingable` tinyint(1) DEFAULT '0',
  `dscp` tinyint(1) DEFAULT NULL,
  `firmware` tinyint(1) DEFAULT NULL,
  `marvin` tinyint(1) DEFAULT NULL,
  `ping6` tinyint(1) DEFAULT NULL,
  `upd_pending` tinyint(1) DEFAULT NULL,
  `stats` tinyint(1) DEFAULT NULL,
  `prox` tinyint(1) DEFAULT NULL,
  `mtu` tinyint(1) DEFAULT NULL,
  `ccr_online` tinyint(1) DEFAULT NULL,
  `rascal` tinyint(1) DEFAULT NULL,
  `chr` int(11) DEFAULT NULL,
  `cdu` int(11) DEFAULT NULL,
  `ort_errors` int(11) NOT NULL DEFAULT '-1',
  `mbps_out` int(11) DEFAULT '0',
  `clients_connected` int(11) DEFAULT '0',
  `server` int(11) NOT NULL,
  `last_recycle_date` timestamp NULL DEFAULT NULL,
  `last_recycle_duration_hrs` int(11) DEFAULT '0',
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`,`server`),
  UNIQUE KEY `server` (`server`),
  UNIQUE KEY `ses_id_UNIQUE` (`id`),
  KEY `fk_serverstatus_server1` (`server`),
  CONSTRAINT `fk_serverstatus_server1` FOREIGN KEY (`server`) REFERENCES `server` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=4180784 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `staticdnsentry`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `staticdnsentry` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `host` varchar(45) NOT NULL,
  `address` varchar(45) NOT NULL,
  `type` int(11) NOT NULL,
  `ttl` int(11) NOT NULL DEFAULT '3600',
  `deliveryservice` int(11) NOT NULL,
  `location` int(11) NOT NULL,
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `combi_UNIQUE` (`host`,`address`,`deliveryservice`,`location`),
  KEY `fk_staticdnsentry_type` (`type`),
  KEY `fk_staticdnsentry_ds` (`deliveryservice`),
  KEY `fk_staticdnsentry_location` (`location`),
  CONSTRAINT `fk_staticdnsentry_ds` FOREIGN KEY (`deliveryservice`) REFERENCES `deliveryservice` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `fk_staticdnsentry_location` FOREIGN KEY (`location`) REFERENCES `location` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `fk_staticdnsentry_type` FOREIGN KEY (`type`) REFERENCES `type` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `status`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `status` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(45) NOT NULL,
  `description` varchar(256) DEFAULT NULL,
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `tm_user`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tm_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(128) DEFAULT NULL,
  `role` int(11) DEFAULT NULL,
  `uid` int(11) DEFAULT NULL,
  `gid` int(11) DEFAULT NULL,
  `local_passwd` varchar(40) DEFAULT NULL,
  `confirm_local_passwd` varchar(40) DEFAULT NULL,
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `company` varchar(256) DEFAULT NULL,
  `email` varchar(128) DEFAULT NULL,
  `full_name` varchar(256) DEFAULT NULL,
  `new_user` tinyint(1) NOT NULL DEFAULT '1',
  `address_line1` varchar(256) DEFAULT NULL,
  `address_line2` varchar(256) DEFAULT NULL,
  `city` varchar(128) DEFAULT NULL,
  `state_or_province` varchar(128) DEFAULT NULL,
  `phone_number` varchar(25) DEFAULT NULL,
  `postal_code` varchar(11) DEFAULT NULL,
  `country` varchar(256) DEFAULT NULL,
  `local_user` tinyint(1) NOT NULL DEFAULT '0',
  `token` varchar(50) DEFAULT NULL,
  `registration_sent` timestamp NOT NULL DEFAULT '1999-01-01 00:00:00',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username_UNIQUE` (`username`),
  KEY `fk_user_1` (`role`),
  CONSTRAINT `fk_user_1` FOREIGN KEY (`role`) REFERENCES `role` (`id`) ON DELETE SET NULL ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=57 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `type`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `type` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(45) NOT NULL,
  `description` varchar(45) NOT NULL,
  `use_in_table` varchar(45) DEFAULT NULL,
  `last_updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `NAME_UNIQUE` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=31 DEFAULT CHARSET=latin1;

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

/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2015-01-14 14:13:34
