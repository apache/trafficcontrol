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

package geolocation;

import org.apache.traffic_control.traffic_router.neustar.NeustarGeolocationService;
import org.apache.traffic_control.traffic_router.neustar.data.NeustarDatabaseUpdater;
import org.apache.traffic_control.traffic_router.geolocation.Geolocation;
import com.maxmind.db.Reader;
import com.quova.bff.reader.io.GPDatabaseReader;
import com.quova.bff.reader.model.GeoPointResponse;
import org.apache.logging.log4j.Appender;
import org.apache.logging.log4j.LogManager;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import java.io.File;
import java.net.InetAddress;

import static org.hamcrest.CoreMatchers.notNullValue;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.core.IsNull.nullValue;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;
import static org.mockito.MockitoAnnotations.initMocks;
import static org.powermock.api.mockito.PowerMockito.whenNew;
import static org.springframework.test.util.AssertionErrors.fail;

@RunWith(PowerMockRunner.class)
@PrepareForTest({NeustarGeolocationService.class, GPDatabaseReader.class, Reader.class})
@PowerMockIgnore("javax.management.*")
public class NeustarGeolocationServiceTest {
	@Mock
	File neustarDatabaseDirectory;

	@InjectMocks
	NeustarGeolocationService service = new NeustarGeolocationService();

	@Before
	public void before() throws Exception {
		// This prevents extraneous output about 'WARN No appenders could be found....'
		LogManager.getRootLogger().addAppender(mock(Appender.class));

		initMocks(this);
		service.init();
	}

	@Test
	public void itNoLongerAllowsVerifyDatabase() throws Exception {
		try {
			service.verifyDatabase(neustarDatabaseDirectory);
			fail("Should have thrown RuntimeException when calling verifyDatabase");
		} catch (RuntimeException e) {
			assertThat(e.getMessage(), equalTo("verifyDatabase is no longer allowed, " + NeustarDatabaseUpdater.class.getSimpleName() + " is used for verification instead"));
		}
	}

	@Test
	public void itReturnsNullWhenDatabaseNotLoaded() throws Exception {
		assertThat(service.isInitialized(), equalTo(false));
		assertThat(service.location("192.168.99.100"), nullValue());
	}

	@Test
	public void itReturnsNullWhenDatabaseDoesNotExist() throws Exception {
		when(neustarDatabaseDirectory.getAbsolutePath()).thenReturn("/path/to/file/");

		assertThat(service.isInitialized(), equalTo(false));
		assertThat(service.location("192.168.99.100"), nullValue());

		service.reloadDatabase();

		assertThat(service.isInitialized(), equalTo(false));
		assertThat(service.location("192.168.99.100"), nullValue());
	}

	@Test
	@PrepareForTest({GPDatabaseReader.Builder.class, NeustarGeolocationService.class})
	public void itReturnsALocationWhenTheDatabaseIsLoaded() throws Exception {
		when(neustarDatabaseDirectory.exists()).thenReturn(true);
		when(neustarDatabaseDirectory.list()).thenReturn(new String[] {"foo.gpdb"});

		GeoPointResponse geoPointResponse = mock(GeoPointResponse.class);
		when(geoPointResponse.getCity()).thenReturn("Springfield");
		when(geoPointResponse.getLatitude()).thenReturn(40.0);
		when(geoPointResponse.getLongitude()).thenReturn(-105.0);
		when(geoPointResponse.getCountry()).thenReturn("United States");
		when(geoPointResponse.getCountryCode()).thenReturn("100");

		GPDatabaseReader gpDatabaseReader = mock(GPDatabaseReader.class);
		when(gpDatabaseReader.ipInfo(InetAddress.getByName("192.168.99.100"))).thenReturn(geoPointResponse);

		GPDatabaseReader.Builder builder = mock(GPDatabaseReader.Builder.class);
		when(builder.build()).thenReturn(gpDatabaseReader);
		whenNew(GPDatabaseReader.Builder.class).withArguments(neustarDatabaseDirectory).thenReturn(builder);

		service.reloadDatabase();

		assertThat(service.isInitialized(), equalTo(true));

		Geolocation geolocation = service.location("192.168.99.100");
		assertThat(geolocation.getCity(), equalTo("Springfield"));
		assertThat(geolocation.getLatitude(), equalTo(40.0));
		assertThat(geolocation.getLongitude(), equalTo(-105.0));
		assertThat(geolocation.getCountryName(), equalTo("United States"));
		assertThat(geolocation.getCountryCode(), equalTo("100"));


		assertThat(service.location("192.168.99.100"), notNullValue());
	}

	@After
	public void after() {
		service.destroy();
	}
}
