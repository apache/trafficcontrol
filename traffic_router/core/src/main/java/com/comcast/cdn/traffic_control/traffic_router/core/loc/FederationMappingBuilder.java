package com.comcast.cdn.traffic_control.traffic_router.core.loc;

import com.comcast.cdn.traffic_control.traffic_router.core.util.CidrAddress;
import org.apache.wicket.ajax.json.JSONArray;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;
import org.apache.wicket.ajax.json.JSONTokener;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.List;

public class FederationMappingBuilder {
    private static final Logger LOGGER = LoggerFactory.getLogger(FederationMappingBuilder.class);

    public FederationMapping fromJSON(final String json) throws JSONException, UnknownHostException {
        final JSONObject jsonObject = new JSONObject(new JSONTokener(json));

        final String cname = jsonObject.getString("cname");
        final int ttl = jsonObject.getInt("ttl");

        List<CidrAddress> network = null;
        if (jsonObject.has("resolve4")) {
            final JSONArray networkArray = jsonObject.getJSONArray("resolve4");

            try {
                network = buildAddresses(networkArray);
            }
            catch (JSONException e) {
                LOGGER.warn("Failed getting ipv4 address array");
            }
        }


        List<CidrAddress> network6 = null;
        if (jsonObject.has("resolve6")) {
            final JSONArray network6Array = jsonObject.getJSONArray("resolve6");
            try {
                network6 = buildAddresses(network6Array);
            }
            catch (JSONException e) {
                LOGGER.warn("Failed getting ipv6 address array");
            }
        }

        return new FederationMapping(cname, ttl, network, network6);
    }

    private List<CidrAddress> buildAddresses(final JSONArray networkArray) throws JSONException {
        final List<CidrAddress> network = new ArrayList<CidrAddress>();

        for (int i = 0; i < networkArray.length(); i++) {
            final String addressString = networkArray.getString(i);
            try {
                final CidrAddress cidrAddress = new CidrAddress(addressString);
                network.add(cidrAddress);
            } catch (NetworkNodeException e) {
                LOGGER.warn(e.getMessage());
            }
        }

        return network;
    }
}
