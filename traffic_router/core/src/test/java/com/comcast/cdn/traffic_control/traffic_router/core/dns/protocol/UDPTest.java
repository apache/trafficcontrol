/*
 * Copyright 2015 Comcast Cable Communications Management, LLC
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

package com.comcast.cdn.traffic_control.traffic_router.core.dns.protocol;

import static org.junit.Assert.assertEquals;

import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;
import java.util.Arrays;
import java.util.concurrent.ExecutorService;

import org.hamcrest.Description;
import org.hamcrest.Matcher;
import org.hamcrest.TypeSafeMatcher;
import org.jmock.Expectations;
import org.jmock.Mockery;
import org.jmock.integration.junit4.JMock;
import org.jmock.integration.junit4.JUnit4Mockery;
import org.jmock.lib.legacy.ClassImposteriser;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.xbill.DNS.DClass;
import org.xbill.DNS.Flags;
import org.xbill.DNS.Message;
import org.xbill.DNS.Name;
import org.xbill.DNS.OPTRecord;
import org.xbill.DNS.Rcode;
import org.xbill.DNS.Record;
import org.xbill.DNS.Section;
import org.xbill.DNS.Type;

import com.comcast.cdn.traffic_control.traffic_router.core.dns.NameServer;
import com.comcast.cdn.traffic_control.traffic_router.core.dns.protocol.UDP;
import com.comcast.cdn.traffic_control.traffic_router.core.dns.protocol.UDP.UDPPacketHandler;

@RunWith(JMock.class)
public class UDPTest {

    private final Mockery context = new JUnit4Mockery() {
        {
            setImposteriser(ClassImposteriser.INSTANCE);
        }
    };

    private DatagramSocket datagramSocket;
    private ExecutorService executorService;
    private NameServer nameServer;

    private UDP udp;

    @Before
    public void setUp() throws Exception {
        datagramSocket = context.mock(DatagramSocket.class);
        executorService = context.mock(ExecutorService.class);
        nameServer = context.mock(NameServer.class);
        udp = new UDP();
        udp.setDatagramSocket(datagramSocket);
        udp.setExecutorService(executorService);
        udp.setNameServer(nameServer);

    }

    @Test
    public void testGetMaxResponseLengthNoOPTQuery() throws Exception {
        final Name name = Name.fromString("www.foo.com.");
        final Record question = Record.newRecord(name, Type.A, DClass.IN);
        final Message request = Message.newQuery(question);
        assertEquals(512, udp.getMaxResponseLength(request));
    }

    @Test
    public void testGetMaxResponseLengthNullQuery() {
        assertEquals(512, udp.getMaxResponseLength(null));
    }

    @Test
    public void testGetMaxResponseLengthWithOPTQuery() throws Exception {
        final int size = 1280;
        final Name name = Name.fromString("www.foo.com.");
        final Record question = Record.newRecord(name, Type.A, DClass.IN);
        final OPTRecord options = new OPTRecord(size, 0, 0);
        final Message request = Message.newQuery(question);
        request.addRecord(options, Section.ADDITIONAL);
        assertEquals(size, udp.getMaxResponseLength(request));
    }

    @Test
    public void testSubmit() {
        final Runnable r = context.mock(Runnable.class);
        context.checking(new Expectations() {
            {
                one(executorService).submit(r);
            }
        });
        udp.submit(r);
    }

    @Test
    public void testUDPPacketHandler() throws Exception {
        final InetAddress client = InetAddress.getLocalHost();
        final int port = 11111;

        final Name name = Name.fromString("www.foo.bar.");
        final Record question = Record.newRecord(name, Type.A, DClass.IN);
        final Message request = Message.newQuery(question);
        final byte[] wireRequest = request.toWire();

        final Record aRecord = Record.newRecord(name, Type.A, DClass.IN, 3600);
        final Message response = Message.newQuery(question);
        response.getHeader().setFlag(Flags.QR);
        response.addRecord(aRecord, Section.ANSWER);
        final byte[] wireResponse = response.toWire();

        final DatagramPacket packet = new DatagramPacket(wireRequest, wireRequest.length, client, port);

        context.checking(new Expectations() {
            {
                one(nameServer).query(with(any(Message.class)), with(same(client)));
                will(returnValue(response));

                one(datagramSocket).send(with(aDatagramPacketWithThePayload(wireResponse)));
            }
        });
        final UDPPacketHandler handler = udp.new UDPPacketHandler(packet);
        handler.run();
    }

    @Test
    public void testUDPPacketHandlerBadMessage() throws Exception {
        final InetAddress client = InetAddress.getLocalHost();
        final int port = 11111;

        final byte[] wireRequest = new byte[0];

        final DatagramPacket packet = new DatagramPacket(wireRequest, wireRequest.length, client, port);

        final UDPPacketHandler handler = udp.new UDPPacketHandler(packet);
        handler.run();
    }

    @Test
    public void testUDPPacketHandlerQueryFail() throws Exception {
        final InetAddress client = InetAddress.getLocalHost();
        final int port = 11111;

        final Name name = Name.fromString("www.foo.bar.");
        final Record question = Record.newRecord(name, Type.A, DClass.IN);
        final Message request = Message.newQuery(question);
        final byte[] wireRequest = request.toWire();

        final Message response = new Message();
        response.setHeader(request.getHeader());
        for (int i = 0; i < 4; i++) {
            response.removeAllRecords(i);
        }
        response.addRecord(question, Section.QUESTION);
        response.getHeader().setRcode(Rcode.SERVFAIL);
        final byte[] wireResponse = response.toWire();

        final DatagramPacket packet = new DatagramPacket(wireRequest, wireRequest.length, client, port);

        context.checking(new Expectations() {
            {
                one(nameServer).query(with(any(Message.class)), with(same(client)));
                will(throwException(new Exception()));

                one(datagramSocket).send(with(aDatagramPacketWithThePayload(wireResponse)));
            }
        });
        final UDPPacketHandler handler = udp.new UDPPacketHandler(packet);
        handler.run();
    }

    private static Matcher<DatagramPacket> aDatagramPacketWithThePayload(final byte[] payload) {
        return new DatagramPacketPayloadMatcher(payload);
    }

    private static class DatagramPacketPayloadMatcher extends TypeSafeMatcher<DatagramPacket> {

        private final byte[] payload;

        private DatagramPacketPayloadMatcher(final byte[] payload) {
            this.payload = payload;
        }

        @Override
        public void describeTo(final Description description) {
            description.appendText("a DatagramPacket with the specified payload.");
        }

        @Override
        public boolean matchesSafely(final DatagramPacket item) {
            return Arrays.equals(payload, item.getData());
        }

    }
}
