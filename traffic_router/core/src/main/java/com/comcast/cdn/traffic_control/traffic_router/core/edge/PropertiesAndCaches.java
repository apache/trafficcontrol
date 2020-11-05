package com.comcast.cdn.traffic_control.traffic_router.core.edge;

import java.util.ArrayList;
import java.util.Map;

public class PropertiesAndCaches {
    public Map<String, String> properties;
    public ArrayList<String> caches;

    public PropertiesAndCaches(CacheLocation cacheLocation) {
        properties = cacheLocation.getProperties();
        caches = new ArrayList<>();
        for (Cache cache : cacheLocation.getCaches()) {
            caches.add(cache.getId());
        }
    }

    public Map<String, String> getProperties() {
        return properties;
    }

    public void setProperties(Map<String, String> properties) {
        this.properties = properties;
    }

    public ArrayList<String> getCaches() {
        return caches;
    }

    public void setCaches(ArrayList<String> caches) {
        this.caches = caches;
    }
}
