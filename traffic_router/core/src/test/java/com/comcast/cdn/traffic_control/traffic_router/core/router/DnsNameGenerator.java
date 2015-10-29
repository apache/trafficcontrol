package com.comcast.cdn.traffic_control.traffic_router.core.router;

import org.json.JSONArray;
import org.json.JSONObject;

import java.util.ArrayList;
import java.util.List;

// Attempts to generate names like 'www.[foo].kabletown.com' to do dns queries against traffic router
// Tries to pull 'whole' words from the regex of cr-config
public class DnsNameGenerator {
    public List<String> getNames(JSONObject deliveryServicesConfig) throws Exception {
        List<String> names = new ArrayList<String>();

        for (String deliveryServiceId : JSONObject.getNames(deliveryServicesConfig)) {
            final JSONArray matchsets = deliveryServicesConfig
                .getJSONObject(deliveryServiceId)
                .getJSONArray("matchsets");

            for (int i = 0; i < matchsets.length(); i++) {
                final JSONObject matchset = matchsets.getJSONObject(i);

                if (!"DNS".equals(matchset.getString("protocol"))) {
                    continue;
                }

                final JSONArray list = matchset.getJSONArray("matchlist");
                for (int j = 0; j < list.length(); j++) {
                    // Not bulletproof
                    final String name = list.getJSONObject(j).getString("regex")
                        .replaceAll("\\.", "")
                        .replaceAll("\\*", "")
                        .replaceAll("\\\\", "");

                    names.add("www." + name + ".kabletown.com");
                }
            }
        }

        return names;
    }
}
