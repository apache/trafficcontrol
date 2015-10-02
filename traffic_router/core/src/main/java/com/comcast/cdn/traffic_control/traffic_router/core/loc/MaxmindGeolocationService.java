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
import java.util.concurrent.locks.ReadWriteLock;
import java.util.concurrent.locks.ReentrantReadWriteLock;

import org.apache.log4j.Logger;

import com.maxmind.geoip2.DatabaseReader;
import com.maxmind.geoip2.exception.AddressNotFoundException;
import com.maxmind.geoip2.model.CityResponse;

public class MaxmindGeolocationService implements GeolocationService {
	private static final Logger LOGGER = Logger.getLogger(MaxmindGeolocationService.class);
	private final ReadWriteLock lock = new ReentrantReadWriteLock();
	private String databaseName;
	private DatabaseReader databaseReader;
	private boolean initialized = false;

	@Override
	@SuppressWarnings("PMD.EmptyCatchBlock")
	public Geolocation location(final String ip) throws GeolocationException {
		lock.readLock().lock();

		try {
			if (databaseReader != null) {
				final String[] parts = ip.split("/");
				final InetAddress address = InetAddress.getByName(parts[0]);
				final CityResponse response = databaseReader.city(address);

				if (isResponseValid(response)) {
					return new Geolocation(response);
				}
			}
		} catch (AddressNotFoundException ex) {
			// this is fine; we'll just return null below and send them to Chicago
		} catch (Exception ex) {
			LOGGER.error(ex, ex);
			throw new GeolocationException("Caught exception while attempting to determine location: " + ex.getMessage(), ex);
		} finally {
			lock.readLock().unlock();
		}

		return null;
	}

	private boolean isResponseValid(final CityResponse response) {
		if (response == null) {
			return false;
		} else if (response.getLocation() == null) {
			return false;
		} else if (response.getLocation().getLatitude() == null) {
			return false;
		} else if (response.getLocation().getLongitude() == null) {
			return false;
		}

		return true;
	}

	protected DatabaseReader createDatabaseReader() throws IOException {
		final File database = new File(getDatabaseName());
		if (database.exists()) {
			LOGGER.info("Loading MaxMind db: " + database);
			final DatabaseReader reader = new DatabaseReader.Builder(database).build();
			setInitialized(true);
			return reader;
		} else {
			LOGGER.warn(database + " does not exist yet!");
			return null;
		}
	}

	public void init() {
		lock.writeLock().lock();

		try {
			databaseReader = createDatabaseReader();
		} catch (final IOException ex) {
			LOGGER.fatal("Caught exception while trying to open geolocation database " + getDatabaseName() + ": " + ex.getMessage(), ex);
		} finally {
			lock.writeLock().unlock();
		}
	}

	public void destroy() {
		lock.writeLock().lock();

		try {
			if (databaseReader != null) {
				databaseReader.close();
				databaseReader = null;
			}
		} catch (IOException ex) {
			LOGGER.warn("Caught exception while trying to close geolocation database reader: " + ex.getMessage(), ex);
		} finally {
			lock.writeLock().unlock();
		}
	}

	@Override
	public void reloadDatabase() throws IOException {
		lock.writeLock().lock();

		try {
			if (databaseReader != null) {
				databaseReader.close();
			}

			databaseReader = createDatabaseReader();
		} finally {
			lock.writeLock().unlock();
		}
	}

	@Override
	public void verifyDatabase(final File dbFile) throws IOException {
		LOGGER.info("Attempting to verify " + dbFile.getAbsolutePath());
		final DatabaseReader dbr = new DatabaseReader.Builder(dbFile).build();
		dbr.close();
	}

	public String getDatabaseName() {
		return databaseName;
	}

	public void setDatabaseName(final String databaseName) {
		this.databaseName = databaseName;
	}

	@Override
	public boolean isInitialized() {
		return initialized;
	}

	private void setInitialized(final boolean initialized) {
		this.initialized = initialized;
	}
}