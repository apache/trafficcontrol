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

package com.comcast.cdn.traffic_control.traffic_router.core.router;

import com.comcast.cdn.traffic_control.traffic_router.core.util.IntegrationTest;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.junit.Before;
import org.junit.BeforeClass;
import org.junit.Test;
import org.junit.experimental.categories.Category;
import org.springframework.context.ApplicationContext;

import com.comcast.cdn.traffic_control.traffic_router.core.TestBase;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.GeolocationDatabaseUpdater;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.NetworkUpdater;
import com.comcast.cdn.traffic_control.traffic_router.core.request.HTTPRequest;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker.Track;

import java.nio.file.Files;
import java.nio.file.Paths;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;

@Category(IntegrationTest.class)
public class StatelessTrafficRouterTest {
	private static final Logger LOGGER = LogManager.getLogger(StatelessTrafficRouterTest.class);
	private TrafficRouterManager trafficRouterManager;
	private GeolocationDatabaseUpdater geolocationDatabaseUpdater;
	private NetworkUpdater networkUpdater;
	private static ApplicationContext context;

	@BeforeClass
	public static void setUpBeforeClass() throws Exception {
		assertThat("Copy core/src/main/conf/traffic_monitor.properties to core/src/test/conf and set 'traffic_monitor.bootstrap.hosts' to a real traffic monitor", Files.exists(Paths.get(TestBase.monitorPropertiesPath)), equalTo(true));
		context = TestBase.getContext();
	}

	@Before
	public void setUp() throws Exception {
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
	}

	@Test
	public void testRouteHTTPRequestTrack() throws Exception {
		HTTPRequest req = new HTTPRequest();
		req.setClientIP("10.0.0.1");
		req.setPath("/QualityLevels(96000)/Fragments(audio_eng=20720000000)");
		req.setQueryString("");
		req.setHostname("somehost.cdn.net");
		req.setRequestedUrl("http://somehost.cdn.net/QualityLevels(96000)/Fragments(audio_eng=20720000000)");
		Track track = StatTracker.getTrack();
		trafficRouterManager.getTrafficRouter().route(req, track);
	}

}
