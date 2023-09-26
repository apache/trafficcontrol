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

package org.apache.traffic_control.traffic_router.core.dns;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.node.ArrayNode;
import com.fasterxml.jackson.databind.node.JsonNodeFactory;
import com.fasterxml.jackson.databind.node.ObjectNode;
import com.google.common.cache.CacheBuilder;
import com.google.common.cache.CacheLoader;
import com.google.common.cache.LoadingCache;
import org.apache.traffic_control.traffic_router.core.ds.DeliveryService;
import org.apache.traffic_control.traffic_router.core.edge.CacheRegister;
import org.apache.traffic_control.traffic_router.core.edge.InetRecord;
import org.apache.traffic_control.traffic_router.core.request.DNSRequest;
import org.apache.traffic_control.traffic_router.core.router.DNSRouteResult;
import org.apache.traffic_control.traffic_router.core.router.StatTracker;
import org.apache.traffic_control.traffic_router.core.router.StatTracker.Track.ResultType;
import org.apache.traffic_control.traffic_router.core.router.TrafficRouterManager;
import org.apache.traffic_control.traffic_router.core.util.TrafficOpsUtils;
import org.apache.traffic_control.traffic_router.core.router.TrafficRouter;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.stubbing.Answer;
import org.powermock.api.mockito.PowerMockito;
import org.powermock.core.classloader.annotations.PowerMockIgnore;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;
import org.xbill.DNS.ARecord;
import org.xbill.DNS.DClass;
import org.xbill.DNS.NSECRecord;
import org.xbill.DNS.Name;
import org.xbill.DNS.RRSIGRecord;
import org.xbill.DNS.Record;
import org.xbill.DNS.SOARecord;
import org.xbill.DNS.SetResponse;
import org.xbill.DNS.Type;
import org.xbill.DNS.Zone;
import org.xbill.DNS.NSRecord;
import org.xbill.DNS.CNAMERecord;

import java.io.File;
import java.net.InetAddress;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Date;
import java.util.List;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyInt;
import static org.mockito.Mockito.*;
import static org.powermock.api.mockito.PowerMockito.doCallRealMethod;
import static org.powermock.api.mockito.PowerMockito.whenNew;

@RunWith(PowerMockRunner.class)
@PrepareForTest({ZoneManager.class, SignatureManager.class, InetAddress.class})
@PowerMockIgnore("javax.management.*")
public class ZoneManagerUnitTest {
    ZoneManager zoneManager;
    TrafficRouter trafficRouter;
    SignatureManager signatureManager;
    CacheRegister cacheRegister;
    @Before
    public void before() throws Exception {
        trafficRouter = mock(TrafficRouter.class);
        cacheRegister = mock(CacheRegister.class);
        when(trafficRouter.getCacheRegister()).thenReturn(cacheRegister);

        PowerMockito.spy(ZoneManager.class);
        PowerMockito.stub(PowerMockito.method(ZoneManager.class, "initTopLevelDomain")).toReturn(null);
        PowerMockito.stub(PowerMockito.method(ZoneManager.class, "initZoneCache")).toReturn(null);

        signatureManager = PowerMockito.mock(SignatureManager.class);
        whenNew(SignatureManager.class).withArguments(any(ZoneManager.class), any(CacheRegister.class), any(TrafficOpsUtils.class), any(TrafficRouterManager.class)).thenReturn(signatureManager);

        zoneManager = spy(new ZoneManager(trafficRouter, new StatTracker(), null, mock(TrafficRouterManager.class)));

    }

    @Test
    public void testNegativeCachingTTLGetterAndSetter() throws Exception {
        final File file = new File("src/test/resources/publish/CrConfig5.json");
        final ObjectMapper mapper = new ObjectMapper();
        final JsonNode jo = mapper.readTree(file);
        zoneManager.setNegativeCachingTTL(jo);
        assertThat(zoneManager.getNegativeCachingTTL(), equalTo(1200L));
    }

