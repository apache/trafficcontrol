package com.comcast.cdn.traffic_control.traffic_router.core.router;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.Matchers.greaterThan;
import static org.mockito.Matchers.any;
import static org.mockito.Matchers.anyBoolean;
import static org.mockito.Matchers.anyString;
import static org.mockito.Mockito.mock;
import static org.powermock.api.mockito.PowerMockito.doCallRealMethod;
import static org.powermock.api.mockito.PowerMockito.doReturn;
import static org.powermock.api.mockito.PowerMockito.spy;
import static org.powermock.api.mockito.PowerMockito.whenNew;

import java.io.File;
import java.io.FileReader;
import java.net.InetAddress;
import java.util.HashMap;
import java.util.HashSet;
import java.util.Iterator;
import java.util.List;
import java.util.Map;
import java.util.Random;
import java.util.Set;

import org.apache.commons.pool.impl.GenericObjectPool;
import org.json.JSONArray;
import org.json.JSONObject;
import org.json.JSONTokener;
import org.junit.Before;
import org.junit.Test;
import org.junit.experimental.categories.Category;
import org.junit.runner.RunWith;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;
import org.springframework.context.ApplicationContext;

import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheLocation;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheRegister;
import com.comcast.cdn.traffic_control.traffic_router.core.dns.ZoneManager;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryService;
import com.comcast.cdn.traffic_control.traffic_router.core.hash.MD5HashFunctionPoolableObjectFactory;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.FederationRegistry;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.NetworkNode;
import com.comcast.cdn.traffic_control.traffic_router.core.request.DNSRequest;
import com.comcast.cdn.traffic_control.traffic_router.core.request.Request;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker.Track;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker.Track.ResultType;
import com.comcast.cdn.traffic_control.traffic_router.core.util.IntegrationTest;
import com.comcast.cdn.traffic_control.traffic_router.core.util.TrafficOpsUtils;
import com.comcast.cdn.traffic_control.traffic_router.geolocation.Geolocation;
import com.comcast.cdn.traffic_control.traffic_router.geolocation.GeolocationService;
import com.google.common.net.InetAddresses;

@Category(IntegrationTest.class)
@RunWith(PowerMockRunner.class)
@PrepareForTest(TrafficRouter.class)
public class DnsRoutePerformanceTest {

    private TrafficRouter trafficRouter;
    private Map<String, Set<String>> hostMap = new HashMap<String, Set<String>>();
    private Set<String> coverageZoneRouted = new HashSet<String>();

    long minimumTPS = Long.parseLong(System.getProperty("minimumTPS", "140"));
    private List<String> names;

