package com.comcast.cdn.traffic_control.traffic_router.utils;

import org.apache.log4j.Logger;

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
