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

package org.apache.traffic_control.traffic_router.core.dns.protocol;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.junit.Assert.assertEquals;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.doAnswer;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.LinkedBlockingQueue;
import java.util.concurrent.ThreadPoolExecutor;
import java.util.concurrent.atomic.AtomicInteger;

import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.core.classloader.annotations.PowerMockIgnore;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;
import org.xbill.DNS.DClass;
import org.xbill.DNS.Flags;
import org.xbill.DNS.Message;
import org.xbill.DNS.Name;
import org.xbill.DNS.OPTRecord;
import org.xbill.DNS.Rcode;
import org.xbill.DNS.Record;
import org.xbill.DNS.Section;
import org.xbill.DNS.Type;

import org.apache.traffic_control.traffic_router.core.dns.NameServer;
import org.apache.traffic_control.traffic_router.core.dns.protocol.UDP.UDPPacketHandler;
import org.apache.traffic_control.traffic_router.core.dns.DNSAccessRecord;

@RunWith(PowerMockRunner.class)
@PrepareForTest({AbstractProtocol.class, Message.class})
@PowerMockIgnore("javax.management.*")
public class UDPTest {

    private DatagramSocket datagramSocket;
    private ThreadPoolExecutor executorService;
    private ExecutorService cancelService;
    private LinkedBlockingQueue queue;
    private NameServer nameServer;

    private UDP udp;

    @Before
    public void setUp() throws Exception {
        datagramSocket = mock(DatagramSocket.class);
        executorService = mock(ThreadPoolExecutor.class);
        cancelService = mock(ExecutorService.class);
        queue = mock(LinkedBlockingQueue.class);
        nameServer = mock(NameServer.class);
        udp = new UDP();
        udp.setDatagramSocket(datagramSocket);
        udp.setExecutorService(executorService);
        udp.setCancelService(cancelService);
        udp.setNameServer(nameServer);

        when(executorService.getQueue()).thenReturn(queue);
        when(queue.size()).thenReturn(0);
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
    public void testSubmit() throws Exception {
        final SocketHandler r = mock(SocketHandler.class);
        udp.submit(r);
        verify(executorService).submit(r);
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

        when(nameServer.query(any(Message.class), eq(client), any(DNSAccessRecord.Builder.class))).thenReturn(response);

        final AtomicInteger count = new AtomicInteger(0);
        doAnswer(invocation -> {
            DatagramPacket datagramPacket = (DatagramPacket) invocation.getArguments()[0];
            assertThat(datagramPacket.getData(), equalTo(wireResponse));
            count.incrementAndGet();
            return null;
        }).when(datagramSocket).send(any(DatagramPacket.class));

        final UDPPacketHandler handler = udp.new UDPPacketHandler(packet);
        handler.run();
        assertThat(count.get(), equalTo(1));
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

        final AtomicInteger count = new AtomicInteger(0);

        when(nameServer.query(any(Message.class), eq(client), any(DNSAccessRecord.Builder.class))).thenThrow(new RuntimeException("Boom! UDP Query"));

        doAnswer(invocation -> {
            DatagramPacket datagramPacket = (DatagramPacket) invocation.getArguments()[0];
            assertThat(datagramPacket.getData(), equalTo(wireResponse));
            count.incrementAndGet();
            return null;
        }).when(datagramSocket).send(any(DatagramPacket.class));

        final UDPPacketHandler handler = udp.new UDPPacketHandler(packet);
        handler.run();
        assertThat(count.get(), equalTo(1));
    }
}
