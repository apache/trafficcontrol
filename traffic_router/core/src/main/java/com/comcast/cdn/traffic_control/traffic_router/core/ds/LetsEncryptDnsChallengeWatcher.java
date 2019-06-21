package com.comcast.cdn.traffic_control.traffic_router.core.ds;

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

import com.comcast.cdn.traffic_control.traffic_router.core.config.ConfigHandler;
import com.comcast.cdn.traffic_control.traffic_router.core.util.AbstractResourceWatcher;
import com.comcast.cdn.traffic_control.traffic_router.core.util.JsonUtils;
import com.comcast.cdn.traffic_control.traffic_router.core.util.JsonUtilsException;
import com.fasterxml.jackson.core.JsonFactory;
import com.fasterxml.jackson.core.JsonParseException;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.node.ArrayNode;
import com.fasterxml.jackson.databind.node.ObjectNode;
import org.apache.log4j.Logger;

import java.io.*;
import java.time.Instant;
import java.util.HashMap;
import java.util.List;

public class LetsEncryptDnsChallengeWatcher extends AbstractResourceWatcher {
    private static final Logger LOGGER = Logger.getLogger(LetsEncryptDnsChallengeWatcher.class);
    public static final String DEFAULT_LE_DNS_CHALLENGE_URL = "https://${toHostname}/api/1.4/letsencrypt/dnsrecords/";
    private static final String configFile = "/opt/traffic_router/db/cr-config.json";

    private ConfigHandler configHandler;

    public void setConfigHandler(final ConfigHandler configHandler) {
        this.configHandler = configHandler;
    }

    public LetsEncryptDnsChallengeWatcher() {
        setDatabaseUrl(DEFAULT_LE_DNS_CHALLENGE_URL);
    }

    @Override
    public boolean useData(final String data) {
        try {
            final ObjectMapper mapper = new ObjectMapper(new JsonFactory());
            final HashMap<String, List<LetsEncryptDnsChallenge>> dataMap = mapper.readValue(data, new TypeReference<HashMap<String, List<LetsEncryptDnsChallenge>>>() { });
            final List<LetsEncryptDnsChallenge> challengeList = dataMap.get("response");

            JsonNode mostRecentConfig = mapper.readTree(readConfigFile());
            final ObjectNode deliveryServicesNode = (ObjectNode) JsonUtils.getJsonNode(mostRecentConfig, ConfigHandler.deliveryServicesKey);


            challengeList.forEach(challenge -> {
                final StringBuilder sb = new StringBuilder();
                sb.append(challenge.getFqdn());
                if (!challenge.getFqdn().endsWith(".")) {
                     sb.append('.');
                }
                final String challengeDomain = sb.toString();
                final String fqdn = challengeDomain.substring(0, challengeDomain.length() - 1).replace("_acme-challenge.", "");

                    final String dsLabel = fqdn.split("\\.")[0];
                    final ObjectNode deliveryServiceConfig = (ObjectNode) deliveryServicesNode.get(dsLabel);

                    ArrayNode staticDnsEntriesNode;

                    if (deliveryServiceConfig.findValue("staticDnsEntries") != null) {
                         staticDnsEntriesNode = (ArrayNode) deliveryServiceConfig.findValue("staticDnsEntries");
                    } else {
                        staticDnsEntriesNode = mapper.createArrayNode();
                    }

                    final ObjectNode newChildNode = mapper.createObjectNode();
                    newChildNode.put("type", "TXT");
                    newChildNode.put("name", "_acme-challenge");
                    newChildNode.put("value", challenge.getRecord());
                    newChildNode.put("ttl", 10);

                    staticDnsEntriesNode.add(newChildNode);

                    deliveryServiceConfig.set("staticDnsEntries", staticDnsEntriesNode);
                    deliveryServicesNode.set(dsLabel, deliveryServiceConfig);

            });

            final ObjectNode statsNode = (ObjectNode) mostRecentConfig.get("stats");
            statsNode.put("date", Instant.now().toEpochMilli() / 1000L);

//                    final ObjectNode fullConfig = mapper.createObjectNode();
            final ObjectNode fullConfig = (ObjectNode) mostRecentConfig;
            fullConfig.set(ConfigHandler.deliveryServicesKey, deliveryServicesNode);
            fullConfig.set("stats", statsNode);

            try {
                configHandler.processConfig(fullConfig.toString());
            } catch (JsonParseException | JsonUtilsException jsonError) {
                LOGGER.error("error processing config: " + jsonError.getMessage());
            }

            return true;
        } catch (Exception e) {
            LOGGER.warn("Failed updating dns challenge txt record with data from " + dataBaseURL);
        }

        return false;
    }

    @Override
    protected boolean verifyData(final String data) {
        try {
            final ObjectMapper mapper = new ObjectMapper(new JsonFactory());
            mapper.readValue(data, new TypeReference<HashMap<String, List<LetsEncryptDnsChallenge>>>() { });
            return true;
        } catch (Exception e) {
            LOGGER.warn("Failed to build dns challenge data while verifying");
        }

        return false;
    }

    @Override
    public String getWatcherConfigPrefix() {
        return "dnschallengemapping";
    }

    @SuppressWarnings("PMD.AppendCharacterWithChar")
    private String readConfigFile() {
        try {
            final InputStream is = new FileInputStream(configFile);
            final BufferedReader buf = new BufferedReader(new InputStreamReader(is));
            String line = buf.readLine();
            final StringBuilder sb = new StringBuilder();
            while (line != null) {
                sb.append(line).append("\n");
                line = buf.readLine();
            }
            return sb.toString();
        } catch (Exception e) {
            LOGGER.error("Could not read cr-config file.");
            return null;
        }
    }
}
