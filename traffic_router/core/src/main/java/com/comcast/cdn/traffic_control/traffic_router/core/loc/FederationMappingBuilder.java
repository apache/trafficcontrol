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
import org.apache.wicket.ajax.json.JSONArray;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;
import org.apache.wicket.ajax.json.JSONTokener;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.UnknownHostException;

public class FederationMappingBuilder {
    private static final Logger LOGGER = LoggerFactory.getLogger(FederationMappingBuilder.class);

    public FederationMapping fromJSON(final String json) throws JSONException, UnknownHostException {
        final JSONObject jsonObject = new JSONObject(new JSONTokener(json));

        final String cname = jsonObject.getString("cname");
        final int ttl = jsonObject.getInt("ttl");

        final ComparableTreeSet<CidrAddress> network = new ComparableTreeSet<CidrAddress>();
        if (jsonObject.has("resolve4")) {
            final JSONArray networkArray = jsonObject.getJSONArray("resolve4");

            try {
                network.addAll(buildAddresses(networkArray));
            }
            catch (JSONException e) {
                LOGGER.warn("Failed getting ipv4 address array likely due to bad json data: " + e.getMessage());
            }
        }


        final ComparableTreeSet<CidrAddress> network6 = new ComparableTreeSet<CidrAddress>();
        if (jsonObject.has("resolve6")) {
            final JSONArray network6Array = jsonObject.getJSONArray("resolve6");
            try {
                network6.addAll(buildAddresses(network6Array));
            }
            catch (JSONException e) {
                LOGGER.warn("Failed getting ipv6 address array likely due to bad json data: " + e.getMessage());
            }
        }

        return new FederationMapping(cname, ttl, network, network6);
    }

    private ComparableTreeSet<CidrAddress> buildAddresses(final JSONArray networkArray) throws JSONException {
        final ComparableTreeSet<CidrAddress> network = new ComparableTreeSet<CidrAddress>();

        for (int i = 0; i < networkArray.length(); i++) {
            final String addressString = networkArray.getString(i);
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
