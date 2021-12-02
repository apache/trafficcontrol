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

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import java.io.File;
import java.io.IOException;

@SuppressWarnings({ "PMD.AvoidDuplicateLiterals" })
public class AnonymousIpDatabaseUpdater extends AbstractServiceUpdater {
    private static final Logger LOGGER = LogManager.getLogger(AnonymousIpDatabaseUpdater.class);

	private AnonymousIpDatabaseService anonymousIpDatabaseService;

	@Override
	/*
	 * Verifies the anonymous ip database
	 */
	public boolean verifyDatabase(final File dbFile) throws IOException {
		LOGGER.debug("Verifying Anonymous IP Database");
		return anonymousIpDatabaseService.verifyDatabase(dbFile);
	}

	/*
	 * Sets the anonymous ip database file and reloads the database
	 */
	public boolean loadDatabase() throws IOException {
		LOGGER.debug("Loading Anonymous IP Database");
		anonymousIpDatabaseService.setDatabaseFile(databasesDirectory.resolve(databaseName).toFile());
		anonymousIpDatabaseService.reloadDatabase();
		return true;
	}

	@Override
	/*
	 * Returns a boolean with the initialization state of the database
	 */
	public boolean isLoaded() {
		if (anonymousIpDatabaseService != null) {
			return anonymousIpDatabaseService.isInitialized();
		}

		return loaded;
	}
	
	public void setAnonymousIpDatabaseService(final AnonymousIpDatabaseService anonymousIpDatabaseService) {
		this.anonymousIpDatabaseService = anonymousIpDatabaseService;
	}

}