    @Before
    public void before() throws Exception {
        CacheRegister cacheRegister = new CacheRegister();

        JSONTokener healthTokener = new JSONTokener(new FileReader("src/test/db/health.json"));
        JSONObject healthObject = new JSONObject(healthTokener);

        JSONTokener jsonTokener = new JSONTokener(new FileReader("src/test/db/cr-config.json"));
        JSONObject configJson = new JSONObject(jsonTokener);
        JSONObject locationsJo = configJson.getJSONObject("edgeLocations");
        final Set<CacheLocation> locations = new HashSet<CacheLocation>(locationsJo.length());
        for (final String loc : JSONObject.getNames(locationsJo)) {
            final JSONObject jo = locationsJo.getJSONObject(loc);
            locations.add(new CacheLocation(loc, jo.optString("zoneId"), new Geolocation(jo.getDouble("latitude"), jo.getDouble("longitude"))));
        }

        names = new DnsNameGenerator().getNames(configJson.getJSONObject("deliveryServices"), configJson.getJSONObject("config"));

        cacheRegister.setConfig(configJson);
        CacheRegisterBuilder.parseDeliveryServiceConfig(configJson.getJSONObject("deliveryServices"), cacheRegister);

        cacheRegister.setConfiguredLocations(locations);
        CacheRegisterBuilder.parseCacheConfig(configJson.getJSONObject("contentServers"), cacheRegister);

        NetworkNode.generateTree(new File("src/test/db/czmap.json"));

        ZoneManager zoneManager = mock(ZoneManager.class);

        MD5HashFunctionPoolableObjectFactory md5factory = new MD5HashFunctionPoolableObjectFactory();
        GenericObjectPool pool = new GenericObjectPool(md5factory);

        whenNew(ZoneManager.class).withArguments(any(TrafficRouter.class), any(StatTracker.class), any(TrafficOpsUtils.class)).thenReturn(zoneManager);

        trafficRouter = new TrafficRouter(cacheRegister, mock(GeolocationService.class), mock(GeolocationService.class),
            pool, mock(StatTracker.class), mock(TrafficOpsUtils.class), mock(FederationRegistry.class));

        trafficRouter.setApplicationContext(mock(ApplicationContext.class));

        trafficRouter = spy(trafficRouter);

        doCallRealMethod().when(trafficRouter).getCoverageZoneCache(anyString(), any(DeliveryService.class));

        doCallRealMethod().when(trafficRouter).selectCache(any(Request.class), any(DeliveryService.class), any(Track.class));
        doCallRealMethod().when(trafficRouter, "selectCache", any(CacheLocation.class), any(DeliveryService.class));
        doCallRealMethod().when(trafficRouter, "getSupportingCaches", any(List.class), any(DeliveryService.class));
        doCallRealMethod().when(trafficRouter).setState(any(JSONObject.class));
        doCallRealMethod().when(trafficRouter).selectDeliveryService(any(Request.class), anyBoolean());
        doReturn(new Geolocation(39.739167, -104.984722)).when(trafficRouter).getLocation(anyString());

        trafficRouter.setState(healthObject);

        JSONObject coverageZoneMap = new JSONObject(new JSONTokener(new FileReader("src/test/db/czmap.json")));
        JSONObject coverageZones = coverageZoneMap.getJSONObject("coverageZones");

        Iterator iterator = coverageZones.keys();
        Map<String, Integer> seen = new HashMap<String, Integer>();

        while (iterator.hasNext()) {
            String coverageZoneName = (String) iterator.next();
            JSONObject coverageZoneJson = coverageZones.getJSONObject(coverageZoneName);
            JSONArray networks = coverageZoneJson.getJSONArray("network");
            Set<String> hosts = hostMap.get(coverageZoneName);

            if (hosts == null) {
                hosts = new HashSet<String>();
            }

            for (int i = 0; i < networks.length(); i++) {
                String network = networks.getString(i).split("/")[0];
                InetAddress ip = InetAddresses.forString(network);
                ip = InetAddresses.increment(ip);
                String ipstr = InetAddresses.toAddrString(ip);
                hosts.add(ipstr);

                if (seen.containsKey(ipstr)) {
                    seen.put(ipstr, seen.get(ipstr)+1);
                } else {
                    seen.put(ipstr, 1);
                }
            }

            final CacheLocation location = cacheRegister.getCacheLocation(coverageZoneName);

            if (location != null || (location == null && coverageZoneJson.has("coordinates"))) {
                coverageZoneRouted.add(coverageZoneName);
            }

            hostMap.put(coverageZoneName, hosts);
        }

        // remove any dups that will cause a GEO result
        for (String ip : seen.keySet()) {
            if (seen.get(ip) > 1) {
                for (String coverageZoneName : hostMap.keySet()) {
                    Set<String> ips = hostMap.get(coverageZoneName);
                    ips.remove(ip);
                }
            }
        }
    }

    @Test
    public void itSupportsMinimalDNSRouteRequestTPS() throws Exception {
        Track track = StatTracker.getTrack();
        DNSRequest dnsRequest = new DNSRequest();

        long before = System.currentTimeMillis();
        int clients = 0;

        // Make it random within the test run but the same random sequence between two test runs
        Random random = new Random(hostMap.keySet().size());
        for (String cacheGroup : hostMap.keySet()) {
            for (String clientIP : hostMap.get(cacheGroup)) {
                dnsRequest.setHostname(names.get(random.nextInt(names.size())));
                dnsRequest.setClientIP(clientIP);
                trafficRouter.route(dnsRequest, track);
                clients++;
            }
        }

        long tps = clients / ((System.currentTimeMillis() - before) / 1000);
        assertThat(tps, greaterThan(minimumTPS));
    }

    @Test
    public void itUsesCoverageZoneWhenPossible() throws Exception {
        Track track = StatTracker.getTrack();
        DNSRequest dnsRequest = new DNSRequest();

        for (String cacheGroup : hostMap.keySet()) {
            for (String clientIP : hostMap.get(cacheGroup)) {
                dnsRequest.setHostname(names.get(0));
                dnsRequest.setClientIP(clientIP);
                trafficRouter.route(dnsRequest, track);

                if (coverageZoneRouted.contains(cacheGroup)) {
                    assertThat("DNS Request for " + dnsRequest.getHostname() + " " + dnsRequest.getType() + ", client ip " + clientIP + " not found in coverage zone even though " + cacheGroup + " is in coverageZoneRouted data" ,track.getResult(), equalTo(ResultType.CZ));
                } else {
                    assertThat(track.getResult(), equalTo(ResultType.GEO));
                }
            }
        }
    }
}
