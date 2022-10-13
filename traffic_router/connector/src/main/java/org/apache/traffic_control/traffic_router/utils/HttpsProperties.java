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

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.HashMap;
import java.util.Map;

public class HttpsProperties {
    private static final Logger log = LogManager.getLogger(HttpsProperties.class);
    private static final String HTTPS_CERTIFICATE_OU = "https.certificate.organizational.unit";

    private final Map<String, String> httpsPropertiesMap;

    public HttpsProperties(final String fileName) {
        this.httpsPropertiesMap = loadHttpsProperties(fileName);
    }

    public Map<String, String> getHttpsPropertiesMap() {
        return httpsPropertiesMap;
    }

    private static Map<String, String> loadHttpsProperties(final String fileName) {
        try {
            final Map<String, String> httpsProperties = new HashMap<>();
            Files.readAllLines(Paths.get(fileName)).forEach(propString -> {
                if (!propString.startsWith("#")) { // Ignores comments in properties file
                    final String[] props = propString.split("=");
                    if (props.length < 2) {
                        log.error("Property malformed, should be in the form key=value");
                    } else {
                        final String key = props[0];
                        final String val = props[1];
                        if (key.equals(HTTPS_CERTIFICATE_OU)) {
                            if (val.equals("") || val.length() < 2) {
                                log.error("Malformed " + HTTPS_CERTIFICATE_OU + " property value");
                            } else {
                                final String[] orgUnits = val.split(",");
                                String organizationalUnit = "";
                                StringBuilder sb = new StringBuilder(organizationalUnit);
                                for (final String ou : orgUnits) {
                                    sb = sb.append("; OU=" + ou);
                                }
                                organizationalUnit = sb.toString();
                                httpsProperties.put(key, organizationalUnit);
                            }
                        } else {
                            httpsProperties.put(key, val);
                        }
                    }
                }
            });
            return httpsProperties;
        } catch (Exception e) {
            log.error("Error loading https properties file at "+ fileName+ ", error: " +e.getMessage());
            return null;
        }
    }
}
