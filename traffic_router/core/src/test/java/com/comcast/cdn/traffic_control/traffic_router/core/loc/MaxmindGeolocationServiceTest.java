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

import static org.junit.Assert.assertEquals;

import java.io.IOException;

import org.jmock.Expectations;
import org.jmock.Mockery;
import org.jmock.integration.junit4.JMock;
import org.jmock.integration.junit4.JUnit4Mockery;
import org.jmock.lib.legacy.ClassImposteriser;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;

import com.comcast.cdn.traffic_control.traffic_router.core.loc.MaxmindGeolocationService;
import com.maxmind.geoip2.DatabaseReader;

@RunWith(JMock.class)
public class MaxmindGeolocationServiceTest {
	private static final String DATABASE_URL = "database url";

	private final Mockery context = new JUnit4Mockery() {
		{
			setImposteriser(ClassImposteriser.INSTANCE);
		}
	};

	private TestService service;
	private DatabaseReader databaseReader;

	@Before
	public void setUp() throws Exception {
		databaseReader = context.mock(DatabaseReader.class);

		service = new TestService();
		service.setDatabaseName(DATABASE_URL);
		service.init();
	}

	@Test
	public void testDestroy() throws IOException {
		context.checking(new Expectations() {
			{
				oneOf(databaseReader).close();
			}
		});
		service.destroy();
		assertEquals(1, service.numCalls);
	}

//	@Test
//	public void testLocationFound() throws Exception {
//		final String ip = "10.0.0.1";
//
//		final Location loc = new Location();
//		loc.latitude = 10f;
//		loc.longitude = 10f;
//
//		context.checking(new Expectations() {
//			{
//				oneOf(lookupService).getLocation(ip);
//				will(returnValue(loc));
//			}
//		});
//		final Geolocation actual = service.location(ip);
//		assertEquals(loc.latitude, actual.getLatitude(), 0.01);
//		assertEquals(loc.longitude, actual.getLongitude(), 0.01);
//	}

	@Test
	public void testReloadDatabase() throws Exception {
		context.checking(new Expectations() {
			{
				oneOf(databaseReader).close();
			}
		});
		service.reloadDatabase();
		assertEquals(2, service.numCalls);
	}

	private class TestService extends MaxmindGeolocationService {

		private int numCalls = 0;

		@Override
		protected DatabaseReader createDatabaseReader() throws IOException {
			numCalls++;
			return databaseReader;
		}

	}

}
