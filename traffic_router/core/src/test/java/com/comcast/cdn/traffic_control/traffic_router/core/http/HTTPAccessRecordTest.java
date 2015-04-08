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

package com.comcast.cdn.traffic_control.traffic_router.core.http;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;

import java.net.URL;
import java.util.Date;
import java.util.List;

import javax.servlet.http.HttpServletRequest;

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

import com.comcast.cdn.traffic_control.traffic_router.core.http.HTTPAccessRecord;
import com.comcast.cdn.traffic_control.traffic_router.logger.NoLogAppender;

@RunWith(JMock.class)
public class HTTPAccessRecordTest {
    private static final Logger LOGGER = Logger.getLogger("com.comcast.cdn.traffic_control.traffic_router.core.access");
    private static final NoLogAppender APPENDER = new NoLogAppender();

    private final Mockery context = new JUnit4Mockery() {
        {
            setImposteriser(ClassImposteriser.INSTANCE);
        }
    };

    private HTTPAccessRecord access;
    private HttpServletRequest request;

    @Before
    public void setUp() throws Exception {
        access = new HTTPAccessRecord();
        request = context.mock(HttpServletRequest.class);
        APPENDER.clear();
    }

    @Test
    public void testLog302NoQueryString() throws Exception {
        final Date date = new Date(0);
        final String ip = "127.0.0.1";
        final String requestURL = "http://foo.com/stuff";
        final int responseCode = 302;
        final URL responseURL = new URL("http://bar.com/stuff");
        final String expected = String.format("HTTP [01/Jan/1970:00:00:00.000 +0000] %s %s %s %s", ip, requestURL,
                String.valueOf(responseCode), responseURL);
        access.setRequestDate(date);
        access.setRequest(request);
        access.setResponseCode(responseCode);
        access.setResponseURL(responseURL);
        context.checking(new Expectations() {
            {
                allowing(request).getRemoteAddr();
                will(returnValue(ip));
                allowing(request).getRequestURL();
                will(returnValue(new StringBuffer(requestURL)));
                allowing(request).getQueryString();
                will(returnValue(null));
            }
        });
        access.log();
        final List<LoggingEvent> events = APPENDER.getEvents();
        assertNotNull(events);
        assertEquals(1, events.size());
        final LoggingEvent event = events.get(0);
        assertEquals(Level.INFO, event.getLevel());
        assertEquals(expected, event.getMessage());
    }

    @Test
    public void testLog302WithQueryString() throws Exception {
        final Date date = new Date(0);
        final String ip = "127.0.0.1";
        final String requestURL = "http://foo.com/stuff";
        final String queryString = "foo=bar&stuff=lots";
        final int responseCode = 302;
        final URL responseURL = new URL("http://bar.com/stuff");
        final String expected = String.format("HTTP [01/Jan/1970:00:00:00.000 +0000] %s %s?%s %s %s", ip, requestURL,
                queryString, String.valueOf(responseCode), responseURL);
        access.setRequestDate(date);
        access.setRequest(request);
        access.setResponseCode(responseCode);
        access.setResponseURL(responseURL);
        context.checking(new Expectations() {
            {
                allowing(request).getRemoteAddr();
                will(returnValue(ip));
                allowing(request).getRequestURL();
                will(returnValue(new StringBuffer(requestURL)));
                allowing(request).getQueryString();
                will(returnValue(queryString));
            }
        });
        access.log();
        final List<LoggingEvent> events = APPENDER.getEvents();
        assertNotNull(events);
        assertEquals(1, events.size());
        final LoggingEvent event = events.get(0);
        assertEquals(Level.INFO, event.getLevel());
        assertEquals(expected, event.getMessage());
    }

    @Test
    public void testLog503NoQueryString() throws Exception {
        final Date date = new Date(0);
        final String ip = "127.0.0.1";
        final String requestURL = "http://foo.com/stuff";
        final int responseCode = 503;
        final URL responseURL = null;
        final String expected = String.format("HTTP [01/Jan/1970:00:00:00.000 +0000] %s %s %s %s", ip, requestURL,
                String.valueOf(responseCode), "-");
        access.setRequestDate(date);
        access.setRequest(request);
        access.setResponseCode(responseCode);
        access.setResponseURL(responseURL);
        context.checking(new Expectations() {
            {
                allowing(request).getRemoteAddr();
                will(returnValue(ip));
                allowing(request).getRequestURL();
                will(returnValue(new StringBuffer(requestURL)));
                allowing(request).getQueryString();
                will(returnValue(null));
            }
        });
        access.log();
        final List<LoggingEvent> events = APPENDER.getEvents();
        assertNotNull(events);
        assertEquals(1, events.size());
        final LoggingEvent event = events.get(0);
        assertEquals(Level.INFO, event.getLevel());
        assertEquals(expected, event.getMessage());
    }

    @Test
    public void testLogErrorNoQueryString() throws Exception {
        final Date date = new Date(0);
        final String ip = "127.0.0.1";
        final String requestURL = "http://foo.com/stuff";
        final int responseCode = -1;
        final URL responseURL = null;
        final String expected = String.format("HTTP [01/Jan/1970:00:00:00.000 +0000] %s %s %s %s", ip, requestURL, "-",
                "-");
        access.setRequestDate(date);
        access.setRequest(request);
        access.setResponseCode(responseCode);
        access.setResponseURL(responseURL);
        context.checking(new Expectations() {
            {
                allowing(request).getRemoteAddr();
                will(returnValue(ip));
                allowing(request).getRequestURL();
                will(returnValue(new StringBuffer(requestURL)));
                allowing(request).getQueryString();
                will(returnValue(null));
            }
        });
        access.log();
        final List<LoggingEvent> events = APPENDER.getEvents();
        assertNotNull(events);
        assertEquals(1, events.size());
        final LoggingEvent event = events.get(0);
        assertEquals(Level.INFO, event.getLevel());
        assertEquals(expected, event.getMessage());
    }

    @BeforeClass
    public static void setUpBeforeClass() {
        LOGGER.addAppender(APPENDER);
        LOGGER.setLevel(Level.INFO);
        LOGGER.setAdditivity(false);
    }

    @AfterClass
    public static void tearDownAfterClass() {
        LOGGER.removeAppender(APPENDER);
        LOGGER.setLevel(null);
        LOGGER.setAdditivity(true);
    }
}
