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

import java.net.InetAddress;
import java.security.PrivateKey;
import java.util.ArrayList;
import java.util.Collections;
import java.util.Date;
import java.util.List;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ConcurrentMap;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.Matchers.notNullValue;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.argThat;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.ArgumentMatcher;
import org.powermock.api.mockito.PowerMockito;
import org.powermock.core.classloader.annotations.PowerMockIgnore;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;
import org.xbill.DNS.ARecord;
import org.xbill.DNS.DClass;
import org.xbill.DNS.DNSKEYRecord;
import org.xbill.DNS.Name;
import org.xbill.DNS.RRSIGRecord;
import org.xbill.DNS.RRset;
import org.xbill.DNS.Record;
import org.xbill.DNS.Type;

@RunWith(PowerMockRunner.class)
@PrepareForTest(ZoneSignerImpl.class)
@PowerMockIgnore("javax.management.*")
public class ZoneSignerImplTest {

    static class IsRRsetTypeA implements ArgumentMatcher<RRset> {
        @Override
        public boolean matches(RRset rRset) {
            return rRset.getType() == Type.A;
        }
    }

    static class IsRRsetTypeNSEC implements ArgumentMatcher<RRset> {
        @Override
        public boolean matches(RRset rRset) {
            return rRset.getType() == Type.NSEC;
        }
    }

