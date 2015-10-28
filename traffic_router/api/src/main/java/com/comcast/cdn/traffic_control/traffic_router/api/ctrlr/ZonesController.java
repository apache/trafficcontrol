package com.comcast.cdn.traffic_control.traffic_router.api.ctrlr;

import com.comcast.cdn.traffic_control.traffic_router.api.util.TrafficRouterImporter;
import org.apache.log4j.Logger;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.ResponseBody;

@Controller
@RequestMapping("/stats/zones")
public class ZonesController {
    private static final Logger LOGGER = Logger.getLogger(ZonesController.class);
    @Autowired
    private TrafficRouterImporter trafficRouterImporter;

    @RequestMapping(value = "/caches")
    @ResponseBody
    public ResponseEntity<String> getAllCachesStats() {
        try {
            String json = trafficRouterImporter.fetchZoneStats();
            return new ResponseEntity<String>(json, HttpStatus.OK);
        }
        catch (Exception e) {
            LOGGER.error("Failed collecting zone cache statistics");
            return new ResponseEntity<String>(HttpStatus.INTERNAL_SERVER_ERROR);
        }
    }

    @RequestMapping(value = "/caches/{filter:static|dynamic}" , method = RequestMethod.GET)
    @ResponseBody
    public ResponseEntity<String> getCachesStats(@PathVariable("filter") String filter) {
        String json = null;

        try {
            if ("static".equals(filter)) {
                json = trafficRouterImporter.fetchStaticZoneStats();
            }

            if ("dynamic".equals(filter)) {
                json = trafficRouterImporter.fetchDynamicZoneStats();
            }
            return new ResponseEntity<String>(json, HttpStatus.OK);
        }
        catch (Exception e) {
            LOGGER.error("Failed collecting zone cache statistics");
            return new ResponseEntity<String>(HttpStatus.INTERNAL_SERVER_ERROR);
        }
    }
}