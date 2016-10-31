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

import org.apache.wicket.ajax.json.JSONArray;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;
import org.apache.wicket.ajax.json.JSONTokener;

import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.List;

public class FederationsBuilder {

    public List<Federation> fromJSON(final String jsonString) throws JSONException, UnknownHostException {
        final List<Federation> federations = new ArrayList<Federation>();

        final JSONObject jsonObject = new JSONObject(new JSONTokener(jsonString));
        final JSONArray federationsArray = jsonObject.getJSONArray("response");

        for (int i = 0; i < federationsArray.length(); i++) {
            final JSONObject jsonObject1 = federationsArray.getJSONObject(i);
            final String deliveryService = jsonObject1.getString("deliveryService");

            final List<FederationMapping> mappings = new ArrayList<FederationMapping>();

            final JSONArray mappingsArray = jsonObject1.getJSONArray("mappings");
            final FederationMappingBuilder federationMappingBuilder = new FederationMappingBuilder();

            for (int j = 0; j < mappingsArray.length(); j++) {
                mappings.add(federationMappingBuilder.fromJSON(mappingsArray.getJSONObject(j).toString()));
            }

            final Federation federation = new Federation(deliveryService, mappings);
            federations.add(federation);
        }

        return federations;
    }

}
