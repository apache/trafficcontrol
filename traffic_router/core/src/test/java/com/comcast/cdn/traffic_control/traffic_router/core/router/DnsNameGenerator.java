/*
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package com.comcast.cdn.traffic_control.traffic_router.core.router;

import org.json.JSONArray;
import org.json.JSONObject;

import java.util.ArrayList;
import java.util.List;

// Attempts to generate names like 'www.[foo].kabletown.com' to do dns queries against traffic router
// Tries to pull 'whole' words from the regex of cr-config
public class DnsNameGenerator {
    public List<String> getNames(JSONObject deliveryServicesConfig, JSONObject cdnConfig) throws Exception {
        List<String> names = new ArrayList<String>();

        String domainName = cdnConfig.getString("domain_name");

        for (String deliveryServiceId : JSONObject.getNames(deliveryServicesConfig)) {
            final JSONArray matchsets = deliveryServicesConfig.getJSONObject(deliveryServiceId).getJSONArray("matchsets");

            for (int i = 0; i < matchsets.length(); i++) {
                final JSONObject matchset = matchsets.getJSONObject(i);

                if (!"DNS".equals(matchset.getString("protocol"))) {
                    continue;
                }

                final JSONArray list = matchset.getJSONArray("matchlist");
                for (int j = 0; j < list.length(); j++) {
                    // Not bulletproof
                    final String name = list.getJSONObject(j).getString("regex")
                        .replaceAll("\\.", "")
                        .replaceAll("\\*", "")
                        .replaceAll("\\\\", "");

                    names.add("edge." + name + "." + domainName);
                }
            }
        }

        return names;
    }
}
