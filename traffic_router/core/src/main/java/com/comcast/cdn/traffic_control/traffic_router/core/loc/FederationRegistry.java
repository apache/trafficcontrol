package com.comcast.cdn.traffic_control.traffic_router.core.loc;

import com.comcast.cdn.traffic_control.traffic_router.core.cache.InetRecord;
import com.comcast.cdn.traffic_control.traffic_router.core.util.CidrAddress;
import com.comcast.cdn.traffic_control.traffic_router.core.util.ComparableTreeSet;

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

            ComparableTreeSet<CidrAddress> cidrAddresses;
            if (cidrAddress.isIpV6()) {
                cidrAddresses = federationMapping.getResolve6();
            }
            else {
                cidrAddresses = federationMapping.getResolve4();
            }

            for (CidrAddress resolverAddress : cidrAddresses) {
                if (resolverAddress.includesAddress(cidrAddress)) {
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

		for (Federation federation : federations) {
			if (federation.containsCidrAddress(cidrAddress)) {
				results.add(federation);
			}
		}

		return results;
	}
}
