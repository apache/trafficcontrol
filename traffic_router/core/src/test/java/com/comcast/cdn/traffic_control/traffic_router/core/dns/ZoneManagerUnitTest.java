/*
t *
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

package com.comcast.cdn.traffic_control.traffic_router.core.dns;

import com.comcast.cdn.traffic_control.traffic_router.core.edge.CacheRegister;
import com.comcast.cdn.traffic_control.traffic_router.core.edge.InetRecord;
import com.comcast.cdn.traffic_control.traffic_router.core.config.ConfigHandler;
import com.comcast.cdn.traffic_control.traffic_router.core.config.SnapshotEventsProcessor;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryService;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.SteeringWatcher;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.AnonymousIpDatabaseService;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.FederationRegistry;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.FederationsWatcher;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.MaxmindGeolocationService;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker.Track.ResultType;
import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouter;
import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouterManager;
import com.comcast.cdn.traffic_control.traffic_router.core.util.JsonUtils;
import com.comcast.cdn.traffic_control.traffic_router.core.util.TrafficOpsUtils;
import com.comcast.cdn.traffic_control.traffic_router.geolocation.GeolocationService;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.google.common.cache.LoadingCache;
import org.apache.commons.io.IOUtils;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.Mockito;
import org.powermock.api.mockito.PowerMockito;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;
import org.powermock.reflect.Whitebox;
import org.xbill.DNS.ARecord;
import org.xbill.DNS.Name;
import org.xbill.DNS.Record;
import org.xbill.DNS.SetResponse;
import org.xbill.DNS.Type;
import org.xbill.DNS.Zone;

import java.io.File;
import java.io.InputStream;
import java.net.InetAddress;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.lessThan;
import static org.hamcrest.Matchers.notNullValue;
import static org.junit.Assert.*;
import static org.mockito.Matchers.*;
import static org.mockito.Mockito.verify;
import static org.powermock.api.mockito.PowerMockito.*;

@RunWith(PowerMockRunner.class)
@PrepareForTest({ConfigHandler.class, CacheRegister.class, ZoneManager.class, SignatureManager.class,
		TrafficRouterManager.class, TrafficRouter.class })
public class ZoneManagerUnitTest {
    ZoneManager zoneManager;
    SignatureManager signatureManager;
    TrafficRouter trafficRouter;
    CacheRegister cacheRegister;
    Map<String, DeliveryService> changes;
	LoadingCache<ZoneKey, Zone> dynamicZoneCache;
	LoadingCache<ZoneKey, Zone> zoneCache;
	DeliveryService deliveryService;
	JsonNode newDsSnapJo = null;
	JsonNode updateJo = null;
	JsonNode baselineJo = null;
	TrafficRouterManager trafficRouterManager;

    @Before
    public void before() throws Exception {
    	try {
		    String resourcePath = "unit/ExistingConfig.json";
		    InputStream inputStream = getClass().getClassLoader().getResourceAsStream(resourcePath);
		    if (inputStream == null) {
			    fail("Could not find file '" + resourcePath + "' needed for test from the current classpath as a resource!");
		    }
		    String baseDb = IOUtils.toString(inputStream);

		    resourcePath = "unit/UpdateDsSnap.json";
		    inputStream = getClass().getClassLoader().getResourceAsStream(resourcePath);
		    if (inputStream == null) {
			    fail("Could not find file '" + resourcePath + "' needed for test from the current classpath as a resource!");
		    }
		    String updateDb = IOUtils.toString(inputStream);

		    resourcePath = "unit/NewDsSnap.json";
		    inputStream = getClass().getClassLoader().getResourceAsStream(resourcePath);
		    if (inputStream == null) {
			    fail("Could not find file '" + resourcePath + "' needed for test from the current classpath as a resource!");
		    }
		    String newDsSnapDb = IOUtils.toString(inputStream);

		    final ObjectMapper mapper = new ObjectMapper();
		    assertThat(newDsSnapDb, notNullValue());
		    assertThat(baseDb, notNullValue());
		    assertThat(updateDb, notNullValue());

		    newDsSnapJo = mapper.readTree(newDsSnapDb);
		    assertThat(newDsSnapJo, notNullValue());
		    updateJo = mapper.readTree(updateDb);
		    assertThat(updateJo, notNullValue());
		    baselineJo = mapper.readTree(baseDb);
		    assertThat(baselineJo, notNullValue());

		    GeolocationService geolocationService = new MaxmindGeolocationService();
		    AnonymousIpDatabaseService anonymousIpDatabaseService = new AnonymousIpDatabaseService();
		    FederationRegistry federationRegistry = new FederationRegistry();
		    TrafficOpsUtils trafficOpsUtils = new TrafficOpsUtils();
		    trafficRouterManager = PowerMockito.spy(new TrafficRouterManager());
		    trafficRouterManager.setAnonymousIpService(anonymousIpDatabaseService);
		    trafficRouterManager.setGeolocationService(geolocationService);
		    trafficRouterManager.setGeolocationService6(geolocationService);
		    trafficRouterManager.setFederationRegistry(federationRegistry);
		    trafficRouterManager.setTrafficOpsUtils(trafficOpsUtils);
		    trafficRouterManager.setNameServer(new NameServer());
		    SnapshotEventsProcessor snapshotEventsProcessor = SnapshotEventsProcessor.diffCrConfigs(baselineJo, null);
		    ConfigHandler configHandler = PowerMockito.spy(new ConfigHandler());
		    StatTracker statTracker = new StatTracker();
		    ZoneManager.setZoneDirectory(new File("src/test/resources/unit/zonemanager"));
		    configHandler.setTrafficRouterManager(trafficRouterManager);
		    configHandler.setStatTracker(statTracker);
		    configHandler.setFederationsWatcher(new FederationsWatcher());
		    configHandler.setSteeringWatcher(new SteeringWatcher());
		    final Map<String, DeliveryService> deliveryServiceMap = snapshotEventsProcessor.getCreationEvents();
		    final JsonNode config = JsonUtils.getJsonNode(baselineJo, ConfigHandler.CONFIG_KEY);
		    final JsonNode stats = JsonUtils.getJsonNode(baselineJo, "stats");
		    cacheRegister = spy(new CacheRegister());
		    cacheRegister.setTrafficRouters(JsonUtils.getJsonNode(baselineJo, "contentRouters"));
		    cacheRegister.setConfig(config);
		    cacheRegister.setStats(stats);
		    Whitebox.invokeMethod(configHandler, "parseCertificatesConfig", config);
		    final List<DeliveryService> deliveryServices = new ArrayList<>();

		    if (deliveryServiceMap != null && !deliveryServiceMap.values().isEmpty()) {
			    deliveryServices.addAll(deliveryServiceMap.values());
		    }

		    Whitebox.invokeMethod(configHandler, "parseDeliveryServiceMatchSets", deliveryServiceMap, cacheRegister);
		    Whitebox.invokeMethod(configHandler, "parseLocationConfig", JsonUtils.getJsonNode(baselineJo, "edgeLocations"),
			    cacheRegister);
		    Whitebox.invokeMethod(configHandler, "parseCacheConfig", JsonUtils.getJsonNode(baselineJo,
			    ConfigHandler.CONTENT_SERVERS_KEY), cacheRegister);
		    Whitebox.invokeMethod(configHandler, "parseMonitorConfig", JsonUtils.getJsonNode(baselineJo, "monitors"));
		    FederationsWatcher federationsWatcher = new FederationsWatcher();
		    federationsWatcher.configure(config);
		    configHandler.setFederationsWatcher(federationsWatcher);
		    SteeringWatcher steeringWatcher = new SteeringWatcher();
		    steeringWatcher.configure(config);
		    configHandler.setSteeringWatcher(steeringWatcher);
		    trafficRouter = PowerMockito.spy( new TrafficRouter(cacheRegister, geolocationService, geolocationService,
		    anonymousIpDatabaseService, federationRegistry, trafficRouterManager));
		    zoneManager =
				    PowerMockito.spy(ZoneManager.initialInstance(trafficRouter,statTracker,trafficOpsUtils,trafficRouterManager));
		    Whitebox.setInternalState(trafficRouter, "zoneManager", zoneManager);
		    Whitebox.setInternalState(trafficRouterManager, "trafficRouter", trafficRouter);
		    trafficRouterManager.getNameServer().setEcsEnable(JsonUtils.optBoolean(config, "ecsEnable", false));
		    trafficRouterManager.getTrafficRouter().configurationChanged();
		    assertNotNull(zoneManager);
		    assertEquals(zoneManager, trafficRouterManager.getTrafficRouter().getZoneManager());
		    zoneCache = Whitebox.getInternalState(ZoneManager.class, "zoneCache");
		    dynamicZoneCache = Whitebox.getInternalState(ZoneManager.class, "dynamicZoneCache");
		    signatureManager = Whitebox.getInternalState(ZoneManager.class, "signatureManager");
		    assertNotNull(dynamicZoneCache);
		    assertNotNull(zoneCache);
		    assertNotNull(signatureManager);
		    signatureManager = spy(new SignatureManager(zoneManager,cacheRegister, trafficOpsUtils,
				    trafficRouterManager));
		    Whitebox.setInternalState(ZoneManager.class, "signatureManager", signatureManager);
	    }
	    catch (Exception ex)
	    {
	    	ex.printStackTrace();
	    	fail();
	    }
    }

    @Test
    public void itMarksResultTypeAndLocationInDNSAccessRecord() throws Exception {
        final Name qname = Name.fromString("edge.www.google.com.");
        final InetAddress client = InetAddress.getByName("192.168.56.78");

        SetResponse setResponse = mock(SetResponse.class);
        when(setResponse.isSuccessful()).thenReturn(false);

        Zone zone = mock(Zone.class);
        when(zone.findRecords(any(Name.class), anyInt())).thenReturn(setResponse);
	    when(zone.getOrigin()).thenReturn(new Name(qname, 1));

        DNSAccessRecord.Builder builder = new DNSAccessRecord.Builder(1L, client);
        builder = spy(builder);

        doReturn(zone).when(zoneManager).getZone(qname, Type.A);
        doCallRealMethod().when(zoneManager).getZone(qname, Type.A, client, false, builder);

        zoneManager.getZone(qname, Type.A, client, false, builder);
        verify(builder).resultType(any(ResultType.class));
        verify(builder).resultLocation(null);
    }

	@Test
	public void snapshotReplacesZoneCaches() {
		try {
			SnapshotEventsProcessor snapshotEventsProcessor = mock(SnapshotEventsProcessor.class);
			when(snapshotEventsProcessor.getChangeEvents()).thenReturn(new HashMap<>());
			zoneManager = PowerMockito.spy(ZoneManager
					.snapshotInstance(trafficRouter,  new StatTracker(), snapshotEventsProcessor ));
			assertNotEquals("expected snapshotInstance to replace the dynamicZoneCache, but it is unchanged",
					dynamicZoneCache, Whitebox.getInternalState( ZoneManager.class, "dynamicZoneCache" ));
			assertNotEquals("expected snaphshotInstance to replace the zoneCache, but it is unchanged",
					zoneCache, Whitebox.getInternalState( ZoneManager.class, "zoneCache" ));
			assertEquals("zoneCache size == but it is: "+ ((LoadingCache<ZoneKey, Zone>)Whitebox.getInternalState( ZoneManager.class, "zoneCache" )).toString(),
					zoneCache.size(),
					((LoadingCache<ZoneKey, Zone>)Whitebox.getInternalState( ZoneManager.class, "zoneCache" )).size());

		} catch (Exception ioe) {
			ioe.printStackTrace();
			fail("In snapshotReplacesZoneCaches - " + ioe.toString());
		}
	}

	@Test
	public void processDsChanges() {
		try {
			SnapshotEventsProcessor snapshotEventsProcessor = mock(SnapshotEventsProcessor.class);
			deliveryService = mock(DeliveryService.class);
			when(deliveryService.getId()).thenReturn("MockDs");
			when(deliveryService.getDomain()).thenReturn("mockds.mockcdn.moc");
			when(deliveryService.isDns()).thenReturn(true);
			changes = new HashMap<>();
			changes.put("MockDs", deliveryService);
			when(cacheRegister.getDeliveryService(anyString())).thenReturn(deliveryService);
			when(cacheRegister.getDeliveryServices()).thenReturn(changes);
			Set<String> routingNames = new HashSet<>();
			Whitebox.setInternalState(ZoneManager.class, "dnsRoutingNames", routingNames);
			dynamicZoneCache = ZoneManager.createZoneCache(ZoneManager.ZoneCacheType.DYNAMIC);
			Whitebox.setInternalState(ZoneManager.class, "dynamicZoneCache", dynamicZoneCache);
			zoneCache = ZoneManager.createZoneCache(ZoneManager.ZoneCacheType.STATIC);
			Whitebox.setInternalState(ZoneManager.class, "zoneCache", zoneCache);
			when(snapshotEventsProcessor.getChangeEvents()).thenReturn(changes);
			zoneManager.processChangeEvents(snapshotEventsProcessor);
			assertNotEquals("expected processDsChanges to replace the dynamicZoneCache, but it is unchanged",
					dynamicZoneCache, Whitebox.getInternalState(ZoneManager.class, "dynamicZoneCache"));
			assertNotEquals("expected processDsChanges to replace the zoneCache, but it is unchanged",
					zoneCache, Whitebox.getInternalState(ZoneManager.class, "zoneCache"));
		} catch (Exception ioe) {
			ioe.printStackTrace();
			fail("In processDsChanges - " + ioe.toString());
		}
	}

	@Test
	public void verifyChangedZonesAppened() {
		try {
			assertTrue(zoneCache != null);
			assertTrue(zoneCache.size() > 0);
			LoadingCache<ZoneKey, Zone> prevZc = zoneCache;
			SnapshotEventsProcessor snapshotEventsProcessor = mock(SnapshotEventsProcessor.class);
			deliveryService = mock(DeliveryService.class);
			when(deliveryService.getId()).thenReturn("http-and-https-test");
			when(deliveryService.getDomain()).thenReturn("thecdn.example.com");
			when(deliveryService.isDns()).thenReturn(true);
			when(deliveryService.getLocationLimit()).thenReturn(50);
			changes = new HashMap<>();
			changes.put("http-and-https-test", deliveryService);
			when(cacheRegister.getDeliveryService(anyString())).thenReturn(deliveryService);
			when(cacheRegister.getDeliveryServices()).thenReturn(changes);
			Set<String> routingNames = new HashSet<>();
			Whitebox.setInternalState(ZoneManager.class, "dnsRoutingNames", routingNames);
			when(snapshotEventsProcessor.getChangeEvents()).thenReturn(changes);
			zoneManager.processChangeEvents(snapshotEventsProcessor);
			assertNotEquals("expected processDsChanges to replace the dynamicZoneCache, but it is unchanged",
					dynamicZoneCache, Whitebox.getInternalState(ZoneManager.class, "dynamicZoneCache"));
			assertNotEquals("expected processDsChanges to replace the zoneCache, but it is unchanged",
					zoneCache, Whitebox.getInternalState(ZoneManager.class, "zoneCache"));
			int zize = ((LoadingCache<ZoneKey, Zone>)Whitebox.getInternalState(ZoneManager.class,
					"zoneCache")).asMap().size();
			assertThat( prevZc.asMap().size(), lessThan(zize));
		} catch (Exception ioe) {
			ioe.printStackTrace();
			fail("In processDsChanges - " + ioe.toString());
		}
	}
	@Test
	public void resolve() {
		try {
			Name keyName = new Name("mockds.mockcdn.moc");
			Name hostName = new Name("edge.mockds.mockcdn.com");
			Record dnsRecord = mock(ARecord.class);
			when(dnsRecord.getType()).thenReturn(Type.NS);
			when(dnsRecord.getName()).thenReturn(hostName);
			InetAddress resolved = InetAddress.getByAddress(new byte[]{123, 4, 5, 6});
			when(((ARecord) dnsRecord).getAddress()).thenReturn(resolved);
			when(((ARecord) dnsRecord).getTTL()).thenReturn(60l);
			List<Record> records = new ArrayList<>();
			records.add(dnsRecord);
			Zone retZone = PowerMockito.mock(Zone.class);
			when(retZone.getOrigin()).thenReturn(keyName);
			PowerMockito.doReturn(retZone).when(zoneManager, "getZone", any(Name.class));
			List<InetRecord> lookupResult = new ArrayList<>();
			InetRecord resultrec = new InetRecord(resolved, 60l);
			lookupResult.add(resultrec);
			PowerMockito.doReturn(lookupResult).when(zoneManager, "lookup", any(Name.class), any(Zone.class),
					eq(Type.A));
			List<InetRecord> results = Whitebox.invokeMethod(zoneManager, "resolve", "edge.mockds.mocdn.moc");
			assertNotNull(results);
			assertTrue(results.size() > 0);
			assertEquals(((ARecord) dnsRecord).getAddress().toString(), results.get(0).getAddress().toString());
			when(zoneManager.getZone(any(Name.class))).thenReturn(retZone);
			results = Whitebox.invokeMethod(zoneManager, "resolve", "edge.mockds.mocdn.moc");
			assertNotNull(results);
			assertEquals("/123.4.5.6", results.get(0).getAddress().toString());
		} catch (Exception ioe) {
			ioe.printStackTrace();
			fail("In resolve - " + ioe.toString());
		}
	}

	@Test
	public void updateZoneCacheWithDsChanges() {
		try {
			assertEquals(5, dynamicZoneCache.size());
			assertEquals(19, zoneCache.size());
			// Modify 2 add 1
			SnapshotEventsProcessor snapshotEventsProcessor = SnapshotEventsProcessor.diffCrConfigs(newDsSnapJo,
					baselineJo);
			zoneManager.updateZoneCache(snapshotEventsProcessor.getChangeEvents());
			verify(signatureManager, Mockito.times(3)).generateZoneKey(any(), any());
			zoneCache = Whitebox.getInternalState(ZoneManager.class, "zoneCache");
			dynamicZoneCache = Whitebox.getInternalState(ZoneManager.class, "dynamicZoneCache");
			assertEquals(5, dynamicZoneCache.size());
			assertEquals(20, zoneCache.size());

			// change 15 and delete one
			snapshotEventsProcessor = SnapshotEventsProcessor.diffCrConfigs(updateJo,
					newDsSnapJo);
			zoneManager.updateZoneCache(snapshotEventsProcessor.getChangeEvents());
			verify(signatureManager, Mockito.times(3+16)).generateZoneKey(any(), any());
			zoneCache = Whitebox.getInternalState(ZoneManager.class, "zoneCache");
			dynamicZoneCache = Whitebox.getInternalState(ZoneManager.class, "dynamicZoneCache");
			assertEquals(5, dynamicZoneCache.size());
			assertEquals(20, zoneCache.size());

			// change none
			snapshotEventsProcessor = SnapshotEventsProcessor.diffCrConfigs(updateJo,
					updateJo);
			zoneManager.updateZoneCache(snapshotEventsProcessor.getChangeEvents());
			verify(signatureManager, Mockito.times(3+16+0)).generateZoneKey(any(), any());
			zoneCache = Whitebox.getInternalState(ZoneManager.class, "zoneCache");
			dynamicZoneCache = Whitebox.getInternalState(ZoneManager.class, "dynamicZoneCache");
			assertEquals(5, dynamicZoneCache.size());
			assertEquals(20, zoneCache.size());
		}
		catch (Exception exp)
		{
			exp.printStackTrace();
			fail("In updateZoneCacheWithKeyList - " + exp.toString());
		}
	}
}
