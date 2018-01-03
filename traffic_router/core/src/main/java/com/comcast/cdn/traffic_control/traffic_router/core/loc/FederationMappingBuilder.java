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

package com.comcast.cdn.traffic_control.traffic_router.core.loc;

import com.comcast.cdn.traffic_control.traffic_router.core.util.CidrAddress;
import com.comcast.cdn.traffic_control.traffic_router.core.util.ComparableTreeSet;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;

public class FederationMappingBuilder {
    private static final Logger LOGGER = LoggerFactory.getLogger(FederationMappingBuilder.class);

    @SuppressWarnings("PMD.NPathComplexity")
    public FederationMapping fromJSON(final String json) throws IOException {
        final ObjectMapper mapper = new ObjectMapper();
        final JsonNode jsonNode = mapper.readTree(json);

        final String cname = jsonNode.has("cname") ? jsonNode.get("cname").asText() : "";
        final int ttl = jsonNode.has("ttl") ? jsonNode.get("ttl").asInt() : 0;

        final ComparableTreeSet<CidrAddress> network = new ComparableTreeSet<CidrAddress>();
        if (jsonNode.has("resolve4")) {
            final JsonNode networkList = jsonNode.get("resolve4");

            try {
                network.addAll(buildAddresses(networkList));
            }
            catch (Exception e) {
                LOGGER.warn("Failed getting ipv4 address array likely due to bad json data: " + e.getMessage());
            }
        }

        final ComparableTreeSet<CidrAddress> network6 = new ComparableTreeSet<CidrAddress>();
        if (jsonNode.has("resolve6")) {
            final JsonNode network6List = jsonNode.get("resolve6");
            try {
                network6.addAll(buildAddresses(network6List));
            }
            catch (Exception e) {
                LOGGER.warn("Failed getting ipv6 address array likely due to bad json data: " + e.getMessage());
            }
        }

        return new FederationMapping(cname, ttl, network, network6);
    }

    private ComparableTreeSet<CidrAddress> buildAddresses(final JsonNode networkArray) {
        final ComparableTreeSet<CidrAddress> network = new ComparableTreeSet<CidrAddress>();

        for (final JsonNode currNetwork : networkArray) {
            final String addressString = currNetwork.asText();
            try {
                final CidrAddress cidrAddress = CidrAddress.fromString(addressString);
                network.add(cidrAddress);
            } catch (NetworkNodeException e) {
                LOGGER.warn(e.getMessage());
            }
        }

        return network;
    }
}
