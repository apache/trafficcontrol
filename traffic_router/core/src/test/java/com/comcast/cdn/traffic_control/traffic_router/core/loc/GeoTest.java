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

package com.comcast.cdn.traffic_control.traffic_router.core.loc;

import com.comcast.cdn.traffic_control.traffic_router.core.util.IntegrationTest;
import org.apache.log4j.Logger;
import org.junit.Before;
import org.junit.BeforeClass;
import org.junit.Test;
import org.junit.experimental.categories.Category;
import org.springframework.context.ApplicationContext;

import org.junit.Assert;

import com.comcast.cdn.traffic_control.traffic_router.core.TestBase;
import com.comcast.cdn.traffic_control.traffic_router.geolocation.Geolocation;

@Category(IntegrationTest.class)
public class GeoTest {
	private static final Logger LOGGER = Logger.getLogger(GeoTest.class);

	private GeolocationDatabaseUpdater geolocationDatabaseUpdater;
	private MaxmindGeolocationService geolocationService;
	private static ApplicationContext context;

	@BeforeClass
	public static void setUpBeforeClass() throws Exception {
		try {
			context = TestBase.getContext();
		} catch(Exception e) {
			e.printStackTrace();
		}
	}

	@Before
	public void setUp() throws Exception {
		geolocationDatabaseUpdater = (GeolocationDatabaseUpdater) context.getBean("geolocationDatabaseUpdater");
		geolocationService = (MaxmindGeolocationService) context.getBean("GeolocationService");

		while (!geolocationDatabaseUpdater.isLoaded()) {
			LOGGER.info("Waiting for a valid Maxmind database before proceeding");
			Thread.sleep(1000);
		}

	}

	@Test
	public void testIps() {
		try {
			final String testips[][] = {
					{"40.40.40.40","cache-group-1"}
			};
			for(int i = 0; i < testips.length; i++) {
				Geolocation location = geolocationService.location(testips[i][0]);
				Assert.assertNotNull(location);
				String loc = location.toString();
				LOGGER.info(String.format("result for ip=%s: %s\n",testips[i], loc));
			}
		} catch (Exception e) {
			e.printStackTrace();
		}
	}
}