    @Test
    public void signZoneWithRRSIGCacheTest() throws Exception {
        ZoneSignerImpl zoneSigner = PowerMockito.spy(new ZoneSignerImpl());
        List<Record> records = new ArrayList<>();
        Record ARecord1 = new ARecord(new Name("foo.example.com."), DClass.IN, 60, InetAddress.getByName("1.2.3.4"));
        Record ARecord2 = new ARecord(new Name("foo.example.com."), DClass.IN, 60, InetAddress.getByName("1.2.3.5"));
        Record ARecord3 = new ARecord(new Name("foo.example.com."), DClass.IN, 60, InetAddress.getByName("1.2.3.6"));
        Record ARecord4 = new ARecord(new Name("foo.example.com."), DClass.IN, 60, InetAddress.getByName("1.2.3.7"));
        DnsSecKeyPair zskPair = mock(DnsSecKeyPairImpl.class);
        DNSKEYRecord zskDnskey = mock(DNSKEYRecord.class);
        when(zskPair.getDNSKEYRecord()).thenReturn(zskDnskey);
        PrivateKey zskKey = mock(PrivateKey.class);
        when(zskPair.getPrivate()).thenReturn(zskKey);
        when(zskKey.getEncoded()).thenReturn(new byte[] {1});
        when(zskDnskey.getAlgorithm()).thenReturn(1);
        List<DnsSecKeyPair> kskPairs = new ArrayList<>();
        List<DnsSecKeyPair> zskPairs = Collections.singletonList(zskPair);

        Date inception = new Date();
        Date expire = Date.from(inception.toInstant().plusSeconds(100000));
        RRSIGRecord aRRSigRecord = new RRSIGRecord(new Name("foo.example.com."), DClass.IN, 60, Type.A, 1, 60, inception, expire, 1, new Name("example.com."), new byte[]{1});
        RRSIGRecord nsecRRSigRecord = new RRSIGRecord(new Name("foo.example.com."), DClass.IN, 60, Type.NSEC, 1, 60, inception, expire, 1, new Name("example.com."), new byte[]{2});
        PowerMockito.doReturn(aRRSigRecord).when(zoneSigner, "sign", argThat(new IsRRsetTypeA()), any(DNSKEYRecord.class), any(PrivateKey.class), eq(inception), eq(expire));
        PowerMockito.doReturn(nsecRRSigRecord).when(zoneSigner, "sign", argThat(new IsRRsetTypeNSEC()), any(DNSKEYRecord.class), any(PrivateKey.class), eq(inception), eq(expire));

        Date newInception = Date.from(inception.toInstant().plusSeconds(100));
        Date newExpire = Date.from(newInception.toInstant().plusSeconds(100000));
        RRSIGRecord newARRSigRecord = new RRSIGRecord(new Name("foo.example.com."), DClass.IN, 60, Type.A, 1, 60, newInception, newExpire, 1, new Name("example.com."), new byte[]{3});
        RRSIGRecord newNSECRRSigRecord = new RRSIGRecord(new Name("foo.example.com."), DClass.IN, 60, Type.NSEC, 1, 60, newInception, newExpire, 1, new Name("example.com."), new byte[]{4});
        PowerMockito.doReturn(newARRSigRecord).when(zoneSigner, "sign", argThat(new IsRRsetTypeA()), any(DNSKEYRecord.class), any(PrivateKey.class), eq(newInception), eq(newExpire));
        PowerMockito.doReturn(newNSECRRSigRecord).when(zoneSigner, "sign", argThat(new IsRRsetTypeNSEC()), any(DNSKEYRecord.class), any(PrivateKey.class), eq(newInception), eq(newExpire));

        Date expiresSoonInception = Date.from(inception.toInstant().minusSeconds(100));
        Date expiresSoonExpire = Date.from(inception.toInstant().plusSeconds(50));
        RRSIGRecord expiresSoonARRSigRecord = new RRSIGRecord(new Name("foo.example.com."), DClass.IN, 60, Type.A, 1, 60, expiresSoonInception, expiresSoonExpire, 1, new Name("example.com."), new byte[]{5});
        RRSIGRecord expiresSoonNSECRRSigRecord = new RRSIGRecord(new Name("foo.example.com."), DClass.IN, 60, Type.NSEC, 1, 60, expiresSoonInception, expiresSoonExpire, 1, new Name("example.com."), new byte[]{6});
        PowerMockito.doReturn(expiresSoonARRSigRecord).when(zoneSigner, "sign", argThat(new IsRRsetTypeA()), any(DNSKEYRecord.class), any(PrivateKey.class), eq(expiresSoonInception), eq(expiresSoonExpire));
        PowerMockito.doReturn(expiresSoonNSECRRSigRecord).when(zoneSigner, "sign", argThat(new IsRRsetTypeNSEC()), any(DNSKEYRecord.class), any(PrivateKey.class), eq(expiresSoonInception), eq(expiresSoonExpire));

        ConcurrentMap<RRSIGCacheKey, ConcurrentMap<RRsetKey, RRSIGRecord>> RRSIGCache = new ConcurrentHashMap<>();

        records.add(ARecord1);
        records.add(ARecord2);
        List<Record> signedRecords = zoneSigner.signZone(records, kskPairs, zskPairs, inception, expire, RRSIGCache);
        RRSIGRecord ret = (RRSIGRecord) signedRecords.stream().filter(r -> r instanceof RRSIGRecord && ((RRSIGRecord) r).getTypeCovered() == Type.A).findFirst().orElse(null);
        assertThat(ret, notNullValue());
        assertThat(ret, equalTo(aRRSigRecord));

        // re-signing the same RRset with new timestamps should reuse the cached RRSIG record
        records.clear();
        records.add(ARecord1);
        records.add(ARecord2);
        signedRecords = zoneSigner.signZone(records, kskPairs, zskPairs, newInception, newExpire, RRSIGCache);
        ret = (RRSIGRecord) signedRecords.stream().filter(r -> r instanceof RRSIGRecord && ((RRSIGRecord) r).getTypeCovered() == Type.A).findFirst().orElse(null);
        assertThat(ret, notNullValue());
        assertThat(ret, equalTo(aRRSigRecord));

        // changed RRset should be re-signed
        records.clear();
        records.add(ARecord1);
        records.add(ARecord2);
        records.add(ARecord3);
        records.add(ARecord4);
        signedRecords = zoneSigner.signZone(records, kskPairs, zskPairs, newInception, newExpire, RRSIGCache);
        ret = (RRSIGRecord) signedRecords.stream().filter(r -> r instanceof RRSIGRecord && ((RRSIGRecord) r).getTypeCovered() == Type.A).findFirst().orElse(null);
        assertThat(ret, notNullValue());
        assertThat(ret, equalTo(newARRSigRecord));

        // re-signing 1st RRset again should reuse the cached RRSIG record
        records.clear();
        records.add(ARecord1);
        records.add(ARecord2);
        signedRecords = zoneSigner.signZone(records, kskPairs, zskPairs, newInception, newExpire, RRSIGCache);
        ret = (RRSIGRecord) signedRecords.stream().filter(r -> r instanceof RRSIGRecord && ((RRSIGRecord) r).getTypeCovered() == Type.A).findFirst().orElse(null);
        assertThat(ret, notNullValue());
        assertThat(ret, equalTo(aRRSigRecord));

        // re-signing RRset that has a cached RRSIG record that is close to expiring should be re-signed
        records.clear();
        records.add(ARecord3);
        records.add(ARecord4);
        signedRecords = zoneSigner.signZone(records, kskPairs, zskPairs, expiresSoonInception, expiresSoonExpire, RRSIGCache);
        ret = (RRSIGRecord) signedRecords.stream().filter(r -> r instanceof RRSIGRecord && ((RRSIGRecord) r).getTypeCovered() == Type.A).findFirst().orElse(null);
        assertThat(ret, notNullValue());
        assertThat(ret, equalTo(expiresSoonARRSigRecord));
        records.clear();
        records.add(ARecord3);
        records.add(ARecord4);
        signedRecords = zoneSigner.signZone(records, kskPairs, zskPairs, newInception, newExpire, RRSIGCache);
        ret = (RRSIGRecord) signedRecords.stream().filter(r -> r instanceof RRSIGRecord && ((RRSIGRecord) r).getTypeCovered() == Type.A).findFirst().orElse(null);
        assertThat(ret, notNullValue());
        assertThat(ret, equalTo(newARRSigRecord));
    }
}
