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

package org.apache.traffic_control.traffic_router.core.ds;

import org.apache.traffic_control.traffic_router.core.request.HTTPRequest;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.Test;
import org.powermock.reflect.Whitebox;

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
        assert d.getConsistentHashQueryParams() != null;
        assert d.getConsistentHashQueryParams().size() == 0;
    }

    @Test
    public void itHandlesDuplicatesInConsistentHashQueryParams() throws Exception {
        final ObjectMapper mapper = new ObjectMapper();
        final JsonNode json = mapper.readTree("{\"routingName\":\"edge\",\"coverageZoneOnly\":false,\"consistentHashQueryParams\":[\"test\", \"quest\", \"test\"]}");
        DeliveryService d = new DeliveryService("test", json);
        assert d.getConsistentHashQueryParams() != null;
        assert d.getConsistentHashQueryParams().size() == 2;
        assert d.getConsistentHashQueryParams().contains("test");
        assert d.getConsistentHashQueryParams().contains("quest");
    }

    @Test
    public void itExtractsQueryParams() throws Exception {
        final JsonNode json = (new ObjectMapper()).readTree("{\"routingName\":\"edge\",\"coverageZoneOnly\":false,\"consistentHashQueryParams\":[\"test\", \"quest\"]}");
        final HTTPRequest r = new HTTPRequest();
        r.setPath("/path1234/some_stream_name1234/some_other_info.m3u8");
        r.setQueryString("test=value&foo=fizz&quest=oth%20ervalue&bar=buzz");
        assert (new DeliveryService("test", json)).extractSignificantQueryParams(r).equals("quest=oth ervaluetest=value");
    }

    @Test
    public void itConfiguresRequestHeadersFromJSON() throws Exception {
        final ObjectMapper mapper = new ObjectMapper();
        final String jsonStr = "{\"routingName\":\"edge\",\"requestHeaders\":[\"Cookie\",\"Cache-Control\",\"If-Modified-Since\",\"Content-Type\"],\"coverageZoneOnly\":false}";
        final JsonNode jsonConfiguration = mapper.readTree(jsonStr);

        DeliveryService deliveryService = new DeliveryService("a-delivery-service", jsonConfiguration);

        assertThat(deliveryService.getRequestHeaders(), containsInAnyOrder("Cache-Control", "Cookie", "Content-Type", "If-Modified-Since"));
    }

    @Test
    public void itAddsRequiredCapabilities() throws Exception {
        final ObjectMapper mapper = new ObjectMapper();
        final JsonNode jsonConfiguration = mapper.readTree("{\"requiredCapabilities\":[\"all-read\",\"all-write\",\"cdn-read\"],\"routingName\":\"edge\",\"coverageZoneOnly\":false}");
        final DeliveryService deliveryService = new DeliveryService("has-required-capabilities", jsonConfiguration);

        assertThat(Whitebox.getInternalState(deliveryService, "requiredCapabilities"), containsInAnyOrder("all-read", "all-write", "cdn-read"));
    }
}
