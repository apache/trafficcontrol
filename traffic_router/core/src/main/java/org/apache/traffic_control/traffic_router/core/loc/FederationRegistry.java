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

import org.apache.traffic_control.traffic_router.core.edge.InetRecord;
import org.apache.traffic_control.traffic_router.core.util.CidrAddress;
import org.apache.traffic_control.traffic_router.core.util.ComparableTreeSet;

import java.util.ArrayList;
import java.util.List;

public class FederationRegistry {
    private List<Federation> federations = new ArrayList<Federation>();

    public void setFederations(final List<Federation> federations) {
        synchronized (this.federations) {
            this.federations = federations;
        }
    }

    public List<InetRecord> findInetRecords(final String deliveryServiceId, final CidrAddress cidrAddress) {
        Federation targetFederation = null;

        synchronized (this.federations) {
            for (final Federation federation : federations) {
                if (deliveryServiceId.equals(federation.getDeliveryService())) {
                    targetFederation = federation;
                    break;
                }
            }
        }

        if (targetFederation == null) {
            return null;
        }

        for (final FederationMapping federationMapping : targetFederation.getFederationMappings()) {

            final ComparableTreeSet<CidrAddress> cidrAddresses = federationMapping.getResolveAddresses(cidrAddress);

            if (cidrAddresses == null) {
                continue;
            }

            for (final CidrAddress resolverAddress : cidrAddresses) {
                if (resolverAddress.equals(cidrAddress) || resolverAddress.includesAddress(cidrAddress)) {
                    return createInetRecords(federationMapping);
                }
            }
        }

        return null;
    }

    protected List<InetRecord> createInetRecords(final FederationMapping federationMapping) {
        final InetRecord inetRecord = new InetRecord(federationMapping.getCname(), federationMapping.getTtl());
        final List<InetRecord> inetRecords = new ArrayList<InetRecord>();
        inetRecords.add(inetRecord);
        return inetRecords;
    }

	public List<Federation> findFederations(final CidrAddress cidrAddress) {
		final List<Federation> results = new ArrayList<Federation>();

		for (final Federation federation : federations) {
			if (federation.containsCidrAddress(cidrAddress)) {
				results.add(federation);
			}
		}

		return results;
	}
}
