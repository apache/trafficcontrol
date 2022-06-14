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

import org.apache.traffic_control.traffic_router.core.util.AbstractResourceWatcher;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.apache.traffic_control.traffic_router.core.util.TrafficOpsUtils;

public class FederationsWatcher extends AbstractResourceWatcher {
    private static final Logger LOGGER = LogManager.getLogger(FederationsWatcher.class);
    private FederationRegistry federationRegistry;

    public static final String DEFAULT_FEDERATION_DATA_URL = "https://${toHostname}/api/"+TrafficOpsUtils.TO_API_VERSION+"/federations/all";
    public FederationsWatcher() {
        setDatabaseUrl(DEFAULT_FEDERATION_DATA_URL);
        setDefaultDatabaseUrl(DEFAULT_FEDERATION_DATA_URL);
    }

    @Override
    public boolean useData(final String data) {
        try {
            federationRegistry.setFederations(new FederationsBuilder().fromJSON(data));
            return true;
        } catch (Exception e) {
            LOGGER.warn("Failed updating federations data from " + dataBaseURL);
        }

        return false;
    }

    @Override
    protected boolean verifyData(final String data) {
        try {
            new FederationsBuilder().fromJSON(data);
            return true;
        } catch (Exception e) {
            LOGGER.warn("Failed to build federations data from " + dataBaseURL);
        }

        return false;
    }

    public void setFederationRegistry(final FederationRegistry federationRegistry) {
        this.federationRegistry = federationRegistry;
    }

    public String getWatcherConfigPrefix() {
        return "federationmapping";
    }
}
