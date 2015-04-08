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

package com.comcast.cdn.traffic_control.traffic_router.api.util;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import org.apache.log4j.Logger;
import org.springframework.web.servlet.HandlerInterceptor;
import org.springframework.web.servlet.ModelAndView;
import org.springframework.web.servlet.mvc.WebContentInterceptor;

public class ResponseCachingInterceptor extends WebContentInterceptor implements HandlerInterceptor {
	private static final Logger LOGGER = Logger.getLogger(ResponseCachingInterceptor.class);

	public final void postHandle(final HttpServletRequest request, final HttpServletResponse response, final Object handler, final ModelAndView modelAndView) {
		try {
			final DataImporter di = new DataImporter("traffic-router:name=dataExporter");
			final int maxAge = (Integer) di.invokeOperation("getCacheControlMaxAge");

			if (maxAge > 0) {
				cacheForSeconds(response, maxAge);
			} else if (maxAge == 0) {
				// this is our existing behavior
				response.addHeader("Pragma", "no-cache");
				response.addHeader("Cache-Control", "no-cache, no-store, max-age=0");
				response.addDateHeader("Expires", 1L);
			}
			// if -1, no cache control headers will be applied
		} catch (DataImporterException e) {
			LOGGER.error(e, e);
		}
	}
}
