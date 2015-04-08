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

package com.comcast.cdn.traffic_control.traffic_router.core.dns;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;

import java.net.InetAddress;
import java.util.Date;
import java.util.List;

import org.apache.log4j.Level;
import org.apache.log4j.Logger;
import org.apache.log4j.spi.LoggingEvent;
import org.jmock.Expectations;
import org.jmock.Mockery;
import org.jmock.integration.junit4.JMock;
import org.jmock.integration.junit4.JUnit4Mockery;
import org.jmock.lib.legacy.ClassImposteriser;
import org.junit.AfterClass;
import org.junit.Before;
import org.junit.BeforeClass;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.xbill.DNS.Header;
import org.xbill.DNS.Message;
import org.xbill.DNS.Name;
import org.xbill.DNS.Rcode;
import org.xbill.DNS.Record;
import org.xbill.DNS.Section;
import org.xbill.DNS.Type;

import com.comcast.cdn.traffic_control.traffic_router.core.dns.DNSAccessRecord;
import com.comcast.cdn.traffic_control.traffic_router.logger.NoLogAppender;

@RunWith(JMock.class)
public class DNSAccessRecordTest {

    private static final Logger LOGGER = Logger.getLogger("com.comcast.cdn.traffic_control.traffic_router.core.access");
    private static final NoLogAppender APPENDER = new NoLogAppender();

    private final Mockery context = new JUnit4Mockery() {
        {
            setImposteriser(ClassImposteriser.INSTANCE);
        }
    };

    private DNSAccessRecord record;
    private Message request;
    private Message response;

    @Before
    public void setUp() throws Exception {
        record = new DNSAccessRecord();
        request = context.mock(Message.class, "request");
        response = context.mock(Message.class, "response");
        APPENDER.clear();
    }

    @Test
    public void testLogNoRequestNoResponse() throws Exception {
        final Date date = new Date(0);
        final InetAddress client = InetAddress.getByName("127.0.0.1");

        final String expected = String.format("DNS [01/Jan/1970:00:00:00.000 +0000] %s %s %s %s \"%s\"",
                client.getHostAddress(), "-", "-", "-", "-");

        record.setRequestDate(date);
        record.setClient(client);
        record.setRequest(null);
        record.setResponse(null);

        context.checking(new Expectations() {
        });

        record.log();

        final List<LoggingEvent> events = APPENDER.getEvents();
        assertNotNull(events);
        assertEquals(1, events.size());
        final LoggingEvent event = events.get(0);
        assertEquals(Level.INFO, event.getLevel());
        assertEquals(expected, event.getMessage());
    }

    @Test
    public void testLogWithRequestNoResponse() throws Exception {
        final Date date = new Date(0);
        final InetAddress client = InetAddress.getByName("127.0.0.1");
        final String requestName = "www.foo.com.";
        final int requestType = Type.A;

        final Record question = context.mock(Record.class, "question");

        final String expected = String.format("DNS [01/Jan/1970:00:00:00.000 +0000] %s %s %s %s \"%s\"",
                client.getHostAddress(), Type.string(requestType), requestName, "-", "-");

        record.setRequestDate(date);
        record.setClient(client);
        record.setRequest(request);
        record.setResponse(null);

        context.checking(new Expectations() {
            {
                allowing(request).getQuestion();
                will(returnValue(question));

                allowing(question).getType();
                will(returnValue(requestType));
                allowing(question).getName();
                will(returnValue(Name.fromString(requestName)));
            }
        });

        record.log();

        final List<LoggingEvent> events = APPENDER.getEvents();
        assertNotNull(events);
        assertEquals(1, events.size());
        final LoggingEvent event = events.get(0);
        assertEquals(Level.INFO, event.getLevel());
        assertEquals(expected, event.getMessage());
    }

    @Test
    public void testLogWithRequestWithResponse() throws Exception {
        final Date date = new Date(0);
        final InetAddress client = InetAddress.getByName("127.0.0.1");
        final String requestName = "www.foo.com.";
        final int requestType = Type.A;
        final int rcode = Rcode.NOERROR;
        final String rdata = "10.0.0.1";

        final Record question = context.mock(Record.class, "question");
        final Header respHeader = context.mock(Header.class, "respHeader");
        final Record answer = context.mock(Record.class, "answer");

        final String expected = String.format("DNS [01/Jan/1970:00:00:00.000 +0000] %s %s %s %s \"%s\"",
                client.getHostAddress(), Type.string(requestType), requestName, Rcode.string(rcode), rdata);

        record.setRequestDate(date);
        record.setClient(client);
        record.setRequest(request);
        record.setResponse(response);

        context.checking(new Expectations() {
            {
                allowing(request).getQuestion();
                will(returnValue(question));

                allowing(question).getType();
                will(returnValue(requestType));
                allowing(question).getName();
                will(returnValue(Name.fromString(requestName)));

                allowing(response).getHeader();
                will(returnValue(respHeader));
                allowing(response).getSectionArray(Section.ANSWER);
                will(returnValue(new Record[] { answer }));

                allowing(respHeader).getRcode();
                will(returnValue(rcode));

                allowing(answer).rdataToString();
                will(returnValue(rdata));
            }
        });

        record.log();

        final List<LoggingEvent> events = APPENDER.getEvents();
        assertNotNull(events);
        assertEquals(1, events.size());
        final LoggingEvent event = events.get(0);
        assertEquals(Level.INFO, event.getLevel());
        assertEquals(expected, event.getMessage());
    }

    @Test
    public void testLogWithRequestWithResponseNoAnswer() throws Exception {
        final Date date = new Date(0);
        final InetAddress client = InetAddress.getByName("127.0.0.1");
        final String requestName = "www.foo.com.";
        final int requestType = Type.AAAA;
        final int rcode = Rcode.NOERROR;

        final Record question = context.mock(Record.class, "question");
        final Header respHeader = context.mock(Header.class, "respHeader");

        final String expected = String.format("DNS [01/Jan/1970:00:00:00.000 +0000] %s %s %s %s \"%s\"",
                client.getHostAddress(), Type.string(requestType), requestName, Rcode.string(rcode), "-");

        record.setRequestDate(date);
        record.setClient(client);
        record.setRequest(request);
        record.setResponse(response);

        context.checking(new Expectations() {
            {
                allowing(request).getQuestion();
                will(returnValue(question));

                allowing(question).getType();
                will(returnValue(requestType));
                allowing(question).getName();
                will(returnValue(Name.fromString(requestName)));

                allowing(response).getHeader();
                will(returnValue(respHeader));
                allowing(response).getSectionArray(Section.ANSWER);
                will(returnValue(new Record[] {}));

                allowing(respHeader).getRcode();
                will(returnValue(rcode));
            }
        });

        record.log();

        final List<LoggingEvent> events = APPENDER.getEvents();
        assertNotNull(events);
        assertEquals(1, events.size());
        final LoggingEvent event = events.get(0);
        assertEquals(Level.INFO, event.getLevel());
        assertEquals(expected, event.getMessage());
    }

    @BeforeClass
    public static void setUpBeforeClass() throws Exception {
        LOGGER.addAppender(APPENDER);
        LOGGER.setLevel(Level.INFO);
        LOGGER.setAdditivity(false);
    }

    @AfterClass
    public static void tearDownAfterClass() throws Exception {
        LOGGER.removeAppender(APPENDER);
        LOGGER.setLevel(null);
        LOGGER.setAdditivity(true);
    }

}
