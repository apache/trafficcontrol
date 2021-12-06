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

import java.io.File;
import java.io.IOException;
import java.net.InetAddress;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import com.maxmind.geoip2.DatabaseReader;
import com.maxmind.geoip2.exception.AddressNotFoundException;
import com.maxmind.geoip2.model.CityResponse;

import org.apache.traffic_control.traffic_router.geolocation.Geolocation;
import org.apache.traffic_control.traffic_router.geolocation.GeolocationException;
import org.apache.traffic_control.traffic_router.geolocation.GeolocationService;

public class MaxmindGeolocationService implements GeolocationService {
	private static final Logger LOGGER = LogManager.getLogger(MaxmindGeolocationService.class);
	private boolean initialized = false;
	private DatabaseReader databaseReader;
	private File databaseFile;

	private CityResponse getCityResponse(final DatabaseReader databaseReader, final String address) throws GeolocationException {
		try {
			return databaseReader.city(InetAddress.getByName(address));
		} catch (AddressNotFoundException e) {
			return null;
		} catch (Exception e) {
			throw new GeolocationException("Caught exception while attempting to determine location: " + e.getMessage(), e);
		}
	}

	@Override
	public Geolocation location(final String ip) throws GeolocationException {
		if (databaseReader == null) {
			return null;
		}

		final CityResponse response = getCityResponse(databaseReader, ip.split("/")[0]);

		return (isResponseValid(response)) ? createGeolocation(response) : null;
	}

	private boolean isResponseValid(final CityResponse response) {
		return (response != null && response.getLocation() != null &&
			response.getLocation().getLatitude() != null && response.getLocation().getLongitude() != null);
	}

	@Override
	public void reloadDatabase() throws IOException {
		if (databaseReader != null) {
			databaseReader.close();
		}

		if (databaseFile != null) {
			databaseReader = createDatabaseReader(databaseFile);
			if (databaseReader != null) {
				initialized = true;
			}
		}
	}

	@Override
	public boolean verifyDatabase(final File databaseFile) throws IOException {
		return createDatabaseReader(databaseFile) != null;
	}

	@Override
	public void setDatabaseFile(final File databaseFile) {
		this.databaseFile = databaseFile;
	}

	@SuppressWarnings("PMD.AvoidUsingHardCodedIP")
	private DatabaseReader createDatabaseReader(final File databaseFile) throws IOException {
		if (!databaseFile.exists()) {
			LOGGER.warn(databaseFile.getAbsolutePath() + " does not exist yet!");
			return null;
		}

		if (databaseFile.isDirectory()) {
			LOGGER.error(databaseFile + " is a directory, need a file");
			return null;
		}

		LOGGER.info("Loading MaxMind db: " + databaseFile.getAbsolutePath());

		try {
			final DatabaseReader reader = new DatabaseReader.Builder(databaseFile).build();
			getCityResponse(reader, "127.0.0.1");
			return reader;
		} catch (Exception e) {
			LOGGER.error(databaseFile.getAbsolutePath() + " is not a valid Maxmind data file.  " + e.getMessage());
			return null;
		}

	}

	@Override
	public boolean isInitialized() {
		return initialized;
	}

	public void init() {
	}

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

	public Geolocation createGeolocation(final CityResponse response) {
		final double latitude = response.getLocation().getLatitude();
		final double longitude = response.getLocation().getLongitude();

		final Geolocation geolocation = new Geolocation(latitude, longitude);
		if (response.getPostal() != null) {
			geolocation.setPostalCode(response.getPostal().getCode());
		}

		if (response.getCity() != null) {
			geolocation.setCity(response.getCity().getName());
		}

		if (response.getCountry() != null) {
			geolocation.setCountryCode(response.getCountry().getIsoCode());
			geolocation.setCountryName(response.getCountry().getName());
		}

		if (geolocation.getCity() == null && geolocation.getPostalCode() == null && response.getSubdivisions().isEmpty()) {
			geolocation.setDefaultLocation(true);
		}

		return geolocation;
	}
}