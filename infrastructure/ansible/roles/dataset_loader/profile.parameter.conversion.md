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
# Conversion of existing profile paramters to yaml dictionary

## Non-ATS
1. Install `jq`
2. Install `ruby`
3. Curl to get the response of https://{{ TO_BASE_URL }}/api/{{ dl_to_api_version }}/profiles/name/{{ TARGET_PROFILE }}/parameters and save that to a file named `profile.parameters.json`
4. `jq -r '[.[] | sort_by(.configFile,.name,.value)[] | {name: .name,configFile: .configFile,value: .value,secure: (if .secure then 1 else 0 end)}]' profile.parameters.json > profile.parameters.filtered.json`
5. `ruby -ryaml -rjson -e 'puts YAML.dump(JSON.parse(STDIN.read))' < profile.parameters.filtered.json > profile.parameters.yml`
6. Transplant into appropriate ansible var file

## ATS
ATS profiles are both huge and overlap eachother, even with only one edge+mid profile per ATS version.  As a result they were pulled out into a separate normalized structure with redundant parameters removed for a ~87% reduction.

In order to update them, it's important to know the scope of the change you're making such as all instances of ATS vs only edges running ATS7.  This hopefully helps with maintaining modern recommendations backward into older ATS versions as applicable.

### analyze.ats.profiles.yml
This is a helper ansible playbook which takes two existing traffic_ops profiles and provides information on changes and redundancy of parameters.

