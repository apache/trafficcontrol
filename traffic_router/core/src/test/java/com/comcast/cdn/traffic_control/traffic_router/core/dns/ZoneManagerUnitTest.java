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

package com.comcast.cdn.traffic_control.traffic_router.core.dns;

import com.comcast.cdn.traffic_control.traffic_router.core.edge.CacheRegister;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker.Track.ResultType;
import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouterManager;
import com.comcast.cdn.traffic_control.traffic_router.core.util.TrafficOpsUtils;
import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouter;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
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

import java.net.InetAddress;
import java.util.Arrays;
import java.util.Date;
import java.util.List;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.mockito.Matchers.any;
import static org.mockito.Matchers.anyInt;
import static org.mockito.Mockito.*;
import static org.powermock.api.mockito.PowerMockito.doCallRealMethod;
import static org.powermock.api.mockito.PowerMockito.whenNew;

@RunWith(PowerMockRunner.class)
@PrepareForTest({ZoneManager.class, SignatureManager.class})
@PowerMockIgnore("javax.management.*")
public class ZoneManagerUnitTest {
    ZoneManager zoneManager;

    @Before
    public void before() throws Exception {
        TrafficRouter trafficRouter = mock(TrafficRouter.class);
        CacheRegister cacheRegister = mock(CacheRegister.class);
        when(trafficRouter.getCacheRegister()).thenReturn(cacheRegister);

        PowerMockito.spy(ZoneManager.class);
        PowerMockito.doNothing().when(ZoneManager.class, "initTopLevelDomain", cacheRegister);
        PowerMockito.doNothing().when(ZoneManager.class, "initZoneCache", trafficRouter);

        SignatureManager signatureManager = PowerMockito.mock(SignatureManager.class);
        whenNew(SignatureManager.class).withArguments(any(ZoneManager.class), any(CacheRegister.class), any(TrafficOpsUtils.class), any(TrafficRouterManager.class)).thenReturn(signatureManager);

        zoneManager = spy(new ZoneManager(trafficRouter, new StatTracker(), null, mock(TrafficRouterManager.class)));
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
