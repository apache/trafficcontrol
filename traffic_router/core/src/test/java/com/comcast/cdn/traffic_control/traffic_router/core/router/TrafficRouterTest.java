package com.comcast.cdn.traffic_control.traffic_router.core.router;

import com.comcast.cdn.traffic_control.traffic_router.core.cache.Cache;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheLocation;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheRegister;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.InetRecord;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryService;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.Dispersion;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.FederationRegistry;
import com.comcast.cdn.traffic_control.traffic_router.geolocation.Geolocation;
import com.comcast.cdn.traffic_control.traffic_router.core.request.DNSRequest;
import com.comcast.cdn.traffic_control.traffic_router.core.request.HTTPRequest;
import com.comcast.cdn.traffic_control.traffic_router.core.request.Request;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker.Track;
import com.comcast.cdn.traffic_control.traffic_router.core.util.CidrAddress;
import org.junit.Before;
import org.junit.Test;
import org.powermock.reflect.Whitebox;
import org.xbill.DNS.Type;

import java.util.ArrayList;
import java.util.Collection;
import java.util.List;

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
    private DeliveryService deliveryService;
    private FederationRegistry federationRegistry;

    @Before
    public void before() throws Exception {
        deliveryService = mock(DeliveryService.class);
        when(deliveryService.isAvailable()).thenReturn(true);
        when(deliveryService.isCoverageZoneOnly()).thenReturn(false);
        when(deliveryService.getDispersion()).thenReturn(mock(Dispersion.class));

        when(deliveryService.createURIString(any(HTTPRequest.class), any(Cache.class))).thenReturn("http://atscache.kabletown.net/index.html");

        List<InetRecord> inetRecords = new ArrayList<InetRecord>();
        InetRecord inetRecord = new InetRecord("cname1", 12345);
        inetRecords.add(inetRecord);

        federationRegistry = mock(FederationRegistry.class);
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
        Cache cache = mock(Cache.class);
        when(cache.hasDeliveryService(anyString())).thenReturn(true);
        CacheLocation cacheLocation = new CacheLocation("", "some-zone-id", new Geolocation(50,50));

        cacheLocation.addCache(cache);

        Collection<CacheLocation> cacheLocationCollection = new ArrayList<CacheLocation>();
        cacheLocationCollection.add(cacheLocation);

        CacheRegister cacheRegister = mock(CacheRegister.class);
        when(cacheRegister.getCacheLocations(null)).thenReturn(cacheLocationCollection);

        when(trafficRouter.getCacheRegister()).thenReturn(cacheRegister);
        when(deliveryService.isLocationAvailable(cacheLocation)).thenReturn(true);

        when(trafficRouter.selectCache(any(Request.class), any(DeliveryService.class), any(Track.class))).thenCallRealMethod();
        when(trafficRouter.selectCachesByGeo(any(Request.class), any(DeliveryService.class), any(CacheLocation.class), any(Track.class))).thenCallRealMethod();

        Geolocation clientLocation = new Geolocation(40, -100);
        when(trafficRouter.getClientLocation(any(Request.class), any(DeliveryService.class), any(CacheLocation.class), any(Track.class))).thenReturn(clientLocation);

        when(trafficRouter.getCachesByGeo(any(Request.class), any(DeliveryService.class), any(Geolocation.class), any(Track.class))).thenCallRealMethod();
        when(trafficRouter.orderCacheLocations(any(Collection.class), any(DeliveryService.class), any(Geolocation.class))).thenCallRealMethod();
        when(trafficRouter.getSupportingCaches(any(List.class), any(DeliveryService.class))).thenCallRealMethod();

        HTTPRequest httpRequest = new HTTPRequest();
        httpRequest.setClientIP("192.168.10.11");
        httpRequest.setHostname("ccr.example.com");

        Track track = spy(StatTracker.getTrack());

        trafficRouter.route(httpRequest, track);

        assertThat(track.getResult(), equalTo(Track.ResultType.GEO));
        assertThat(track.getResultLocation(), equalTo(new Geolocation(50, 50)));

        when(federationRegistry.findInetRecords(anyString(), any(CidrAddress.class))).thenReturn(null);

        DNSRequest dnsRequest = new DNSRequest();
        dnsRequest.setClientIP("192.168.1.2");
        dnsRequest.setClientIP("10.10.10.10");
        dnsRequest.setQtype(Type.A);

        track = StatTracker.getTrack();
        trafficRouter.route(dnsRequest, track);

        assertThat(track.getResult(), equalTo(Track.ResultType.GEO));
        assertThat(track.getResultLocation(), equalTo(new Geolocation(50, 50)));
    }
}
