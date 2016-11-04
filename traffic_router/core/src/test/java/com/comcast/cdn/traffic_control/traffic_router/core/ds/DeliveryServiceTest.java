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
