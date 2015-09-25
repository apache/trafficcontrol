package com.comcast.cdn.traffic_control.traffic_router.core.router;

import com.comcast.cdn.traffic_control.traffic_router.core.request.HTTPRequest;
import org.junit.Before;
import org.junit.Test;

import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.verify;
import static org.powermock.api.mockito.PowerMockito.doCallRealMethod;
import static org.powermock.api.mockito.PowerMockito.spy;

public class TrafficRouterHTTPRoutingMissesTest {
    private HTTPRequest request;
    private TrafficRouter trafficRouter;
    private StatTracker.Track track;

    @Before
    public void before() throws Exception {
        request = new HTTPRequest();
        request.setClientIP("192.168.34.56");

        trafficRouter = mock(TrafficRouter.class);

        track = spy(StatTracker.getTrack());
        doCallRealMethod().when(trafficRouter).route(request, track);
    }

    @Test
    public void itSetsDetailsWhenNoDeliveryService() throws Exception {
        trafficRouter.route(request, track);

        verify(track).setResult(StatTracker.Track.ResultType.DS_MISS);
        verify(track).setResultDetails(StatTracker.Track.ResultDetails.DS_NOT_FOUND);
    }
}
