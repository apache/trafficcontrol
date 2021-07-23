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

public class GeolocationDatabaseUpdater extends AbstractServiceUpdater {
	private MaxmindGeolocationService maxmindGeolocationService;

	@Override
	public boolean verifyDatabase(final File dbFile) throws IOException {
		return maxmindGeolocationService.verifyDatabase(dbFile);
	}

	public boolean loadDatabase() throws IOException {
		maxmindGeolocationService.setDatabaseFile(databasesDirectory.resolve(databaseName).toFile());
		maxmindGeolocationService.reloadDatabase();
		return true;
	}

	@Override
	public boolean isLoaded() {
		if (maxmindGeolocationService != null) {
			return maxmindGeolocationService.isInitialized();
		}

		return loaded;
	}

	public void setMaxmindGeolocationService(final MaxmindGeolocationService maxmindGeolocationService) {
		this.maxmindGeolocationService = maxmindGeolocationService;
	}
}
