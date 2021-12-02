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

public class AnonymousIpConfigUpdater extends AbstractServiceUpdater {
    private static final Logger LOGGER = LogManager.getLogger(AnonymousIpConfigUpdater.class);

    public AnonymousIpConfigUpdater() {
        LOGGER.debug("init...");
        sourceCompressed = false;
        tmpPrefix = "anonymousip";
        tmpSuffix = ".json";
    }
    
    @Override
    /*
     * Loads the anonymous ip config file
     */
    public boolean loadDatabase() throws IOException {
    	LOGGER.debug("AnonymousIpConfigUodater loading config");
        final File existingDB = databasesDirectory.resolve(databaseName).toFile();
        return AnonymousIp.parseConfigFile(existingDB, false);
    }

    @Override
    /*
     * Verifies the anonymous ip config file
     */
    public boolean verifyDatabase(final File dbFile) throws IOException {
    	LOGGER.debug("AnonymousIpConfigUpdater verifying config");
        return AnonymousIp.parseConfigFile(dbFile, true);
    }

}
