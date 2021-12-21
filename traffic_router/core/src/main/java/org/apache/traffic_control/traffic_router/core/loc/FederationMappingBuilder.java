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

package org.apache.traffic_control.traffic_router.core.loc;

import org.apache.traffic_control.traffic_router.core.util.CidrAddress;
import org.apache.traffic_control.traffic_router.core.util.ComparableTreeSet;
import org.apache.traffic_control.traffic_router.core.util.JsonUtils;
import org.apache.traffic_control.traffic_router.core.util.JsonUtilsException;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import java.io.IOException;

public class FederationMappingBuilder {
    private final static Logger LOGGER = LogManager.getLogger(FederationMapping.class);


    public FederationMapping fromJSON(final String json) throws JsonUtilsException, IOException {
        final ObjectMapper mapper = new ObjectMapper();
        final JsonNode jsonNode = mapper.readTree(json);

        final String cname = JsonUtils.getString(jsonNode, "cname");
        final int ttl = JsonUtils.getInt(jsonNode, "ttl");

        final ComparableTreeSet<CidrAddress> network = new ComparableTreeSet<CidrAddress>();
        if (jsonNode.has("resolve4")) {
            final JsonNode networkList = JsonUtils.getJsonNode(jsonNode, "resolve4");
            network.addAll(buildAddresses(networkList));

        }

        final ComparableTreeSet<CidrAddress> network6 = new ComparableTreeSet<CidrAddress>();
        if (jsonNode.has("resolve6")) {
            final JsonNode network6List = JsonUtils.getJsonNode(jsonNode, "resolve6");
            network6.addAll(buildAddresses(network6List));
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
