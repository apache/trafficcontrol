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

package com.comcast.cdn.traffic_control.traffic_router.core.router;

import java.net.URL;

import org.apache.log4j.Logger;
import org.junit.After;
import org.junit.AfterClass;
import org.junit.Before;
import org.junit.BeforeClass;
import org.junit.Test;
import org.springframework.context.ApplicationContext;

import com.comcast.cdn.traffic_control.traffic_router.core.TestBase;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.GeolocationDatabaseUpdater;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.NetworkUpdater;
import com.comcast.cdn.traffic_control.traffic_router.core.request.HTTPRequest;
import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouterManager;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker.Track;

public class StatelessTrafficRouterTest {
	private static final Logger LOGGER = Logger.getLogger(StatelessTrafficRouterTest.class);
	private TrafficRouterManager trafficRouterManager;
	private GeolocationDatabaseUpdater geolocationDatabaseUpdater;
	private NetworkUpdater networkUpdater;
	private static ApplicationContext context;

	@BeforeClass
	public static void setUpBeforeClass() throws Exception {
		try {
			context = TestBase.getContext();
			Logger root = Logger.getRootLogger();
			boolean rootIsConfigured = root.getAllAppenders().hasMoreElements();
			System.out.println("rootIsConfigured: "+rootIsConfigured);
		} catch(Exception e) {
			e.printStackTrace();
		}
	}

	@AfterClass
	public static void tearDownAfterClass() throws Exception {
	}

	@Before
	public void setUp() throws Exception {
		try {
			trafficRouterManager = (TrafficRouterManager) context.getBean("trafficRouterManager");
			geolocationDatabaseUpdater = (GeolocationDatabaseUpdater) context.getBean("geolocationDatabaseUpdater");
			networkUpdater = (NetworkUpdater) context.getBean("networkUpdater");

			while (!networkUpdater.isLoaded()) {
				LOGGER.info("Waiting for a valid location database before proceeding");
				Thread.sleep(1000);
			}

			while (!geolocationDatabaseUpdater.isLoaded()) {
				LOGGER.info("Waiting for a valid Maxmind database before proceeding");
				Thread.sleep(1000);
			}
		} catch(Exception e) {
			e.printStackTrace();
		}
	}

	@After
	public void tearDown() throws Exception {
	}

	@Test
	public void testRouteDNSRequestTrack() {
		//		fail("Not yet implemented");
	}

	@Test
	public void testRouteHTTPRequestTrack() {
		HTTPRequest req = new HTTPRequest();
		req.setClientIP("10.0.0.1");
		req.setPath("/QualityLevels(96000)/Fragments(audio_eng=20720000000)");
		req.setQueryString("");
		req.setHostname("somehost.cdn.net");
		req.setRequestedUrl("http://somehost.cdn.net/QualityLevels(96000)/Fragments(audio_eng=20720000000)");
		Track track = StatTracker.getTrack();
		try {
			HTTPRouteResult routeResult = trafficRouterManager.getTrafficRouter().route(req, track);
			if (routeResult == null) {
//				fail("HTTP route returned null");
				System.out.println("HTTP route returned null");
			} else {
				System.out.println(routeResult.getUrl());
			}
		} catch (Exception e2) {
			e2.printStackTrace();
//			fail(e2.toString());
		}
	}

	@Test
	public void testConsistentHash() {
		//		fail("Not yet implemented");
	}

	@Test
	public void testSelectDeliveryService() {
		//		fail("Not yet implemented");
	}

}
