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

package com.comcast.cdn.traffic_control.traffic_router.core.router;

import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryService;
import com.comcast.cdn.traffic_control.traffic_router.core.request.HTTPRequest;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker.Track;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker.Track.ResultDetails;

import org.json.JSONObject;
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
    private JSONObject bypassDestination;

    @Before
    public void before() throws Exception {
        JSONObject unusedByTest = mock(JSONObject.class);
        JSONObject ttls = mock(JSONObject.class);
        when(unusedByTest.optJSONObject("ttls")).thenReturn(ttls);
        deliveryService = new DeliveryService("ignoredbytest", unusedByTest);
        httpRequest = mock(HTTPRequest.class);
        track = StatTracker.getTrack();
        bypassDestination = mock(JSONObject.class);
        setInternalState(deliveryService, "bypassDestination", bypassDestination);
    }

    @Test
    public void itSetsDetailsWhenNoBypass() throws Exception {
        JSONObject nullBypassDestination = null;
        setInternalState(deliveryService, "bypassDestination", nullBypassDestination);
        deliveryService.getFailureHttpResponse(httpRequest, track);
        assertThat(track.getResultDetails(), equalTo(ResultDetails.DS_NO_BYPASS));
        assertThat(track.getResult(), equalTo(Track.ResultType.MISS));
    }

    @Test
    public void itSetsDetailsWhenNoHTTPBypass() throws Exception {
        when(bypassDestination.optJSONObject("HTTP")).thenReturn(null);

        deliveryService.getFailureHttpResponse(httpRequest, track);
        assertThat(track.getResultDetails(), equalTo(ResultDetails.DS_NO_BYPASS));
        assertThat(track.getResult(), equalTo(Track.ResultType.MISS));
    }

    @Test
    public void itSetsDetailsWhenNoFQDNBypass() throws Exception {
        JSONObject httpJsonObject = new JSONObject();
        httpJsonObject = spy(httpJsonObject);
        doReturn(null).when(httpJsonObject).optString("fqdn");

        when(bypassDestination.optJSONObject("HTTP")).thenReturn(httpJsonObject);

        deliveryService.getFailureHttpResponse(httpRequest, track);

        verify(httpJsonObject).optString("fqdn");

        assertThat(track.getResultDetails(), equalTo(ResultDetails.DS_NO_BYPASS));
        assertThat(track.getResult(), equalTo(Track.ResultType.MISS));
    }
}
