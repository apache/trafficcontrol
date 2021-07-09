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

import org.apache.traffic_control.traffic_router.core.util.CidrAddress;

import java.util.List;

public class Federation implements Comparable<Federation> {

    private final String deliveryService;
    private final List<FederationMapping> federationMappings;

    public Federation(final String deliveryService, final List<FederationMapping> federationMappings) {
        this.deliveryService = deliveryService;
        this.federationMappings = federationMappings;
    }

    public String getDeliveryService() {
        return deliveryService;
    }

    public List<FederationMapping> getFederationMappings() {
        return federationMappings;
    }

    @Override
    @SuppressWarnings({"PMD.IfStmtsMustUseBraces"})
    public boolean equals(final Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;

        final Federation that = (Federation) o;

        if (deliveryService != null ? !deliveryService.equals(that.deliveryService) : that.deliveryService != null)
            return false;
        return !(federationMappings != null ? !federationMappings.equals(that.federationMappings) : that.federationMappings != null);

    }

    @Override
    public int hashCode() {
        int result = deliveryService != null ? deliveryService.hashCode() : 0;
        result = 31 * result + (federationMappings != null ? federationMappings.hashCode() : 0);
        return result;
    }

    @Override
    public int compareTo(final Federation other) {
        return deliveryService.compareTo(other.deliveryService);
    }

    public boolean containsCidrAddress(final CidrAddress cidrAddress) {
        for (final FederationMapping federationMapping : federationMappings) {
            if (federationMapping.containsCidrAddress(cidrAddress)) {
                return true;
            }
        }

        return false;
    }
}
