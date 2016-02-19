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

import java.io.File;
import java.io.IOException;
import java.net.InetAddress;

import org.apache.log4j.Logger;

import com.maxmind.geoip2.DatabaseReader;
import com.maxmind.geoip2.exception.AddressNotFoundException;
import com.maxmind.geoip2.model.CityResponse;

public class MaxmindGeolocationService implements GeolocationService {
	private static final Logger LOGGER = Logger.getLogger(MaxmindGeolocationService.class);
	private boolean initialized = false;
	private DatabaseReader databaseReader;
	private File databaseFile;

	private CityResponse getCityResponse(final String address) throws GeolocationException {
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

		final CityResponse response = getCityResponse(ip.split("/")[0]);

		return (isResponseValid(response)) ? new Geolocation(response) : null;
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
		}
	}

	@Override
	public void verifyDatabase(final File databaseFile) throws IOException {
		databaseReader = createDatabaseReader(databaseFile);
		this.databaseFile = databaseFile;
	}

	private DatabaseReader createDatabaseReader(final File databaseFile) throws IOException {
		if (!databaseFile.exists()) {
			LOGGER.warn(databaseFile + " does not exist yet!");
			return null;
		}

		LOGGER.info("Loading MaxMind db: " + databaseFile);
		final DatabaseReader reader = new DatabaseReader.Builder(databaseFile).build();
		initialized = true;
		return reader;
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
}