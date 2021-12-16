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

import org.apache.traffic_control.traffic_router.core.edge.Cache;
import org.apache.traffic_control.traffic_router.core.ds.DeliveryService;
import org.apache.traffic_control.traffic_router.core.ds.SteeringResult;
import org.apache.traffic_control.traffic_router.core.ds.SteeringTarget;
import org.apache.traffic_control.traffic_router.geolocation.Geolocation;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.core.classloader.annotations.PowerMockIgnore;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import static org.junit.Assert.assertEquals;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyList;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.spy;
import static org.mockito.Mockito.doCallRealMethod;
import static org.mockito.Mockito.never;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

@RunWith(PowerMockRunner.class)
@PrepareForTest({Collections.class})
@PowerMockIgnore("javax.management.*")
public class GeoSortSteeringResultsTest {

    private TrafficRouter trafficRouter;
    private List<SteeringResult> steeringResults;
    private Geolocation clientLocation;
    private DeliveryService deliveryService;

    @Before
    public void before() {
        trafficRouter = mock(TrafficRouter.class);
        steeringResults = spy(new ArrayList<>());
        clientLocation = new Geolocation(47.0, -122.0);
        deliveryService = mock(DeliveryService.class);
        doCallRealMethod().when(trafficRouter).geoSortSteeringResults(anyList(), anyString(), any(DeliveryService.class));
        when(trafficRouter.getClientLocationByCoverageZoneOrGeo(anyString(), any(DeliveryService.class))).thenReturn(clientLocation);
    }

    @Test
    public void testNullClientIP() {
        trafficRouter.geoSortSteeringResults(steeringResults, null, deliveryService);
        verify(trafficRouter, never()).getClientLocationByCoverageZoneOrGeo(null, deliveryService);
    }

    @Test
    public void testEmptyClientIP() {
        trafficRouter.geoSortSteeringResults(steeringResults, "", deliveryService);
        verify(trafficRouter, never()).getClientLocationByCoverageZoneOrGeo("", deliveryService);
    }

    @Test
    public void testNoSteeringTargetsHaveGeolocations() {
        steeringResults.add(new SteeringResult(new SteeringTarget(), deliveryService));
        trafficRouter.geoSortSteeringResults(steeringResults, "::1", deliveryService);
        verify(trafficRouter, never()).getClientLocationByCoverageZoneOrGeo("::1", deliveryService);
    }

    @Test
    public void testClientGeolocationIsNull() {
        SteeringTarget steeringTarget = new SteeringTarget();
        steeringTarget.setGeolocation(clientLocation);
        steeringResults.add(new SteeringResult(steeringTarget, deliveryService));
        when(trafficRouter.getClientLocationByCoverageZoneOrGeo(anyString(), any(DeliveryService.class))).thenReturn(null);

        trafficRouter.geoSortSteeringResults(steeringResults, "::1", deliveryService);

        verify(steeringResults, never()).sort(any());
    }

    @Test
    public void testGeoSortingMixedWithNonGeoTargets() {
        Cache cache = new Cache("fake-id", "fake-hash-id",1, clientLocation);
        SteeringTarget target;

        target = new SteeringTarget();
        target.setOrder(-1);
        SteeringResult resultNoGeoNegativeOrder = new SteeringResult(target, deliveryService);
        resultNoGeoNegativeOrder.setCache(cache);
        steeringResults.add(resultNoGeoNegativeOrder);

        target = new SteeringTarget();
        target.setOrder(1);
        SteeringResult resultNoGeoPositiveOrder = new SteeringResult(target, deliveryService);
        resultNoGeoPositiveOrder.setCache(cache);
        steeringResults.add(resultNoGeoPositiveOrder);

        target = new SteeringTarget();
        target.setOrder(0);
        SteeringResult resultNoGeoZeroOrder = new SteeringResult(target, deliveryService);
        resultNoGeoZeroOrder.setCache(cache);
        steeringResults.add(resultNoGeoZeroOrder);

        target = new SteeringTarget();
        target.setGeolocation(clientLocation);
        target.setOrder(0);
        SteeringResult resultGeo = new SteeringResult(target, deliveryService);
        resultGeo.setCache(cache);
        steeringResults.add(resultGeo);

        trafficRouter.geoSortSteeringResults(steeringResults, "::1", deliveryService);

        assertEquals(resultNoGeoNegativeOrder, steeringResults.get(0));
        assertEquals(resultGeo, steeringResults.get(1));
        assertEquals(resultNoGeoZeroOrder, steeringResults.get(2));
        assertEquals(resultNoGeoPositiveOrder, steeringResults.get(3));
    }

}

