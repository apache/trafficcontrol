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

package org.apache.traffic_control.traffic_router.neustar;

import org.apache.traffic_control.traffic_router.neustar.data.NeustarDatabaseUpdater;
import org.apache.traffic_control.traffic_router.geolocation.Geolocation;
import org.apache.traffic_control.traffic_router.geolocation.GeolocationException;
import org.apache.traffic_control.traffic_router.geolocation.GeolocationService;
import com.quova.bff.reader.exception.AddressNotFoundException;
import com.quova.bff.reader.io.GPDatabaseReader;
import com.quova.bff.reader.model.GeoPointResponse;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import java.io.File;
import java.io.IOException;
import java.net.InetAddress;

@Component
public class NeustarGeolocationService implements GeolocationService {
	private static final Logger LOGGER = LogManager.getLogger(NeustarGeolocationService.class);
	private GPDatabaseReader databaseReader;

	@Autowired
	private File neustarDatabaseDirectory;

	@Override
	public Geolocation location(final String ip) throws GeolocationException {
		if (databaseReader == null) {
			return null;
		}

		try {
			GeoPointResponse geoPointResponse = databaseReader.ipInfo(InetAddress.getByName(ip.split("/")[0]));
			return createGeolocation(geoPointResponse);
		} catch (AddressNotFoundException e) {
			return null;
		} catch (Exception e) {
			throw new GeolocationException("Caught exception while attempting to determine location: " + e.getMessage(), e);
		}
	}

	@Override
	public boolean verifyDatabase(final File databaseDirectory) throws IOException {
		throw new RuntimeException("verifyDatabase is no longer allowed, " + NeustarDatabaseUpdater.class.getSimpleName() + " is used for verification instead");
	}

	@Override
	public void reloadDatabase() throws IOException {
		GPDatabaseReader gpDatabaseReader = createDatabaseReader(neustarDatabaseDirectory);

		if (databaseReader != null) {
			databaseReader.close();
		}
		databaseReader = gpDatabaseReader;
	}

	private GPDatabaseReader createDatabaseReader(File databaseDirectory) {
		LOGGER.info("Loading Neustar db: " + databaseDirectory);
		final long t1 = System.currentTimeMillis();
		try {
			GPDatabaseReader gpDatabaseReader = new GPDatabaseReader.Builder(databaseDirectory).build();
			LOGGER.info((System.currentTimeMillis() - t1) + " msec to load Neustar db: " + databaseDirectory);
			return gpDatabaseReader;
		}
		catch (Exception e) {
			String path = databaseDirectory != null ? databaseDirectory.getAbsolutePath() : "NULL";
			LOGGER.error("Database Directory " + path + " is not a valid Neustar database. " + e.getMessage());
		}
		return null;
	}

	// Used by traffic router application context
	public void init() {
	}

	// Used by traffic router application context
	public void destroy() {
		if (databaseReader == null) {
			return;
		}

		try {
			databaseReader.close();
			databaseReader = null;
		} catch (IOException ex) {
			LOGGER.warn("Caught exception while trying to close geolocation database reader: " + ex.getMessage(), ex);
		}
	}

	@Override
	public boolean isInitialized() {
		return databaseReader != null;
	}

	@Override
	public void setDatabaseFile(File file) {
		// Do nothing, this is just here for the interface.
		// The Maxmind version needs it due to how the GeolocationDatabaseUpdater class sets this up.
		// Once TR is running though this isn't going to change.
		// So instead we're just autowiring in the same file for both this and NeustarDatabaseUpdater
	}

	private Geolocation createGeolocation(GeoPointResponse response) {
		Geolocation geolocation = new Geolocation(response.getLatitude(), response.getLongitude());
		geolocation.setCity(response.getCity());
		geolocation.setCountryCode(response.getCountryCode());
		geolocation.setCountryName(response.getCountry());
		geolocation.setPostalCode(response.getPostalCode());

		return geolocation;
	}
}