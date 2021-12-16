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

package org.apache.traffic_control.traffic_router.core.loc;

import com.maxmind.geoip2.DatabaseReader;
import com.maxmind.geoip2.model.CityResponse;
import com.maxmind.geoip2.record.Location;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.api.mockito.PowerMockito;
import org.powermock.core.classloader.annotations.PowerMockIgnore;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import java.io.File;
import java.net.InetAddress;

import static org.hamcrest.CoreMatchers.equalTo;
import static org.hamcrest.CoreMatchers.notNullValue;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.core.IsNull.nullValue;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;
import static org.powermock.api.mockito.PowerMockito.whenNew;

@RunWith(PowerMockRunner.class)
@PowerMockIgnore("javax.management.*")
public class MaxmindGeolocationServiceTest {

	private MaxmindGeolocationService service;

	@Before
	public void before() throws Exception {
		service = new MaxmindGeolocationService();
		service.init();
	}

	@Test
	public void itReturnsNullWhenDatabaseNotLoaded() throws Exception {
		assertThat(service.isInitialized(), equalTo(false));
		assertThat(service.location("192.168.99.100"), nullValue());
	}

	@Test
	public void itReturnsNullWhenDatabaseDoesNotExist() throws Exception {
		service.verifyDatabase(mock(File.class));
		assertThat(service.isInitialized(), equalTo(false));
		assertThat(service.location("192.168.99.100"), nullValue());
		service.reloadDatabase();
		assertThat(service.isInitialized(), equalTo(false));
		assertThat(service.location("192.168.99.100"), nullValue());
	}

	@PrepareForTest({MaxmindGeolocationService.class, DatabaseReader.Builder.class, Location.class, CityResponse.class})
	@Test
	public void itReturnsALocationWhenTheDatabaseIsLoaded() throws Exception {
		File databaseFile = mock(File.class);
		when(databaseFile.exists()).thenReturn(true);

		Location location = PowerMockito.mock(Location.class);
		when(location.getLatitude()).thenReturn(40.0);
		when(location.getLongitude()).thenReturn(-105.0);

		CityResponse cityResponse = PowerMockito.mock(CityResponse.class);
		when(cityResponse.getLocation()).thenReturn(location);

		DatabaseReader databaseReader = mock(DatabaseReader.class);
		when(databaseReader.city(InetAddress.getByName("192.168.99.100"))).thenReturn(cityResponse);

		DatabaseReader.Builder builder = mock(DatabaseReader.Builder.class);
		when(builder.build()).thenReturn(databaseReader);

		whenNew(DatabaseReader.Builder.class).withArguments(databaseFile).thenReturn(builder);
		service.setDatabaseFile(databaseFile);
		service.reloadDatabase();

		assertThat(service.isInitialized(), equalTo(true));

		assertThat(service.location("192.168.99.100"), notNullValue());
	}

	@After
	public void after() {
		service.destroy();
	}
}
