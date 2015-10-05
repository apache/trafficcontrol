package com.comcast.cdn.traffic_control.traffic_router.core.loc;

import java.util.List;

public class Federation {

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
}
