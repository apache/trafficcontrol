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
