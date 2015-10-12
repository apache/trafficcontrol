package com.comcast.cdn.traffic_control.traffic_router.core.dns;

import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker.Track.ResultType;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker.Track.ResultDetails;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;
import org.xbill.DNS.*;

import java.net.Inet4Address;
import java.net.InetAddress;
import java.util.Random;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;
import static org.powermock.api.mockito.PowerMockito.mockStatic;
import static org.powermock.api.mockito.PowerMockito.spy;
import static org.powermock.api.mockito.PowerMockito.whenNew;

@RunWith(PowerMockRunner.class)
@PrepareForTest({Random.class, Header.class, DNSAccessEventBuilder.class})
public class DNSAccessEventBuilderTest {

    private InetAddress client;

    @Before
    public void before() throws Exception {
        mockStatic(System.class);

        Random random = mock(Random.class);
        when(random.nextInt(0xffff)).thenReturn(65535);
        whenNew(Random.class).withNoArguments().thenReturn(random);

        client = mock(InetAddress.class);
        when(client.getHostAddress()).thenReturn("192.168.10.11");
    }

    @Test
    public void itCreatesRequestErrorData() throws Exception {
        when(System.currentTimeMillis()).thenReturn(144140678789L);

        DNSAccessRecord dnsAccessRecord = new DNSAccessRecord.Builder(144140678000L, client).build();

        String dnsAccessEvent = DNSAccessEventBuilder.create(dnsAccessRecord, new WireParseException("invalid record length"));
        assertThat(dnsAccessEvent, equalTo("144140678.000 qtype=DNS chi=192.168.10.11 ttms=789 xn=- fqdn=- type=- class=- ttl=- rcode=-" +
                " rtype=- rdetails=- rerr=\"Bad Request:WireParseException:invalid record length\" ans=\"-\""));
    }

    @Test
    public void itAddsResponseData() throws Exception {
        final Name name = Name.fromString("www.example.com.");

        when(System.currentTimeMillis()).thenReturn(144140678789L).thenReturn(144140678000L);

        final Record question = Record.newRecord(name, Type.A, DClass.IN, 12345L);

        final Message response = spy(Message.newQuery(question));
        response.getHeader().setRcode(Rcode.NOERROR);

        final Record record1 = mock(Record.class);
        when(record1.rdataToString()).thenReturn("foo");
        final Record record2 = mock(Record.class);
        when(record2.rdataToString()).thenReturn("bar");
        final Record record3 = mock(Record.class);
        when(record3.rdataToString()).thenReturn("baz");

        Record[] records = new Record[] {record1, record2, record3};
        when(response.getSectionArray(Section.ANSWER)).thenReturn(records);

        InetAddress answerAddress = Inet4Address.getByName("192.168.1.23");

        ARecord addressRecord = new ARecord(name, DClass.IN, 54321L, answerAddress);
        response.addRecord(addressRecord, Section.ANSWER);

        DNSAccessRecord dnsAccessRecord = new DNSAccessRecord.Builder(144140678000L, client).dnsMessage(response).build();
        String dnsAccessEvent = DNSAccessEventBuilder.create(dnsAccessRecord);

        assertThat(dnsAccessEvent, equalTo("144140678.000 qtype=DNS chi=192.168.10.11 ttms=789" +
                " xn=65535 fqdn=www.example.com. type=A class=IN ttl=12345" +
                " rcode=NOERROR rtype=- rdetails=- rerr=\"-\" ans=\"foo bar baz\""));


        dnsAccessEvent = DNSAccessEventBuilder.create(dnsAccessRecord);

        assertThat(dnsAccessEvent, equalTo("144140678.000 qtype=DNS chi=192.168.10.11 ttms=0" +
                " xn=65535 fqdn=www.example.com. type=A class=IN ttl=12345" +
                " rcode=NOERROR rtype=- rdetails=- rerr=\"-\" ans=\"foo bar baz\""));
    }

