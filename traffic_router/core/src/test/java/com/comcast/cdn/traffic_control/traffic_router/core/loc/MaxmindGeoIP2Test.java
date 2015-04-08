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

import org.junit.After;
import org.junit.Before;
import org.junit.Test;

import com.comcast.cdn.traffic_control.traffic_router.core.loc.GeolocationException;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.MaxmindGeolocationService;

public class MaxmindGeoIP2Test {
	private MaxmindGeolocationService geoip2;
	private final static String mmdb = "src/test/db/GeoIP2-City.mmdb";

	@Before
	public void setUp() throws Exception {
		this.geoip2 = new MaxmindGeolocationService();
		geoip2.setDatabaseName(mmdb);
		geoip2.init();
	}
	@Test
	public void testSerialLookupPerformance() throws GeolocationException {
		long start = System.currentTimeMillis();
		int total = 100000;

		for (int i = 0; i <= total; i++) {
			geoip2.location("10.0.0.1");
		}

		long duration = System.currentTimeMillis() - start;
		double tps = (double) total / ((double) duration / 1000);

		System.out.println("MaxMind2 lookup duration: " + duration + "ms, " + tps + " tps");
	}

	@After
	public void tearDown() throws Exception {
		geoip2.destroy();
	}
}
