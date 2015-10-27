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

import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryService;
import com.comcast.cdn.traffic_control.traffic_router.core.request.HTTPRequest;
import com.comcast.cdn.traffic_control.traffic_router.core.router.HTTPRouteResult;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker;
import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouter;
import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouterManager;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import java.util.Date;
import java.util.HashMap;
import java.util.HashSet;
import java.util.Map;
import java.util.Set;
import java.util.Vector;

import static org.mockito.Matchers.any;
import static org.mockito.Mockito.doReturn;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.spy;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;
import static org.powermock.api.mockito.PowerMockito.whenNew;

@RunWith(PowerMockRunner.class)
@PrepareForTest({TRServlet.class, HTTPAccessRecord.Builder.class, HTTPAccessRecord.class})
public class TRServletTest {

    @Test
    public void itAddsRequestHeadersToAccessLog() throws Exception {

        Set<String> headerNames = new HashSet<String>();
        headerNames.add("If-Modified-Since");

        HttpServletResponse servletResponse = mock(HttpServletResponse.class);

        HttpServletRequest servletRequest = mock(HttpServletRequest.class);
        when(servletRequest.getRequestURL()).thenReturn(new StringBuffer("blah"));
        when(servletRequest.getHeaderNames()).thenReturn(new Vector(headerNames).elements());

        when(servletRequest.getHeader("If-Modified-Since")).thenReturn("Thurs, 15 July 2010 12:00:00 UTC");

        HTTPAccessRecord.Builder builder = new HTTPAccessRecord.Builder(new Date(), servletRequest);
        HTTPAccessRecord httpAccessRecord = builder.build();

        builder = spy(builder);
        doReturn(httpAccessRecord).when(builder).build();

        whenNew(HTTPAccessRecord.Builder.class).withArguments(any(Date.class), any(HttpServletRequest.class)).thenReturn(builder);
        whenNew(HTTPAccessRecord.Builder.class).withArguments(any(HTTPAccessRecord.class)).thenReturn(builder);


        TrafficRouter trafficRouter = mock(TrafficRouter.class);

        TrafficRouterManager trafficRouterManager = mock(TrafficRouterManager.class);
        when(trafficRouterManager.getTrafficRouter()).thenReturn(trafficRouter);

        Set<String> hnames = new HashSet<String>();;
        hnames.add("If-Modified-Since");

        DeliveryService deliveryService = mock(DeliveryService.class);
        when(deliveryService.getRequestHeaders()).thenReturn(hnames);
        when(deliveryService.getResponseHeaders()).thenReturn(new HashMap<String, String>());

        HTTPRouteResult httpRouteResult = mock(HTTPRouteResult.class);
        when(httpRouteResult.getDeliveryService()).thenReturn(deliveryService);

        when(trafficRouter.route(any(HTTPRequest.class), any(StatTracker.Track.class))).thenReturn(httpRouteResult);


        TRServlet trServlet = new TRServlet();
        trServlet.setTrafficRouterManager(trafficRouterManager);
        trServlet.setStatTracker(new StatTracker());
        trServlet.doGet(servletRequest, servletResponse);

        verify(builder).requestHeaders(any(Map.class));
    }

}