    @Test
    public void itCreatesServerErrorData() throws Exception {
        Message query = Message.newQuery(Record.newRecord(Name.fromString("www.example.com."), Type.A, DClass.IN, 12345L));
        when(System.currentTimeMillis()).thenReturn(144140678789L);

        DNSAccessRecord dnsAccessRecord = new DNSAccessRecord.Builder(144140678000L, client).dnsMessage(query).build();
        String dnsAccessEvent = DNSAccessEventBuilder.create(dnsAccessRecord, new RuntimeException("boom it failed"));
        assertThat(dnsAccessEvent, equalTo("144140678.000 qtype=DNS chi=192.168.10.11 ttms=789" +
                " xn=65535 fqdn=www.example.com. type=A class=IN ttl=12345" +
                " rcode=SERVFAIL rtype=- rdetails=- rerr=\"Server Error:RuntimeException:boom it failed\" ans=\"-\""));
    }

    @Test
    public void itAddsResultTypeData() throws Exception {
        final Name name = Name.fromString("www.example.com.");

        when(System.currentTimeMillis()).thenReturn(144140678789L).thenReturn(144140678000L);

        final Record question = Record.newRecord(name, Type.A, DClass.IN, 12345L);
        final Message response = spy(Message.newQuery(question));
        response.getHeader().setRcode(Rcode.NOERROR);

        final Record record1 = mock(Record.class);
        when(record1.rdataToString()).thenReturn("foo");
        final Record record2 = mock(Record.class);
        when(record2.rdataToString()).thenReturn("bar");
        final Record record3 = mock(Record.class);
        when(record3.rdataToString()).thenReturn("baz");

        Record[] records = new Record[] {record1, record2, record3};
        when(response.getSectionArray(Section.ANSWER)).thenReturn(records);

        InetAddress answerAddress = Inet4Address.getByName("192.168.1.23");

        ARecord addressRecord = new ARecord(name, DClass.IN, 54321L, answerAddress);
        response.addRecord(addressRecord, Section.ANSWER);

        ResultType resultType = ResultType.CZ;
        final DNSAccessRecord.Builder builder = new DNSAccessRecord.Builder(144140678000L, client).dnsMessage(response).resultType(resultType);
        DNSAccessRecord dnsAccessRecord = builder.build();
        String dnsAccessEvent = DNSAccessEventBuilder.create(dnsAccessRecord);

        assertThat(dnsAccessEvent, equalTo("144140678.000 qtype=DNS chi=192.168.10.11 ttms=789" +
                " xn=65535 fqdn=www.example.com. type=A class=IN ttl=12345" +
                " rcode=NOERROR rtype=CZ rdetails=- rerr=\"-\" ans=\"foo bar baz\""));

        dnsAccessRecord = builder.resultType(ResultType.GEO).build();
        dnsAccessEvent = DNSAccessEventBuilder.create(dnsAccessRecord);

        assertThat(dnsAccessEvent, equalTo("144140678.000 qtype=DNS chi=192.168.10.11 ttms=0" +
                " xn=65535 fqdn=www.example.com. type=A class=IN ttl=12345" +
                " rcode=NOERROR rtype=GEO rdetails=- rerr=\"-\" ans=\"foo bar baz\""));

        dnsAccessRecord = builder.resultType(ResultType.MISS).resultDetails(ResultDetails.DS_NOT_FOUND).build();
        dnsAccessEvent = DNSAccessEventBuilder.create(dnsAccessRecord);

        assertThat(dnsAccessEvent, equalTo("144140678.000 qtype=DNS chi=192.168.10.11 ttms=0" +
                " xn=65535 fqdn=www.example.com. type=A class=IN ttl=12345" +
                " rcode=NOERROR rtype=MISS rdetails=DS_NOT_FOUND rerr=\"-\" ans=\"foo bar baz\""));
    }
}