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

package org.apache.traffic_control.traffic_router.core.router;

import com.fasterxml.jackson.databind.JsonNode;

import java.util.ArrayList;
import java.util.List;

// Attempts to generate names like 'www.[foo].kabletown.com' to do dns queries against traffic router
// Tries to pull 'whole' words from the regex of cr-config
public class DnsNameGenerator {
    public List<String> getNames(JsonNode deliveryServicesConfig, JsonNode cdnConfig) throws Exception {
        List<String> names = new ArrayList<String>();

        String domainName = cdnConfig.get("domain_name").asText();

        for (final JsonNode matchsets : deliveryServicesConfig.get("matchsets")) {
            for (final JsonNode matchset : matchsets) {
                if (!"DNS".equals(matchset.get("protocol").asText())) {
                    continue;
                }

                for (final JsonNode matchlist : matchset.get("matchlist")) {
                    final String name = matchlist.get("regex").asText()
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
