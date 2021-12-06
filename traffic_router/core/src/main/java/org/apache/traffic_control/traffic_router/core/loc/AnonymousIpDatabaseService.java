/*
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
import com.maxmind.geoip2.exception.GeoIp2Exception;
import com.maxmind.geoip2.model.AnonymousIpResponse;

@SuppressWarnings({ "PMD.AvoidDuplicateLiterals" })
public class AnonymousIpDatabaseService {
	private static final Logger LOGGER = LogManager.getLogger(AnonymousIpDatabaseService.class);

	private boolean initialized = false;
	private File databaseFile;
	private DatabaseReader databaseReader;

	/*
	 * Reloads the anonymous ip database
	 */
	public void reloadDatabase() throws IOException {
		if (databaseReader != null) {
			databaseReader.close();
		}

		if (databaseFile != null) {
			databaseReader = createDatabaseReader(databaseFile);
			if (databaseReader != null) {
				initialized = true;
			} else {
				throw new IOException("Could not create database reader");
			}
		}
	}

	public void setDatabaseFile(final File databaseFile) {
		this.databaseFile = databaseFile;
	}

	/*
	 * Verifies the database by attempting to recreate it
	 */
	public boolean verifyDatabase(final File databaseFile) throws IOException {
		return createDatabaseReader(databaseFile) != null;
	}

	/*
	 * Creates a DatabaseReader object using an input database file
	 */
	private DatabaseReader createDatabaseReader(final File databaseFile) throws IOException {
		if (!databaseFile.exists()) {
			LOGGER.warn(databaseFile.getAbsolutePath() + " does not exist yet!");
			return null;
		}

		if (databaseFile.isDirectory()) {
			LOGGER.error(databaseFile + " is a directory, need a file");
			return null;
		}

		LOGGER.info("Loading Anonymous IP db: " + databaseFile.getAbsolutePath());

		try {
			final DatabaseReader reader = new DatabaseReader.Builder(databaseFile).build();
			return reader;
		} catch (Exception e) {
			LOGGER.error(databaseFile.getAbsolutePath() + " is not a valid Anonymous IP data file", e);
			return null;
		}
	}

	/*
	 * Returns an AnonymousIpResponse from looking an ip up in the database
	 */
	public AnonymousIpResponse lookupIp(final InetAddress ipAddress) {
		if (initialized) {
			// Return an anonymousIp object after looking up the ip in the
			// database
			try {
				return databaseReader.anonymousIp(ipAddress);
			} catch (GeoIp2Exception e) {
				LOGGER.debug(String.format("AnonymousIP: IP %s not found in anonymous ip database", ipAddress.getHostAddress()));
				return null;
			} catch (IOException e) {
				LOGGER.error("AnonymousIp ERR: IO Error during lookup of ip in anonymous ip database", e);
				return null;
			}
		} else {
			return null;
		}
	}

	public boolean isInitialized() {
		return initialized;
	}

	/*
	 * Closes the database when the object is destroyed
	 */
	@Override
	protected void finalize() throws Throwable {
		if (databaseReader != null) {
			try {
				databaseReader.close();
				databaseReader = null;
			} catch (IOException e) {
				LOGGER.warn("Caught exception while trying to close anonymous ip database reader: ", e);
			}
		}
		super.finalize();
	}

}
