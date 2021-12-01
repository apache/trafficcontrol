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

package org.apache.traffic_control.traffic_router.utils;

import org.apache.logging.log4j.Logger;

import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.HashMap;
import java.util.Map;

public class HttpsProperties {
    private static final Logger log = Logger.getLogger(HttpsProperties.class);
    private static final String HTTPS_PROPERTIES_FILE = "/opt/traffic_router/conf/https.properties";
    private final Map<String, String> httpsPropertiesMap;

    public HttpsProperties() {
        this.httpsPropertiesMap = loadHttpsProperties();
    }

    public Map<String, String> getHttpsPropertiesMap() {
        return httpsPropertiesMap;
    }

    private static Map<String, String> loadHttpsProperties() {
        try {
            final Map<String, String> httpsProperties = new HashMap<>();
            Files.readAllLines(Paths.get(HTTPS_PROPERTIES_FILE)).forEach(propString -> {
                if (!propString.startsWith("#")) { // Ignores comments in properties file
                    final String[] prop = propString.split("=");
                    httpsProperties.put(prop[0], prop[1]);
                }
            });
            return httpsProperties;
        } catch (Exception e) {
            log.error("Error loading https properties file.");
            return null;
        }
    }
}
