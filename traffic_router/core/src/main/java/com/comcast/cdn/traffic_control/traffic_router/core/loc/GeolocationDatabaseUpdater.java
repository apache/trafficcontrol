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

import org.apache.log4j.Logger;

public class GeolocationDatabaseUpdater extends AbstractServiceUpdater {
	private static final Logger LOGGER = Logger.getLogger(GeolocationDatabaseUpdater.class);

	public GeolocationDatabaseUpdater() {
	}

	private GeolocationService geoLocation;
	public void setGeoLocation(final GeolocationService geoLocation) {
		this.geoLocation = geoLocation;
	}

	public void verifyDatabase(final File dbFile) throws IOException {
		geoLocation.verifyDatabase(dbFile);
	}
	public boolean loadDatabase() throws IOException {
		LOGGER.info("Reloading location database.");
		geoLocation.reloadDatabase();
		LOGGER.info("Successfully reloaded location database.");
		return true;
	}

	@Override
	public boolean isLoaded() {
		if (geoLocation != null) {
			return geoLocation.isInitialized();
		}

		return loaded;
	}

}
