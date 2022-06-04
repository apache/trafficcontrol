package org.apache.traffic_control.traffic_router.core.ds;

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

import org.apache.traffic_control.traffic_router.core.config.ConfigHandler;
import org.apache.traffic_control.traffic_router.core.util.AbstractResourceWatcher;
import org.apache.traffic_control.traffic_router.core.util.JsonUtils;
import org.apache.traffic_control.traffic_router.core.util.JsonUtilsException;
import org.apache.traffic_control.traffic_router.core.util.TrafficOpsUtils;
import com.fasterxml.jackson.core.JsonFactory;
import com.fasterxml.jackson.core.JsonParseException;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.node.ArrayNode;
import com.fasterxml.jackson.databind.node.ObjectNode;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import java.io.*;
import java.time.Instant;
import java.util.HashMap;
import java.util.Iterator;
import java.util.List;

public class LetsEncryptDnsChallengeWatcher extends AbstractResourceWatcher {
    private static final Logger LOGGER = LogManager.getLogger(LetsEncryptDnsChallengeWatcher.class);
    public static final String DEFAULT_LE_DNS_CHALLENGE_URL = "https://${toHostname}/api/"+TrafficOpsUtils.TO_API_VERSION+"/letsencrypt/dnsrecords/";

    private String configFile;
    private ConfigHandler configHandler;

    public LetsEncryptDnsChallengeWatcher() {
        setDatabaseUrl(DEFAULT_LE_DNS_CHALLENGE_URL);
        setDefaultDatabaseUrl(DEFAULT_LE_DNS_CHALLENGE_URL);
    }

    @Override
    public boolean useData(final String data) {
        try {
            final ObjectMapper mapper = new ObjectMapper(new JsonFactory());
            final HashMap<String, List<LetsEncryptDnsChallenge>> dataMap = mapper.readValue(data, new TypeReference<HashMap<String, List<LetsEncryptDnsChallenge>>>() { });
            final List<LetsEncryptDnsChallenge> challengeList = dataMap.get("response");

            final JsonNode mostRecentConfig = mapper.readTree(readConfigFile());
            final ObjectNode deliveryServicesNode = (ObjectNode) JsonUtils.getJsonNode(mostRecentConfig, ConfigHandler.deliveryServicesKey);


            challengeList.forEach(challenge -> {
                final ObjectNode deliveryServiceConfig = (ObjectNode) deliveryServicesNode.get(challenge.getXmlId());
                if (deliveryServiceConfig == null) {
                    LOGGER.error("finding deliveryservice in cr-config for " + challenge.getXmlId());
                    return;
                }

                String staticEntryString = challenge.getFqdn();
                final ArrayNode domains = (ArrayNode) deliveryServiceConfig.get("domains");
                if (domains == null || domains.size() == 0) {
                    LOGGER.error("no domains found in cr-config for deliveryservice " + challenge.getXmlId());
                    return;
                }

                final Iterator<JsonNode> domainIter = domains.iterator();
                while(domainIter.hasNext()) {
                    final JsonNode domainNode = domainIter.next();
                    staticEntryString = staticEntryString.replace(domainNode.asText() + ".", "");
                }

                if (staticEntryString.endsWith(".")) {
                    staticEntryString = staticEntryString.substring(0, staticEntryString.length() - 1);
                }

                final ArrayNode staticDnsEntriesNode = updateStaticEntries(challenge, staticEntryString, mapper, deliveryServiceConfig);

                deliveryServiceConfig.set("staticDnsEntries", staticDnsEntriesNode);
                deliveryServicesNode.set(challenge.getXmlId(), deliveryServiceConfig);

            });

            final ObjectNode statsNode = (ObjectNode) mostRecentConfig.get("stats");
            statsNode.put("date", Instant.now().toEpochMilli() / 1000L);

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
            LOGGER.warn("Failed updating dns challenge txt record with data from " + dataBaseURL + ":", e);
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
            LOGGER.warn("Failed to build dns challenge data while verifying:", e);
        }

        return false;
    }

    @Override
    public String getWatcherConfigPrefix() {
        return "dnschallengemapping";
    }

    private String readConfigFile() {
        try (InputStream is = new FileInputStream(databasesDirectory.resolve(configFile).toString());
             BufferedReader buf = new BufferedReader(new InputStreamReader(is))
        ) {
            String line = buf.readLine();
            final StringBuilder sb = new StringBuilder();
            while (line != null) {
                sb.append(line).append('\n');
                line = buf.readLine();
            }
            return sb.toString();
        } catch (Exception e) {
            LOGGER.error("Could not read cr-config file " + configFile + ":", e);
            return null;
        }
    }

    private ArrayNode updateStaticEntries(final LetsEncryptDnsChallenge challenge, final String name, final ObjectMapper mapper, final ObjectNode deliveryServiceConfig) {
        ArrayNode staticDnsEntriesNode = mapper.createArrayNode();
        ArrayNode newStaticDnsEntriesNode = mapper.createArrayNode();

        if (deliveryServiceConfig.findValue("staticDnsEntries") != null) {
            staticDnsEntriesNode = (ArrayNode) deliveryServiceConfig.findValue("staticDnsEntries");
        }

        if (challenge.getRecord().isEmpty()) {
            for (int i = 0; i < staticDnsEntriesNode.size(); i++) {
                if (!staticDnsEntriesNode.get(i).get("name").equals(name)) {
                    newStaticDnsEntriesNode.add(i);
                }
            }
        } else {
            newStaticDnsEntriesNode = staticDnsEntriesNode;

            final ObjectNode newChildNode = mapper.createObjectNode();
            newChildNode.put("type", "TXT");
            newChildNode.put("name", name);
            newChildNode.put("value", challenge.getRecord());
            newChildNode.put("ttl", 10);

            newStaticDnsEntriesNode.add(newChildNode);
        }

        return newStaticDnsEntriesNode;
    }

    public void setConfigHandler(final ConfigHandler configHandler) {
        this.configHandler = configHandler;
    }
    public ConfigHandler getConfigHandler() {
        return this.configHandler;
    }

    public void setConfigFile(final String configFile) {
        this.configFile = configFile;
    }
}
