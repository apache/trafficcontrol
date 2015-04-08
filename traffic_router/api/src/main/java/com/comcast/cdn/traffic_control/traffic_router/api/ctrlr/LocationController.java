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

import java.util.HashMap;
import java.util.Map;

import javax.servlet.http.HttpServletRequest;

import org.apache.log4j.Logger;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.servlet.ModelAndView;

import com.comcast.cdn.traffic_control.traffic_router.api.util.APIAccessRecord;
import com.comcast.cdn.traffic_control.traffic_router.api.util.DataImporter;
import com.comcast.cdn.traffic_control.traffic_router.api.util.DataImporterException;

/**
 * Performs operations on the Popularity Zones available to the traffic router
 */
@Controller
@RequestMapping("/locations")
public class LocationController {
	private static final Logger LOGGER = Logger.getLogger(LocationController.class);
	private static final String VIEW_NAME = "JSONView";

	@RequestMapping(value = "/{locID}/caches", method = RequestMethod.GET)
	public ModelAndView getCaches(final HttpServletRequest request, @PathVariable("locID") final String locId) {
		final ModelAndView result = new ModelAndView(VIEW_NAME);

		try {
			final DataImporter di = new DataImporter("traffic-router:name=dataExporter");
			result.addObject("caches", di.invokeOperation("getCaches", locId));
		} catch (DataImporterException ex) {
			LOGGER.error(ex, ex);
			result.addObject("error", ex);
		}

		APIAccessRecord.log(request);

		return result;
	}

	/**
	 * Retrieves all of the locations.
	 * 
	 * @return the {@link ModelAndView} that represents the result.
	 */
	@RequestMapping(value = "", method = RequestMethod.GET)
	public ModelAndView getLocations(final HttpServletRequest request) {
		final ModelAndView result = new ModelAndView(VIEW_NAME);

		try {
			final DataImporter di = new DataImporter("traffic-router:name=dataExporter");
			result.addObject("locations", di.invokeOperation("getLocations"));
		} catch (DataImporterException ex) {
			LOGGER.error(ex, ex);
			result.addObject("error", ex);
		}

		APIAccessRecord.log(request);

		return result;
	}

	/**
	 * Retrieves all of the caches.
	 * 
	 * @return the {@link ModelAndView} that represents the result.
	 */
	@RequestMapping(value = "/caches", method = RequestMethod.GET)
	public ModelAndView getCaches(final HttpServletRequest request) {
		final ModelAndView result = new ModelAndView(VIEW_NAME);

		try {
			final DataImporter di = new DataImporter("traffic-router:name=dataExporter");

			/* Need to coerce this data structure into ModelAndView objects
			 *  to ensure that the resultant data structure matches the old version.
			 */
			@SuppressWarnings("unchecked")
			final Map<String, Object> map = (Map<String, Object>) di.invokeOperation("getCaches");
			final Map<String, ModelAndView> models = new HashMap<String, ModelAndView>();

			for (String key : map.keySet()) {
				final ModelAndView m = new ModelAndView(VIEW_NAME);
				m.addObject("caches", map.get(key));
				models.put(key, m);
			}

			result.addObject("locations", models);
		} catch (DataImporterException ex) {
			LOGGER.error(ex, ex);
			result.addObject("error", ex);
		}

		APIAccessRecord.log(request);

		return result;
	}
}