#### Usage
```shell
ansible-playbook -i localhost analyze.ats.profiles.yml -e "dl_to_url=$TO_URL" -e "dl_to_user=$TO_USER" -e "dl_to_user_password=$TO_PASSWORD" -e "TARGET_PROFILE=MID_ATS_8" -e "BASE_PROFILE=MID_ATS_7"
```
##### Required Parameters:
* dl_to_url : TrafficOps URL (http://to.kabletown.invalid)
* dl_to_user : TrafficOps User (user)
* dl_to_password : TrafficOps Password (pass)
* BASE_PROFILE : TrafficOps original base profile name (MID_ATS_7)
* TARGET_PROFILE : TrafficOps target new profile name (MID_ATS_8)

##### Optional Parameters:
* override_target_host_fqdn : Bypass the selection of a target host to obtain defaults by profile association
* override_base_host_fqdn : Bypass the selection of a base host to obtain defaults by profile association
* output_path : Alternate path to send output files.  Defaults to playbook directory.
* exclution_file_path : A path to a variable file containing exclusions to consider when comparing profiles

##### Additional Requirements:
* The target and base profiles must exist in TrafficOps
* Either:
    * A host override optional parameter is supplied
    * A host must be attached and reachable to each profile
* Each specified host must have ATS installed and running

#### Parameter exclusions
Included is a default set of parameter names/configFiles that are excluded from comparison with the dataset loader and TO diff output.  This is because there are some parameters which naturally vary by hardware, CDN, or delivery service assignments.  It does use a simple string contains logic to know about exclusions, so generally being more specific is appropriate if adjusting them.

#### Output
There are two output files from this playbook:

1. ats_profile_diff_output.yml

Simply a denormalized rendering of the normalized target profile parameters.  Useful if you'd like to still use the above non-ats method for comparison as a double verification.

2. ats_profile_diff_output.json

The primary output to help process key information.  There are several important bits to note in the structure:

##### Sample output for an ATS7 vs ATS8 edge
```json
{
  "ats": {
    "ats_default_differences_between_base_and_target_between_base_and_target": [
      {
        "base_value": "CONFIG proxy.config.http.request_via_str STRING ApacheTrafficServer/7.1.4",
        "target_value": "CONFIG proxy.config.http.request_via_str STRING ApacheTrafficServer/8.1.2"
      }
    ],
    "ats_defaults_missing_from_base_in_target": [
      "CONFIG proxy.config.http.negative_caching_list STRING 204 305 403 404 414 500 501 502 503 504"
    ],
    "ats_defaults_missing_from_target_in_base": [
      "CONFIG proxy.config.admin.synthetic_port INT 8083"
    ]
  },
  "base_ats_version": "trafficserver-7.1.4_rc0-26.b5fbc9a.el7.x86_64",
  "base_host": "cdn-ec-ats7.kabletown.invalid",
  "base_profile": "EDGE_ATS7",
  "dataset_loader": {
    "param_value_differences_between_dataset_loader_and_base": [
      {
        "configFile": "records.config",
        "name": "CONFIG proxy.config.http.connect_attempts_max_retries",
        "value": {
          "dataset_loader": "INT 6",
          "traffic_ops": "INT 3"
        }
      }
    ],
    "param_value_differences_between_dataset_loader_and_target": [
      {
        "configFile": "CRConfig.json",
        "name": "weight",
        "value": {
          "dataset_loader": "1.0",
          "traffic_ops": "1.2"
        }
      }
    ],
    "params_missing_from_base_in_dataset_loader": [],
    "params_missing_from_dataset_loader_in_base": [
      {
        "configFile": "records.config",
        "name": "CONFIG proxy.config.http.wait_for_cache",
        "value": "INT 2"
      }
    ],
    "params_missing_from_dataset_loader_in_target": [],
    "params_missing_from_target_in_dataset_loader": [],
    "redundant_parameters_in_both_base_and_target_in_dataset_loader": [],
    "redundant_parameters_with_base_params_and_defaults_included_in_dataset_loader_only": [],
    "redundant_parameters_with_target_params_and_defaults_included_in_dataset_loader_only": []
  },
  "target_ats_version": "trafficserver-8.1.1-21.402826e.el7.x86_64",
  "target_host": "cdn-ec-ats8.kabletown.invalid",
  "target_profile": "EDGE_ATS8",
  "traffic_ops": {
    "redundant_parameters_in_both_base_and_target_in_traffic_ops": [
      {
        "configFile": "records.config",
        "name": "CONFIG proxy.config.admin.admin_user",
        "value": "STRING admin"
      }
    ],
    "redundant_parameters_with_base_params_and_defaults_included_in_traffic_ops_only": [
      {
        "configFile": "records.config",
        "name": "CONFIG proxy.config.admin.synthetic_port",
        "value": "INT 8083"
      }
    ],
    "redundant_parameters_with_target_params_and_defaults_included_in_traffic_ops_only": [
      {
        "configFile": "records.config",
        "name": "CONFIG proxy.config.log.config.filename",
        "value": "STRING logging.yaml"
      }
    ]
  }
}
```

###### Basic info
This is included so looking at the outputs later have enough information to explain where it came from and was comparing.

* base_ats_version : ATS RPM version as determined by probing an associated or overridden directly
* base_host : The host selected for use either by override or just being the first to respond that is attached to the base profile
* base_profile : The original (older) ATS version
* target_ats_version : ATS RPM version as determined by probing an associated or overridden directly
* target_host : The host selected for use either by override or just being the first to respond that is attached to the target profile
* target_profile : The target (newer) ATS version

###### ATS info
This data has no concept of Traffic Ops and is purely based on the defaults as reported by traffic_ctl.

* ats_default_differences_between_base_and_target_between_base_and_target: Defautls present in both with differing values
* ats_defaults_missing_from_base_in_target : New parameters added to the target
* ats_defaults_missing_from_target_in_base: Parameters removed from the target

###### TrafficOps info
This uses the ATS default information in both target and base along with the TrafficOps profiles to identify redundant unnecessary parameters.

* redundant_parameters_in_both_base_and_target_in_traffic_ops
* redundant_parameters_with_base_params_and_defaults_included_in_traffic_ops_only
* redundant_parameters_with_target_params_and_defaults_included_in_traffic_ops_only

###### Dataset Loader info
This compares data from all sources to help provide guidance on how to better maintain the normalized redundancy-free dataset loader defaults.

1. Parameters related to value differences in the dataset loader going both directions
   * param_value_differences_between_dataset_loader_and_base
   * param_value_differences_between_dataset_loader_and_target
2. Parameters dealing with things added or removed from the base version and dataset loader
   * params_missing_from_base_in_dataset_loader
   * params_missing_from_dataset_loader_in_base
3. Parameters dealing with things added or removed from the target version and dataset loader
   * params_missing_from_dataset_loader_in_target
   * params_missing_from_target_in_dataset_loader
4. Parameters Identifying redundant parameters in the dataset loader
   * redundant_parameters_in_both_base_and_target_in_dataset_loader
   * redundant_parameters_with_base_params_and_defaults_included_in_dataset_loader_only
   * redundant_parameters_with_target_params_and_defaults_included_in_dataset_loader_only

##### Interpreting output
1. dataset_loader.redundant* parameters should generally always be empty.  If not, a default has probably changed and made a previous value newly redundant
2. dataset_loader.* these should be weighed carefully as they indicate changes that might be worth modeling as a default update
3. trafficOps.* parameters are mostly there informationally so that you can cleanup TrafficOps if you like
4. ats.* are there so that engineers routinely working with ATS have a concise list of changes between versions without influence from ATC.  It's likely these additions/removals may lead to changes in the dataset loader defaults as appropriate.
5. When updating the dataset loader defaults for ATS it's important to think about the scope of changes.  This is so that best practice updates are correctly backported to older versions of ATS as well, but may diverge by ATS version or tier if needed.
