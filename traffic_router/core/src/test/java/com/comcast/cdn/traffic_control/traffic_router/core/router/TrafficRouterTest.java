package com.comcast.cdn.traffic_control.traffic_router.core.router;

import com.comcast.cdn.traffic_control.traffic_router.core.cache.Cache;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheLocation;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.InetRecord;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryService;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.Dispersion;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.FederationRegistry;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.Geolocation;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.RegionalGeoResult;
import com.comcast.cdn.traffic_control.traffic_router.core.request.DNSRequest;
import com.comcast.cdn.traffic_control.traffic_router.core.request.HTTPRequest;
import com.comcast.cdn.traffic_control.traffic_router.core.request.Request;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker.Track;
import com.comcast.cdn.traffic_control.traffic_router.core.util.CidrAddress;
import org.junit.Before;
import org.junit.Test;
import org.powermock.reflect.Whitebox;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.containsInAnyOrder;
import static org.hamcrest.Matchers.equalTo;
import static org.mockito.Matchers.any;
import static org.mockito.Matchers.anyBoolean;
import static org.mockito.Matchers.anyString;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.spy;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

public class TrafficRouterTest {

    private TrafficRouter trafficRouter;

    @Before
    public void before() throws Exception {
        DeliveryService deliveryService = mock(DeliveryService.class);
        when(deliveryService.isAvailable()).thenReturn(true);
        when(deliveryService.isCoverageZoneOnly()).thenReturn(false);
        when(deliveryService.getDispersion()).thenReturn(mock(Dispersion.class));

        when(deliveryService.createURIString(any(HTTPRequest.class), any(Cache.class))).thenReturn("http://atscache.kabletown.net/index.html");

        List<InetRecord> inetRecords = new ArrayList<InetRecord>();
        InetRecord inetRecord = new InetRecord("cname1", 12345);
        inetRecords.add(inetRecord);

        FederationRegistry federationRegistry = mock(FederationRegistry.class);
        when(federationRegistry.findInetRecords(anyString(), any(CidrAddress.class))).thenReturn(inetRecords);

        trafficRouter = mock(TrafficRouter.class);
        Whitebox.setInternalState(trafficRouter, "federationRegistry", federationRegistry);

        when(trafficRouter.route(any(DNSRequest.class), any(Track.class))).thenCallRealMethod();
        when(trafficRouter.route(any(HTTPRequest.class), any(Track.class))).thenCallRealMethod();
        when(trafficRouter.selectDeliveryService(any(Request.class), anyBoolean())).thenReturn(deliveryService);
    }

    @Test
    public void itCreatesDnsResultsFromFederationMappingHit() throws Exception {
        DNSRequest request = new DNSRequest();
        request.setClientIP("192.168.10.11");
        request.setHostname("edge.example.com");

        Track track = spy(StatTracker.getTrack());

        DNSRouteResult result = trafficRouter.route(request, track);

        assertThat(result.getAddresses(), containsInAnyOrder(new InetRecord("cname1", 12345)));
        verify(track).setRouteType(Track.RouteType.DNS, "edge.example.com");
    }

    @Test
    public void itCreatesHttpResults() throws Exception {
        HTTPRequest httpRequest = new HTTPRequest();
        httpRequest.setClientIP("192.168.10.11");
        httpRequest.setHostname("ccr.example.com");

        Track track = spy(StatTracker.getTrack());

        HTTPRouteResult httpRouteResult = trafficRouter.route(httpRequest, track);

        assertThat(httpRouteResult.getUrl().toString(), equalTo("http://atscache.kabletown.net/index.html"));
    }

    @Test
    public void itSetsResultToGeo() throws Exception {

        when(trafficRouter.selectCache(any(Request.class), any(DeliveryService.class), any(Track.class), any(RegionalGeoResult.class))).thenCallRealMethod();
        when(trafficRouter.selectCachesByGeo(any(Request.class), any(DeliveryService.class), any(CacheLocation.class), any(Track.class), any(RegionalGeoResult.class))).thenCallRealMethod();

        Geolocation clientLocation = new Geolocation(40, -100);
        when(trafficRouter.getClientLocation(any(Request.class), any(DeliveryService.class), any(CacheLocation.class))).thenReturn(clientLocation);

        List<Cache> caches = new ArrayList<Cache>();
        Cache cache = mock(Cache.class);
        caches.add(cache);

        when(trafficRouter.getCachesByGeo(any(Request.class), any(DeliveryService.class), any(Geolocation.class), any(Map.class))).thenReturn(caches);

        HTTPRequest httpRequest = new HTTPRequest();
        httpRequest.setClientIP("192.168.10.11");
        httpRequest.setHostname("ccr.example.com");

        Track track = spy(StatTracker.getTrack());

        trafficRouter.route(httpRequest, track);

        assertThat(track.getResult(), equalTo(Track.ResultType.GEO));
    }
}
