package com.comcast.cdn.traffic_control.traffic_router.core.http;

import javax.servlet.http.HttpServletRequest;
import java.util.HashMap;
import java.util.Map;
import java.util.Set;

public class HttpAccessRequestHeaders {
    public Map<String, String> makeMap(final HttpServletRequest request, final Set<String> headerNames) {
        final Map<String, String> result = new HashMap<String, String>();

        for (String name : headerNames) {
            final String value = request.getHeader(name);
            if (value != null && !value.isEmpty()) {
                result.put(name, value);
            }
        }

        return result;
    }
}
