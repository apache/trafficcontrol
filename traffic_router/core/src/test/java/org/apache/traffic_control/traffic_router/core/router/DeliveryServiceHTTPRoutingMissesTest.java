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

package org.apache.traffic_control.traffic_router.core.router;

import org.apache.traffic_control.traffic_router.core.ds.DeliveryService;
import org.apache.traffic_control.traffic_router.core.request.HTTPRequest;
import org.apache.traffic_control.traffic_router.core.router.StatTracker.Track;
import org.apache.traffic_control.traffic_router.core.router.StatTracker.Track.ResultDetails;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.Before;
import org.junit.Test;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.mockito.Mockito.*;
import static org.powermock.reflect.Whitebox.setInternalState;

public class DeliveryServiceHTTPRoutingMissesTest {

    private DeliveryService deliveryService;
    private HTTPRequest httpRequest;
    private Track track;
    private JsonNode bypassDestination;

    @Before
    public void before() throws Exception {
        ObjectMapper mapper = new ObjectMapper();
        JsonNode unusedByTest = mock(JsonNode.class);
        JsonNode ttls = mock(JsonNode.class);
        when(unusedByTest.get("ttls")).thenReturn(ttls);
        when(unusedByTest.has("routingName")).thenReturn(true);
        when(unusedByTest.get("routingName")).thenReturn(mapper.readTree("\"edge\""));
        when(unusedByTest.has("coverageZoneOnly")).thenReturn(true);
        when(unusedByTest.get("coverageZoneOnly")).thenReturn(mapper.readTree("true"));
        when(unusedByTest.has("deepCachingType")).thenReturn(true);
        when(unusedByTest.get("deepCachingType")).thenReturn(mapper.readTree("\"NEVER\""));
        deliveryService = new DeliveryService("ignoredbytest", unusedByTest);
        httpRequest = mock(HTTPRequest.class);
        track = StatTracker.getTrack();
        bypassDestination = mock(JsonNode.class);
        setInternalState(deliveryService, "bypassDestination", bypassDestination);
    }

    @Test
    public void itSetsDetailsWhenNoBypass() throws Exception {
        JsonNode nullBypassDestination = null;
        setInternalState(deliveryService, "bypassDestination", nullBypassDestination);
        deliveryService.getFailureHttpResponse(httpRequest, track);
        assertThat(track.getResultDetails(), equalTo(ResultDetails.DS_NO_BYPASS));
        assertThat(track.getResult(), equalTo(Track.ResultType.MISS));
    }

    @Test
    public void itSetsDetailsWhenNoHTTPBypass() throws Exception {
        when(bypassDestination.get("HTTP")).thenReturn(null);

        deliveryService.getFailureHttpResponse(httpRequest, track);
        assertThat(track.getResultDetails(), equalTo(ResultDetails.DS_NO_BYPASS));
        assertThat(track.getResult(), equalTo(Track.ResultType.MISS));
    }

    @Test
    public void itSetsDetailsWhenNoFQDNBypass() throws Exception {
        ObjectMapper mapper = new ObjectMapper();
        JsonNode httpJsonObject = mapper.createObjectNode();
        httpJsonObject = spy(httpJsonObject);
        doReturn(null).when(httpJsonObject).get("fqdn");

        when(bypassDestination.get("HTTP")).thenReturn(httpJsonObject);

        deliveryService.getFailureHttpResponse(httpRequest, track);

        verify(httpJsonObject).get("fqdn");

        assertThat(track.getResultDetails(), equalTo(ResultDetails.DS_NO_BYPASS));
        assertThat(track.getResult(), equalTo(Track.ResultType.MISS));
    }
}
