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

import org.apache.traffic_control.traffic_router.core.edge.CacheRegister;
import org.apache.traffic_control.traffic_router.core.request.HTTPRequest;
import org.junit.Before;
import org.junit.Test;

import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.verify;
import static org.powermock.api.mockito.PowerMockito.doCallRealMethod;
import static org.powermock.api.mockito.PowerMockito.spy;
import static org.powermock.reflect.Whitebox.setInternalState;

public class TrafficRouterHTTPRoutingMissesTest {
    private HTTPRequest request;
    private TrafficRouter trafficRouter;
    private StatTracker.Track track;
    private CacheRegister cacheRegister;

    @Before
    public void before() throws Exception {
        request = new HTTPRequest();
        request.setClientIP("192.168.34.56");

        cacheRegister = mock(CacheRegister.class);
        trafficRouter = mock(TrafficRouter.class);

        track = spy(StatTracker.getTrack());
        setInternalState(trafficRouter, "cacheRegister", cacheRegister);
        doCallRealMethod().when(trafficRouter).route(request, track);
        doCallRealMethod().when(trafficRouter).singleRoute(request, track);
    }

    @Test
    public void itSetsDetailsWhenNoDeliveryService() throws Exception {
        trafficRouter.route(request, track);

        verify(track).setResult(StatTracker.Track.ResultType.DS_MISS);
        verify(track).setResultDetails(StatTracker.Track.ResultDetails.DS_NOT_FOUND);
    }
}
