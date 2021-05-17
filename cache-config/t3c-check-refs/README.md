<!--
    Licensed to the Apache Software Foundation (ASF) under one
    or more contributor license agreements.  See the NOTICE file
    distributed with this work for additional information
    regarding copyright ownership.  The ASF licenses this file
    to you under the Apache License, Version 2.0 (the
    "License"); you may not use this file except in compliance
    with the License.  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing,
    software distributed under the License is distributed on an
    "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
    KIND, either express or implied.  See the License for the
    specific language governing permissions and limitations
    under the License.
-->

# t3c-check-refs

This implements the ATS plugin readiness verifier as defined in the
blueprint #4628, see https://github.com/apache/trafficcontrol/pull/4628

## Synopsis
  t3c-check-refs [options] [optional_config_file]

## Description
  The t3c-check-refs app will read an ATS formatted plugin.config or remap.config
  file line by line and verify that the plugin '.so' files are available in the
  filesystem or relative to the ATS plugin installation directory by the
  absolute or relative plugin filename.

  In addition, any plugin parameters that end in '.config', '.cfg', or '.txt'
  are considered to be plugin configuration files and there existence in the
  filesystem or relative to the ATS configuration files directory is verified.

  The configuration file argument is optional.  If no config file argument is
  supplied, t3c-check-refs reads its config file input from 'stdin'

## Options
  --log-location-debug=[value] | -d [value], where to log debugs, default is empty
  --log-location-error=[value], | -e [value], where to log errors, default is 'stderr'
  --log-location-info=[value] | -i [value], where to log infos, default is 'stderr'
  --trafficserver-config-dir=[value] | -c [value], where to find ATS config files, default is '/opt/trafficserver/etc/trafficserver'
  --trafficserver-plugin-dir=[value] | -p [value], where to find ATS plugins, default is '/opt/trafficserver/libexec/trafficserver'
  --help | -h, this help message

## Exit Status
  Returns 0 if no missing plugin DSO or config files are found.
  Otherwise the total number of missing plugin DSO and config files
  are returned.
  
  
