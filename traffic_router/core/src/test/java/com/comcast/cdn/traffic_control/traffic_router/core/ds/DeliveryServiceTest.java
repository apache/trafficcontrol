package com.comcast.cdn.traffic_control.traffic_router.core.ds;

import org.json.JSONObject;
import org.junit.Test;

import java.util.HashSet;
import java.util.Set;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.containsInAnyOrder;
import static org.hamcrest.Matchers.equalTo;

public class DeliveryServiceTest {

    @Test
    public void itHandlesLackOfRequestHeaderNamesInJSON() throws Exception {
        JSONObject jsonConfiguration = new JSONObject();
        jsonConfiguration.put("coverageZoneOnly", false);
        DeliveryService deliveryService = new DeliveryService("a-delivery-service", jsonConfiguration);
        assertThat(deliveryService.getRequestHeaders().size(), equalTo(0));
    }

    @Test
    public void itConfiguresRequestHeadersFromJSON() throws Exception {
        JSONObject jsonConfiguration = new JSONObject();
        jsonConfiguration.put("coverageZoneOnly", false);

        Set<String> requestHeaderNames = new HashSet<String>();
        requestHeaderNames.add("Cache-Control");
        requestHeaderNames.add("Cookie");
        requestHeaderNames.add("Content-Type");
        requestHeaderNames.add("If-Modified-Since");

        jsonConfiguration.put("requestHeaders", requestHeaderNames);

        DeliveryService deliveryService = new DeliveryService("a-delivery-service", jsonConfiguration);

        assertThat(deliveryService.getRequestHeaders(), containsInAnyOrder("Cache-Control", "Cookie", "Content-Type", "If-Modified-Since"));
    }

}
