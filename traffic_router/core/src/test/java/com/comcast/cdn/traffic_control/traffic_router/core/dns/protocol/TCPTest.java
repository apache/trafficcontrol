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

import static org.junit.Assert.assertArrayEquals;
import static org.junit.Assert.assertEquals;

import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.io.DataOutputStream;
import java.io.IOException;
import java.net.InetAddress;
import java.net.ServerSocket;
import java.net.Socket;
import java.util.concurrent.ExecutorService;

import org.jmock.Expectations;
import org.jmock.Mockery;
import org.jmock.integration.junit4.JMock;
import org.jmock.integration.junit4.JUnit4Mockery;
import org.jmock.lib.legacy.ClassImposteriser;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.xbill.DNS.DClass;
import org.xbill.DNS.Message;
import org.xbill.DNS.Name;
import org.xbill.DNS.Rcode;
import org.xbill.DNS.Record;
import org.xbill.DNS.Section;
import org.xbill.DNS.Type;

import com.comcast.cdn.traffic_control.traffic_router.core.dns.NameServer;
import com.comcast.cdn.traffic_control.traffic_router.core.dns.protocol.TCP;
import com.comcast.cdn.traffic_control.traffic_router.core.dns.protocol.TCP.TCPSocketHandler;

@RunWith(JMock.class)
public class TCPTest {
    private final Mockery context = new JUnit4Mockery() {
        {
            setImposteriser(ClassImposteriser.INSTANCE);
        }
    };

    private ServerSocket serverSocket;
    private Socket socket;
    private ExecutorService executorService;
    private NameServer nameServer;

    private TCP tcp;

    @Before
    public void setUp() throws Exception {
        serverSocket = context.mock(ServerSocket.class);
        socket = context.mock(Socket.class);
        executorService = context.mock(ExecutorService.class);
        nameServer = context.mock(NameServer.class);
        tcp = new TCP();
        tcp.setServerSocket(serverSocket);
        tcp.setExecutorService(executorService);
        tcp.setNameServer(nameServer);
    }

    @Test
    public void testGetMaxResponseLength() {
        assertEquals(Integer.MAX_VALUE, tcp.getMaxResponseLength(null));
    }

    @Test
    public void testSubmit() {
        final Runnable r = context.mock(Runnable.class);
        context.checking(new Expectations() {
            {
                one(executorService).submit(r);
            }
        });
        tcp.submit(r);
    }

    @Test
    public void testTCPSocketHandler() throws Exception {
        final InetAddress client = InetAddress.getLocalHost();
        final TCPSocketHandler handler = tcp.new TCPSocketHandler(socket);

        final Name name = Name.fromString("www.foo.bar.");
        final Record question = Record.newRecord(name, Type.A, DClass.IN);
        final Message request = Message.newQuery(question);
        final byte[] wireRequest = request.toWire();

        final ByteArrayOutputStream baos = new ByteArrayOutputStream();
        final DataOutputStream dos = new DataOutputStream(baos);
        dos.writeShort(wireRequest.length);
        dos.write(wireRequest);

        final ByteArrayInputStream in = new ByteArrayInputStream(baos.toByteArray());
        final ByteArrayOutputStream out = new ByteArrayOutputStream();

        context.checking(new Expectations() {
            {
                one(socket).getInetAddress();
                will(returnValue(client));
                one(socket).getInputStream();
                will(returnValue(in));
                one(socket).getOutputStream();
                will(returnValue(out));
                one(socket).close();

                one(nameServer).query(with(any(Message.class)), with(same(client)));
                will(returnValue(request));
            }
        });
        handler.run();
        final byte[] expected = baos.toByteArray();
        final byte[] actual = out.toByteArray();
        assertArrayEquals(expected, actual);
    }

    @Test
    public void testTCPSocketHandlerBadMessage() throws Exception {
        final InetAddress client = InetAddress.getLocalHost();
        final TCPSocketHandler handler = tcp.new TCPSocketHandler(socket);

        final byte[] wireRequest = new byte[0];

        final ByteArrayOutputStream baos = new ByteArrayOutputStream();
        final DataOutputStream dos = new DataOutputStream(baos);
        dos.writeShort(wireRequest.length);
        dos.write(wireRequest);

        final ByteArrayInputStream in = new ByteArrayInputStream(baos.toByteArray());
        final ByteArrayOutputStream out = new ByteArrayOutputStream();

        context.checking(new Expectations() {
            {
                one(socket).getInetAddress();
                will(returnValue(client));
                one(socket).getInputStream();
                will(returnValue(in));
                one(socket).getOutputStream();
                will(returnValue(out));
                one(socket).close();
            }
        });
        handler.run();
        final byte[] expected = new byte[0];
        final byte[] actual = out.toByteArray();
        assertArrayEquals(expected, actual);
    }

    @Test
    public void testTCPSocketHandlerQueryFail() throws Exception {
        final InetAddress client = InetAddress.getLocalHost();
        final TCPSocketHandler handler = tcp.new TCPSocketHandler(socket);

        final Name name = Name.fromString("www.foo.bar.");
        final Record question = Record.newRecord(name, Type.A, DClass.IN);
        final Message request = Message.newQuery(question);
        final byte[] wireRequest = request.toWire();

        final ByteArrayOutputStream baos = new ByteArrayOutputStream();
        final DataOutputStream dos = new DataOutputStream(baos);
        dos.writeShort(wireRequest.length);
        dos.write(wireRequest);

        final ByteArrayInputStream in = new ByteArrayInputStream(baos.toByteArray());
        final ByteArrayOutputStream out = new ByteArrayOutputStream();

        final Message response = new Message();
        response.setHeader(request.getHeader());
        for (int i = 0; i < 4; i++) {
            response.removeAllRecords(i);
        }
        response.addRecord(question, Section.QUESTION);
        response.getHeader().setRcode(Rcode.SERVFAIL);
        final byte[] serverFail = response.toWire();

        final ByteArrayOutputStream baos2 = new ByteArrayOutputStream();
        final DataOutputStream dos2 = new DataOutputStream(baos2);
        dos2.writeShort(serverFail.length);
        dos2.write(serverFail);

        context.checking(new Expectations() {
            {
                one(socket).getInetAddress();
                will(returnValue(client));
                one(socket).getInputStream();
                will(returnValue(in));
                one(socket).getOutputStream();
                will(returnValue(out));
                one(socket).close();
                will(throwException(new IOException()));

                one(nameServer).query(with(any(Message.class)), with(same(client)));
                will(throwException(new Exception()));
            }
        });
        handler.run();
        final byte[] expected = baos2.toByteArray();
        final byte[] actual = out.toByteArray();
        assertArrayEquals(expected, actual);
    }

}
