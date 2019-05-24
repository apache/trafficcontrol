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

import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouter;
import com.comcast.cdn.traffic_control.traffic_router.core.util.AbstractResourceWatcher;
import com.fasterxml.jackson.core.JsonFactory;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.log4j.Logger;
import org.xbill.DNS.*;

import java.time.Duration;
import java.util.HashMap;
import java.util.List;

public class LetsEncryptDnsChallengeWatcher extends AbstractResourceWatcher {
    private static final Logger LOGGER = Logger.getLogger(LetsEncryptDnsChallengeWatcher.class);
    public static final String DEFAULT_LE_DNS_CHALLENGE_URL = "http://${toHostname}/api/1.4/letsencrypt/dnsrecords/";
    private static final long TTL = Duration.ofMinutes(15).getSeconds();
    private TrafficRouter trafficRouter;

    public void setTrafficRouter(final TrafficRouter trafficRouter) {
        this.trafficRouter = trafficRouter;
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

            challengeList.forEach(challenge -> {
                final StringBuilder sb = new StringBuilder();
                sb.append(challenge.getFqdn());
                if (!challenge.getFqdn().endsWith(".")) {
                     sb.append('.');
                }
                final String fqdn = sb.toString();
                try {
                    final Name zoneName = new Name(fqdn);
                    final Record r = new TXTRecord(zoneName, DClass.IN, TTL, challenge.getRecord());
                    trafficRouter.getZoneManager().getZone(zoneName).addRecord(r);
                } catch (TextParseException e) {
                    LOGGER.warn("Failed to parse zone from fqdn: " + fqdn);
                }
            });

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
}
