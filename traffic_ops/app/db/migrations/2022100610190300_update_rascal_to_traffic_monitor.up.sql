/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with this
 * work for additional information regarding copyright ownership.  The ASF
 * licenses this file to you under the Apache License, Version 2.0 (the
 * 'License'); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an 'AS IS' BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  See the
 * License for the specific language governing permissions and limitations under
 * the License.
 */

UPDATE TYPE
SET name = 'TRAFFIC_MONITOR',
    description = 'Traffic Monitor polling & reporting'
WHERE name = 'RASCAL';

UPDATE PARAMETER
SET config_file = REPLACE(config_file, 'rascal', 'traffic_monitor'),
WHERE config_file = 'rascal-config.txt' OR config_file = 'rascal.properties';

UPDATE PARAMETER
SET value = 'TRAFFIC_MONITOR_TOP',
WHERE value = 'RASCAL_TOP' AND name = 'latest_traffic_monitor';

UPDATE PROFILE
SET description = REPLACE(description, 'Rascal', 'Traffic Monitor'), name= REPLACE(name, 'RASCAL', 'TRAFFIC_MONITOR')
WHERE type = 'TM_PROFILE';

