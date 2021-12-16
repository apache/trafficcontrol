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

import org.apache.traffic_control.traffic_router.core.edge.CacheLocation;
import org.apache.traffic_control.traffic_router.core.edge.CacheRegister;
import org.apache.traffic_control.traffic_router.core.edge.Node.IPVersions;
import org.apache.traffic_control.traffic_router.core.ds.DeliveryService;
import org.apache.traffic_control.traffic_router.core.loc.FederationRegistry;
import org.apache.traffic_control.traffic_router.core.request.DNSRequest;
import org.apache.traffic_control.traffic_router.core.request.Request;
import org.apache.traffic_control.traffic_router.core.router.StatTracker.Track;
import org.apache.traffic_control.traffic_router.core.router.StatTracker.Track.ResultType;
import org.apache.traffic_control.traffic_router.core.router.StatTracker.Track.ResultDetails;

import org.apache.traffic_control.traffic_router.core.util.CidrAddress;
import com.fasterxml.jackson.databind.JsonNode;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.core.classloader.annotations.PowerMockIgnore;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;
import org.powermock.reflect.Whitebox;
import org.xbill.DNS.Name;
import org.xbill.DNS.Type;

import static org.mockito.Mockito.*;
import static org.powermock.api.mockito.PowerMockito.doCallRealMethod;
import static org.powermock.api.mockito.PowerMockito.spy;
import static org.powermock.reflect.Whitebox.setInternalState;

@RunWith(PowerMockRunner.class)
@PrepareForTest({DeliveryService.class, TrafficRouter.class})
@PowerMockIgnore("javax.management.*")
public class DNSRoutingMissesTest {

    private DNSRequest request;
    private TrafficRouter trafficRouter;
    private Track track;

    @Before
    public void before() throws Exception {
        Name name = Name.fromString("edge.foo-img.kabletown.com");
        request = new DNSRequest("foo-img.kabletown.com", name, Type.A);

        request.setClientIP("192.168.34.56");
        request.setHostname(name.relativize(Name.root).toString());

        FederationRegistry federationRegistry = mock(FederationRegistry.class);
        when(federationRegistry.findInetRecords(anyString(), any(CidrAddress.class))).thenReturn(null);

        trafficRouter = mock(TrafficRouter.class);
        when(trafficRouter.getCacheRegister()).thenReturn(mock(CacheRegister.class));
        Whitebox.setInternalState(trafficRouter, "federationRegistry", federationRegistry);
        when(trafficRouter.selectCachesByGeo(any(), any(), any(), any(), any())).thenCallRealMethod();

        track = spy(StatTracker.getTrack());
        doCallRealMethod().when(trafficRouter).route(request, track);
    }

    @Test
    public void itSetsDetailsWhenNoDeliveryService() throws Exception {
        trafficRouter.route(request, track);

        verify(track).setResult(ResultType.MISS);
        verify(track).setResultDetails(ResultDetails.LOCALIZED_DNS);
    }

    // When the delivery service is unavailable ...
    @Test
    public void itSetsDetailsWhenNoBypass() throws Exception {
        DeliveryService deliveryService = mock(DeliveryService.class);
        when(deliveryService.isAvailable()).thenReturn(false);
        when(deliveryService.getFailureDnsResponse(request, track)).thenCallRealMethod();
        when(deliveryService.getRoutingName()).thenReturn("edge");
        when(deliveryService.isDns()).thenReturn(true);

        doReturn(deliveryService).when(trafficRouter).selectDeliveryService(request);

        trafficRouter.route(request, track);

        verify(track).setResult(ResultType.MISS);
        verify(track).setResultDetails(ResultDetails.DS_NO_BYPASS);
    }

    @Test
    public void itSetsDetailsWhenBypassDestination() throws Exception {
        DeliveryService deliveryService = mock(DeliveryService.class);
        when(deliveryService.isAvailable()).thenReturn(false);
        when(deliveryService.getFailureDnsResponse(request, track)).thenCallRealMethod();
        when(deliveryService.getRoutingName()).thenReturn("edge");
        when(deliveryService.isDns()).thenReturn(true);

        doReturn(deliveryService).when(trafficRouter).selectDeliveryService(request);

        JsonNode bypassDestination = mock(JsonNode.class);
        when(bypassDestination.get("DNS")).thenReturn(null);

        setInternalState(deliveryService, "bypassDestination", bypassDestination);

        trafficRouter.route(request, track);

        verify(track).setResult(ResultType.DS_REDIRECT);
        verify(track).setResultDetails(ResultDetails.DS_BYPASS);
    }

    // The Delivery Service is available but we don't find the cache in the coverage zone map

    // - and DS doesn't support other lookups
    @Test
    public void itSetsDetailsAboutMissesWhenOnlyCoverageZoneSupported() throws Exception {
        DeliveryService deliveryService = mock(DeliveryService.class);
        doReturn(true).when(deliveryService).isAvailable();
        when(deliveryService.getRoutingName()).thenReturn("edge");
        when(deliveryService.isDns()).thenReturn(true);
        when(deliveryService.isCoverageZoneOnly()).thenReturn(true);

        doReturn(deliveryService).when(trafficRouter).selectDeliveryService(any(Request.class));
        trafficRouter.route(request, track);

        verify(track).setResult(ResultType.MISS);
        verify(track).setResultDetails(ResultDetails.DS_CZ_ONLY);
    }

    // 1. We got an unsupported cache location from the coverage zone map
    // 2. we looked up the client location from maxmind
    // 3. delivery service says the client location is unsupported
    @Test
    public void itSetsDetailsWhenClientGeolocationNotSupported() throws Exception {
        DeliveryService deliveryService = mock(DeliveryService.class);
        doReturn(true).when(deliveryService).isAvailable();
        when(deliveryService.getRoutingName()).thenReturn("edge");
        when(deliveryService.isDns()).thenReturn(true);

        when(deliveryService.isCoverageZoneOnly()).thenReturn(false);

        doReturn(deliveryService).when(trafficRouter).selectDeliveryService(request);

        trafficRouter.route(request, track);

        verify(track).setResult(ResultType.MISS);
        verify(track).setResultDetails(ResultDetails.DS_CLIENT_GEO_UNSUPPORTED);

    }

    @Test
    public void itSetsDetailsWhenCacheNotFoundByGeolocation() throws Exception {
        doCallRealMethod().when(trafficRouter).selectCachesByGeo(anyString(), any(DeliveryService.class), any(CacheLocation.class), any(Track.class), any(IPVersions.class));
        CacheLocation cacheLocation = mock(CacheLocation.class);
        CacheRegister cacheRegister = mock(CacheRegister.class);

        DeliveryService deliveryService = mock(DeliveryService.class);
        doReturn(true).when(deliveryService).isAvailable();
        when(deliveryService.isLocationAvailable(cacheLocation)).thenReturn(false);
        when(deliveryService.isCoverageZoneOnly()).thenReturn(false);
        when(deliveryService.getRoutingName()).thenReturn("edge");
        when(deliveryService.isDns()).thenReturn(true);

        doReturn(deliveryService).when(trafficRouter).selectDeliveryService(request);
        doReturn(cacheLocation).when(trafficRouter).getCoverageZoneCacheLocation("192.168.34.56", deliveryService, IPVersions.IPV4ONLY);
        doReturn(cacheRegister).when(trafficRouter).getCacheRegister();

        trafficRouter.route(request, track);

        verify(track).setResult(ResultType.MISS);
        verify(track).setResultDetails(ResultDetails.DS_CLIENT_GEO_UNSUPPORTED);
    }
}
