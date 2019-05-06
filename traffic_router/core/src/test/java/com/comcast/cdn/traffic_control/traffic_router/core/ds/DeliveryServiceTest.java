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

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.Test;
import org.junit.Assert;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.containsInAnyOrder;
import static org.hamcrest.Matchers.equalTo;

public class DeliveryServiceTest {

    @Test
    public void itHandlesLackOfRequestHeaderNamesInJSON() throws Exception {
        final ObjectMapper mapper = new ObjectMapper();
        final String jsonStr = "{\"routingName\":\"edge\",\"coverageZoneOnly\":false}";
        final JsonNode jsonConfiguration = mapper.readTree(jsonStr);
        DeliveryService deliveryService = new DeliveryService("a-delivery-service", jsonConfiguration);
        assertThat(deliveryService.getRequestHeaders().size(), equalTo(0));
    }

    @Test
    public void itHandlesLackOfConsistentHashQueryParamsInJSON() throws Exception {
        final ObjectMapper mapper = new ObjectMapper();
        final JsonNode json = mapper.readTree("{\"routingName\":\"edge\",\"coverageZoneOnly\":false}");
        DeliveryService d = new DeliveryService("test", json);
        assert d.consistentHashQueryParams != null;
        assert d.consistentHashQueryParams.size() == 0;
    }

    @Test
    public void itHandlesDuplicatesInConsistentHashQueryParams() throws Exception {
        final ObjectMapper mapper = new ObjectMapper();
        final JsonNode json = mapper.readTree("{\"routingName\":\"edge\",\"coverageZoneOnly\":false,\"consistentHashQueryParams\":[\"test\", \"quest\", \"test\"]}");
        DeliveryService d = new DeliveryService("test", json);
        assert d.consistentHashQueryParams != null;
        assert d.consistentHashQueryParams.size() == 2;
        assert d.consistentHashQueryParams.contains("test");
        assert d.consistentHashQueryParams.contains("quest");
    }

    @Test
    public void itConfiguresRequestHeadersFromJSON() throws Exception {
        final ObjectMapper mapper = new ObjectMapper();
        final String jsonStr = "{\"routingName\":\"edge\",\"requestHeaders\":[\"Cookie\",\"Cache-Control\",\"If-Modified-Since\",\"Content-Type\"],\"coverageZoneOnly\":false}";
        final JsonNode jsonConfiguration = mapper.readTree(jsonStr);

        DeliveryService deliveryService = new DeliveryService("a-delivery-service", jsonConfiguration);

        assertThat(deliveryService.getRequestHeaders(), containsInAnyOrder("Cache-Control", "Cookie", "Content-Type", "If-Modified-Since"));
    }

}
