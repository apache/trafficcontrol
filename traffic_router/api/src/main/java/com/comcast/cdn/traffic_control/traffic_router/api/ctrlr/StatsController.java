/*
 * Copyright 2015 Comcast Cable Communications Management, LLC
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

package com.comcast.cdn.traffic_control.traffic_router.api.ctrlr;

import javax.servlet.http.HttpServletRequest;

import org.apache.log4j.Logger;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.servlet.ModelAndView;

import com.comcast.cdn.traffic_control.traffic_router.api.util.APIAccessRecord;
import com.comcast.cdn.traffic_control.traffic_router.api.util.DataImporter;
import com.comcast.cdn.traffic_control.traffic_router.api.util.DataImporterException;

/**
 * Performs operations on the Popularity Zones available to the traffic router
 */
@Controller
@RequestMapping("/stats")
public class StatsController {
	private static final Logger LOGGER = Logger.getLogger(StatsController.class);
	private static final String VIEW_NAME = "JSONView";

	/**
	 * Retrieves all of the locations.
	 * 
	 * @return the {@link ModelAndView} that represents the result.
	 */
	@RequestMapping(value = "", method = RequestMethod.GET)
	public ModelAndView getStats(final HttpServletRequest request) {
		final ModelAndView result = new ModelAndView(VIEW_NAME);

		try {
			final DataImporter di = new DataImporter("traffic-router:name=dataExporter");
			result.addObject("app", di.invokeOperation("getAppInfo"));
			result.addObject("stats", di.invokeOperation("getStatTracker"));
		} catch (DataImporterException ex) {
			LOGGER.error(ex, ex);
			result.addObject("error", ex);
		}

		APIAccessRecord.log(request);

		return result;
	}

	@RequestMapping(value = "/ip/{ip:.+}", method = RequestMethod.GET)
	public ModelAndView getCaches(final HttpServletRequest request, @PathVariable("ip") final String ip,
          @RequestParam(name = "geolocationProvider", required = false, defaultValue = "maxmindGeolocationService") final String geolocationProvider) {
		final ModelAndView result = new ModelAndView(VIEW_NAME);

		try {
			final DataImporter di = new DataImporter("traffic-router:name=dataExporter");
			result.addObject("result", di.invokeOperation("getCachesByIp", ip, geolocationProvider));
		} catch (DataImporterException ex) {
			LOGGER.error(ex, ex);
			result.addObject("error", ex);
		}

		APIAccessRecord.log(request);

		return result;
	}
}
