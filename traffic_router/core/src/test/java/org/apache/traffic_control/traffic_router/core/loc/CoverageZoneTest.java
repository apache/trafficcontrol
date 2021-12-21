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

package org.apache.traffic_control.traffic_router.core.loc;

import org.apache.traffic_control.traffic_router.core.edge.Cache;
import org.apache.traffic_control.traffic_router.core.edge.CacheLocation;
import org.apache.traffic_control.traffic_router.core.edge.CacheLocation.LocalizationMethod;
import org.apache.traffic_control.traffic_router.core.edge.CacheRegister;
import org.apache.traffic_control.traffic_router.core.edge.Node.IPVersions;
import org.apache.traffic_control.traffic_router.core.ds.DeliveryService;
import org.apache.traffic_control.traffic_router.core.router.TrafficRouter;
import org.apache.traffic_control.traffic_router.geolocation.Geolocation;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.api.mockito.PowerMockito;
import org.powermock.core.classloader.annotations.PowerMockIgnore;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;
import org.powermock.reflect.Whitebox;

import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Set;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyList;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

@RunWith(PowerMockRunner.class)
@PrepareForTest(TrafficRouter.class)
@PowerMockIgnore("javax.management.*")
public class CoverageZoneTest {

	private TrafficRouter trafficRouter;

	@Before
	public void before() throws Exception {
		DeliveryService deliveryService = mock(DeliveryService.class);
		when(deliveryService.getId()).thenReturn("delivery-service-1");

		Cache.DeliveryServiceReference deliveryServiceReference = new Cache.DeliveryServiceReference("delivery-service-1", "some.example.com");

		List<Cache.DeliveryServiceReference> deliveryServices = new ArrayList<Cache.DeliveryServiceReference>();
		deliveryServices.add(deliveryServiceReference);

		Geolocation testLocation = new Geolocation(40.0, -101);
		Geolocation farEastLocation = new Geolocation(40.0, -101.5);
		Geolocation eastLocation = new Geolocation(40.0, -100);
		Geolocation westLocation = new Geolocation(40.0, -105);

		Cache farEastCache1 = new Cache("far-east-cache-1", "hashid", 1);
		farEastCache1.setIsAvailable(true);
		Set<LocalizationMethod> lms = new HashSet<>();
		lms.add(LocalizationMethod.GEO);
		CacheLocation farEastCacheGroup = new CacheLocation("far-east-cache-group", farEastLocation, lms);
		farEastCacheGroup.addCache(farEastCache1);
		farEastCache1.setDeliveryServices(deliveryServices);

		Cache eastCache1 = new Cache("east-cache-1", "hashid", 1);
		eastCache1.setIsAvailable(true);
		CacheLocation eastCacheGroup = new CacheLocation("east-cache-group", eastLocation);
		eastCacheGroup.addCache(eastCache1);

		Cache westCache1 = new Cache("west-cache-1", "hashid", 1);
		westCache1.setIsAvailable(true);
		westCache1.setDeliveryServices(deliveryServices);

		CacheLocation westCacheGroup = new CacheLocation("west-cache-group", westLocation);
		westCacheGroup.addCache(westCache1);

		List<CacheLocation> cacheGroups = new ArrayList<CacheLocation>();
		cacheGroups.add(farEastCacheGroup);
		cacheGroups.add(eastCacheGroup);
		cacheGroups.add(westCacheGroup);

		NetworkNode eastNetworkNode = new NetworkNode("12.23.34.0/24", "east-cache-group", testLocation);

		CacheRegister cacheRegister = mock(CacheRegister.class);

		when(cacheRegister.getCacheLocationById("east-cache-group")).thenReturn(eastCacheGroup);

		when(cacheRegister.filterAvailableCacheLocations("delivery-service-1")).thenReturn(cacheGroups);
		when(cacheRegister.getDeliveryService("delivery-service-1")).thenReturn(deliveryService);

		trafficRouter = PowerMockito.mock(TrafficRouter.class);
		Whitebox.setInternalState(trafficRouter, "cacheRegister", cacheRegister);
		when(trafficRouter.getCoverageZoneCacheLocation("12.23.34.45", "delivery-service-1", IPVersions.IPV4ONLY)).thenCallRealMethod();
		when(trafficRouter.getCoverageZoneCacheLocation("12.23.34.45", "delivery-service-1", false, null, IPVersions.IPV4ONLY)).thenCallRealMethod();
		when(trafficRouter.getCacheRegister()).thenReturn(cacheRegister);
		when(trafficRouter.orderLocations(anyList(),any(Geolocation.class))).thenCallRealMethod();
		when(trafficRouter.getSupportingCaches(anyList(), eq(deliveryService), any(IPVersions.class))).thenCallRealMethod();
		when(trafficRouter.filterEnabledLocations(anyList(), any(CacheLocation.LocalizationMethod.class))).thenCallRealMethod();
		PowerMockito.when(trafficRouter, "getNetworkNode", "12.23.34.45").thenReturn(eastNetworkNode);
		PowerMockito.when(trafficRouter, "getClosestCacheLocation", anyList(), any(), any(), any()).thenCallRealMethod();
	}

	@Test
	public void trafficRouterReturnsNearestCacheGroupForDeliveryService() throws Exception {
		CacheLocation cacheLocation = trafficRouter.getCoverageZoneCacheLocation("12.23.34.45", "delivery-service-1", IPVersions.IPV4ONLY);
		assertThat(cacheLocation.getId(), equalTo("west-cache-group"));
		// NOTE: far-east-cache-group is actually closer to the client but isn't enabled for CZ-localization and must be filtered out
	}
}
