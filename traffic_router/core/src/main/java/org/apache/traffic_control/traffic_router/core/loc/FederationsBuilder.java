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

import org.apache.traffic_control.traffic_router.core.util.JsonUtils;
import org.apache.traffic_control.traffic_router.core.util.JsonUtilsException;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

public class FederationsBuilder {

    public List<Federation> fromJSON(final String jsonString) throws JsonUtilsException, IOException {
        final List<Federation> federations = new ArrayList<Federation>();

        final ObjectMapper mapper = new ObjectMapper();
        final JsonNode jsonObject = mapper.readTree(jsonString);
        final JsonNode federationList = JsonUtils.getJsonNode(jsonObject, "response");

        for (final JsonNode currFederation : federationList) {
            final String deliveryService = JsonUtils.getString(currFederation, "deliveryService");

            final List<FederationMapping> mappings = new ArrayList<FederationMapping>();

            final JsonNode mappingsList = JsonUtils.getJsonNode(currFederation, "mappings");
            final FederationMappingBuilder federationMappingBuilder = new FederationMappingBuilder();

            for (final JsonNode mapping : mappingsList) {
                mappings.add(federationMappingBuilder.fromJSON(mapping.toString()));
            }

            final Federation federation = new Federation(deliveryService, mappings);
            federations.add(federation);
        }

        return federations;
    }

}