    @Test
    public void testGetLocalTRHostnameUsesSingleTRHostname() throws Exception {
        JsonNode trs = new ObjectMapper().readTree("{\"tr-01\": {}}");
        when(cacheRegister.getTrafficRouters()).thenReturn(trs);
        InetAddress localhost = mock(InetAddress.class);
        when(localhost.getHostName()).thenReturn("real-local-hostname");
        PowerMockito.mockStatic(InetAddress.class);
        when(InetAddress.getLocalHost()).thenReturn(localhost);
        String actual = ZoneManager.getTRLocalHostname(zoneManager.getTrafficRouter());
        assertThat("hostname of TR server in the CRConfig is returned when there is only one TR in the CRConfig", actual, equalTo("tr-01"));
    }

    @Test
    public void testGetLocalTRHostnameUsesRealLocalHostNameIfMultipleTR() throws Exception {
        JsonNode trs = new ObjectMapper().readTree("{\"tr-01\": {}, \"tr-02\":  {}}");
        when(cacheRegister.getTrafficRouters()).thenReturn(trs);
        InetAddress localhost = mock(InetAddress.class);
        when(localhost.getHostName()).thenReturn("real-local-hostname");
        PowerMockito.mockStatic(InetAddress.class);
        when(InetAddress.getLocalHost()).thenReturn(localhost);
        String actual = ZoneManager.getTRLocalHostname(zoneManager.getTrafficRouter());
        assertThat("real local hostname of TR server is returned when there are multiple TRs in the CRConfig", actual, equalTo("real-local-hostname"));
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
    public void itGetsCorrectNSECRecordFromStaticAndDynamicZones() throws Exception {

        final Name qname = Name.fromString("dns1.example.com.");
        final InetAddress client = InetAddress.getByName("192.168.56.78");

        DNSAccessRecord.Builder builder = new DNSAccessRecord.Builder(1L, client);
        builder = spy(builder);

        Name m_an, m_host, m_admin;
        m_an = Name.fromString("dns1.example.com.");
        m_host = Name.fromString("dns1.example.com.");
        m_admin = Name.fromString("admin.example.com.");
        Record ar;
        NSRecord ns;
        NSECRecord nsec;
        ar = new SOARecord(m_an, DClass.IN, 0x13A8,
                m_host, m_admin, 0xABCDEF12L, 0xCDEF1234L,
                0xEF123456L, 0x12345678L, 0x3456789AL);

        ns = new NSRecord(m_an, DClass.IN, 12345L, m_an);
        nsec = new NSECRecord(m_an, DClass.IN, 12345L, new Name("foobar.dns1.example.com."), new int[]{1});
        Record[] records = new Record[] {ar, ns, nsec};
        m_an = Name.fromString("dns1.example.com.");
        Zone zone = new Zone(m_an, records);
        // static zone
        doReturn(zone).when(zoneManager).getZone(qname, Type.NSEC);

        DNSRouteResult dnsRouteResult = new DNSRouteResult();
        ObjectNode node = JsonNodeFactory.instance.objectNode();
        ArrayNode domainNode = node.putArray("domains");
        domainNode.add("example.com");
        node.put("routingName","edge");
        node.put("coverageZoneOnly", false);
        DeliveryService ds1 = new DeliveryService("ds1", node);


        dnsRouteResult.setDeliveryService(ds1);
        InetRecord address = new InetRecord("cdn-tr.dns1.example.com.", 12345L);
        List<InetRecord> list = new ArrayList<>();
        list.add(address);
        dnsRouteResult.setAddresses(list);

        Record cnameRecord = new CNAMERecord(new Name("dns1.example.com."), DClass.IN, 12345L, new Name("cdn-tr.dns1.example.com."));
        Record nsecRecord = new NSECRecord(new Name("edge.dns1.example.com."), DClass.IN, 12345L, new Name("foobar.dns1.example.com."), new int[]{1});

        // Add records for dynamic zones
        Record[] recordArray = new Record[]{cnameRecord, ar, nsecRecord, ns};
        List<Record> recordList = Arrays.asList(recordArray);
        Zone dynamicZone = new Zone(new Name("dns1.example.com."), recordArray);

        CacheLoader<ZoneKey, Zone> loader;
        loader = new CacheLoader<>() {
            @Override
            public Zone load(ZoneKey zoneKey) {
                return dynamicZone;
            }

        };
        loader.load(new ZoneKey(Name.fromString("dns1.example.com."), Arrays.asList(records)));
        LoadingCache<ZoneKey, Zone> dynamicZoneCache = CacheBuilder.newBuilder().build(loader);

        // stub calls for signatureManager, dynamicZoneCache and generateDynamicZoneKey
        when(ZoneManager.getDynamicZoneCache()).thenReturn(dynamicZoneCache);
        ZoneKey zk = new ZoneKey(Name.fromString("dns1.example.com."), recordList);
        dynamicZoneCache.put(zk, dynamicZone);

        when(ZoneManager.getSignatureManager()).thenReturn(signatureManager);
        Answer<ZoneKey> currentTimeAnswer = invocation -> zk;

        when(ZoneManager.getSignatureManager().generateDynamicZoneKey(
                eq(Name.fromString("dns1.example.com.")),
                anyList(),
                eq(true))).
                then(currentTimeAnswer);
        when(trafficRouter.isEdgeDNSRouting()).thenReturn(true);
        when(trafficRouter.route(any(DNSRequest.class), any(StatTracker.Track.class))).thenReturn(dnsRouteResult);
        Zone resultZone = zoneManager.getZone(qname, Type.NSEC, client, true, builder);
        // make sure the function gets called with the correct records as expected
        verify(ZoneManager.getSignatureManager()).generateDynamicZoneKey( eq(Name.fromString("dns1.example.com.")),
                argThat(t -> t.containsAll(Arrays.asList(nsecRecord, ns, ar))),
                eq(true));
        SetResponse setResponse = resultZone.findRecords(new Name(ds1.getRoutingName() + "." + "dns1.example.com."), Type.NSEC);
        assertThat(setResponse.isNXDOMAIN(), equalTo(false));
    }

    @Test
    public void testZonesAreEqual() throws java.net.UnknownHostException, org.xbill.DNS.TextParseException {
        class TestCase {
            String reason;
            Record[] r1;
            Record[] r2;
            boolean expected;

            TestCase(String r, Record[] a, Record[] b, boolean e) {
                reason = r;
                r1 = a;
                r2 = b;
                expected = e;
            }
        }

        final TestCase[] testCases = {
                new TestCase("empty lists are equal", new Record[]{}, new Record[]{}, true),
                new TestCase("different length lists are unequal", new Record[]{
                        new ARecord(new Name("foo.example.com."), DClass.IN, 60, InetAddress.getByName("1.2.3.4"))
                }, new Record[]{}, false),
                new TestCase("same records but different order lists are equal", new Record[]{
                        new ARecord(new Name("foo.example.com."), DClass.IN, 60, InetAddress.getByName("1.2.3.4")),
                        new ARecord(new Name("bar.example.com."), DClass.IN, 60, InetAddress.getByName("1.2.3.5")),
                }, new Record[]{
                        new ARecord(new Name("bar.example.com."), DClass.IN, 60, InetAddress.getByName("1.2.3.5")),
                        new ARecord(new Name("foo.example.com."), DClass.IN, 60, InetAddress.getByName("1.2.3.4")),
                }, true),
                new TestCase("same non-empty lists are equal", new Record[]{
                        new ARecord(new Name("foo.example.com."), DClass.IN, 60, InetAddress.getByName("1.2.3.4")),
                        new ARecord(new Name("bar.example.com."), DClass.IN, 60, InetAddress.getByName("1.2.3.5")),
                }, new Record[]{
                        new ARecord(new Name("foo.example.com."), DClass.IN, 60, InetAddress.getByName("1.2.3.4")),
                        new ARecord(new Name("bar.example.com."), DClass.IN, 60, InetAddress.getByName("1.2.3.5")),
                }, true),
                new TestCase("lists that only differ in the SOA serial number are equal", new Record[]{
                        new ARecord(new Name("foo.example.com."), DClass.IN, 60, InetAddress.getByName("1.2.3.4")),
                        new ARecord(new Name("bar.example.com."), DClass.IN, 60, InetAddress.getByName("1.2.3.5")),
                        new SOARecord(new Name("example.com."), DClass.IN, 60, new Name("example.com."), new Name("example.com."), 1, 60, 1, 1, 1),
                }, new Record[]{
                        new ARecord(new Name("foo.example.com."), DClass.IN, 60, InetAddress.getByName("1.2.3.4")),
                        new ARecord(new Name("bar.example.com."), DClass.IN, 60, InetAddress.getByName("1.2.3.5")),
                        new SOARecord(new Name("example.com."), DClass.IN, 60, new Name("example.com."), new Name("example.com."), 2, 60, 1, 1, 1),
                }, true),
                new TestCase("lists that differ in the SOA (other than the serial number) are not equal", new Record[]{
                        new ARecord(new Name("foo.example.com."), DClass.IN, 60, InetAddress.getByName("1.2.3.4")),
                        new ARecord(new Name("bar.example.com."), DClass.IN, 60, InetAddress.getByName("1.2.3.5")),
                        new SOARecord(new Name("example.com."), DClass.IN, 60, new Name("example.com."), new Name("example.com."), 1, 60, 1, 1, 1),
                }, new Record[]{
                        new ARecord(new Name("foo.example.com."), DClass.IN, 60, InetAddress.getByName("1.2.3.4")),
                        new ARecord(new Name("bar.example.com."), DClass.IN, 60, InetAddress.getByName("1.2.3.5")),
                        new SOARecord(new Name("example.com."), DClass.IN, 61, new Name("example.com."), new Name("example.com."), 2, 60, 1, 1, 1),
                }, false),
                new TestCase("lists that only differ in NSEC or RRSIG records are equal", new Record[]{
                        new ARecord(new Name("foo.example.com."), DClass.IN, 60, InetAddress.getByName("1.2.3.4")),
                        new ARecord(new Name("bar.example.com."), DClass.IN, 60, InetAddress.getByName("1.2.3.5")),
                        new SOARecord(new Name("example.com."), DClass.IN, 60, new Name("example.com."), new Name("example.com."), 1, 60, 1, 1, 1),
                        new NSECRecord(new Name("foo.example.com."), DClass.IN, 60, new Name("example.com."), new int[]{1}),
                        new RRSIGRecord(new Name("foo.example.com."), DClass.IN, 60, 1, 1, 60, new Date(), new Date(), 1, new Name("example.com."), new byte[]{1})
                }, new Record[]{
                        new ARecord(new Name("foo.example.com."), DClass.IN, 60, InetAddress.getByName("1.2.3.4")),
                        new ARecord(new Name("bar.example.com."), DClass.IN, 60, InetAddress.getByName("1.2.3.5")),
                        new SOARecord(new Name("example.com."), DClass.IN, 60, new Name("example.com."), new Name("example.com."), 2, 60, 1, 1, 1),
                }, true),
                new TestCase("lists that only differ in NSEC or RRSIG records are equal", new Record[]{
                        new ARecord(new Name("foo.example.com."), DClass.IN, 60, InetAddress.getByName("1.2.3.4")),
                        new ARecord(new Name("bar.example.com."), DClass.IN, 60, InetAddress.getByName("1.2.3.5")),
                        new SOARecord(new Name("example.com."), DClass.IN, 60, new Name("example.com."), new Name("example.com."), 1, 60, 1, 1, 1),
                }, new Record[]{
                        new ARecord(new Name("foo.example.com."), DClass.IN, 60, InetAddress.getByName("1.2.3.4")),
                        new ARecord(new Name("bar.example.com."), DClass.IN, 60, InetAddress.getByName("1.2.3.5")),
                        new SOARecord(new Name("example.com."), DClass.IN, 60, new Name("example.com."), new Name("example.com."), 2, 60, 1, 1, 1),
                        new NSECRecord(new Name("foo.example.com."), DClass.IN, 60, new Name("example.com."), new int[]{1}),
                        new RRSIGRecord(new Name("foo.example.com."), DClass.IN, 60, 1, 1, 60, new Date(), new Date(), 1, new Name("example.com."), new byte[]{1})
                }, true),
        };

        for (TestCase t : testCases) {
            List<Record> input1 = Arrays.asList(t.r1);
            List<Record> input2 = Arrays.asList(t.r2);
            List<Record> copy1 = Arrays.asList(t.r1);
            List<Record> copy2 = Arrays.asList(t.r2);
            boolean actual = ZoneManager.zonesAreEqual(input1, input2);
            assertThat(t.reason, actual, equalTo(t.expected));

            // assert that the input lists were not modified
            assertThat("zonesAreEqual input lists should not be modified", input1, equalTo(copy1));
            assertThat("zonesAreEqual input lists should not be modified", input2, equalTo(copy2));
        }
    }
}
