package com.comcast.cdn.traffic_control.traffic_router.api.util;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.stereotype.Component;

import java.util.HashMap;
import java.util.Map;

@Component
public class TrafficRouterImporter {

    public Object fetchData(final String operation) throws Exception{
        return new DataImporter("traffic-router:name=dataExporter").invokeOperation(operation);
    }

    public String fetchStaticZoneStats() throws Exception {
        return new ObjectMapper().writeValueAsString(fetchData("getStaticZoneCacheStats"));
    }

    public String fetchDynamicZoneStats() throws Exception {
        return new ObjectMapper().writeValueAsString(fetchData("getDynamicZoneCacheStats"));
    }

    public String fetchZoneStats() throws Exception {
        final Map <String, Object> map = new HashMap<String, Object>();
        map.put("staticZoneCaches", fetchData("getStaticZoneCacheStats"));
        map.put("dynamicZoneCaches", fetchData("getDynamicZoneCacheStats"));
        return new ObjectMapper().writeValueAsString(map);
    }
}
